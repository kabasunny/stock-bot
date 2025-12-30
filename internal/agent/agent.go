package agent

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/app" // Add this import
	"stock-bot/internal/infrastructure/client"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

// Agent は取引エージェントのメイン構造体
type Agent struct {
	configPath       string
	config           *AgentConfig
	logger           *slog.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	signalPattern    string
	state            *State
	tradeService     TradeService
	eventClient      client.EventClient
	positionRepo     repository.PositionRepository
	executionUseCase app.ExecutionUseCase // New field
	gyouNoToSymbol   map[string]string    // 行番号(文字列)から銘柄コードへのマッピング
}

// NewAgent は新しいエージェントのインスタンスを作成する
// tradeService はトレードサービス（Go APIラッパー）の実装
func NewAgent(configPath string, tradeService TradeService, eventClient client.EventClient, positionRepo repository.PositionRepository, executionUseCase app.ExecutionUseCase) (*Agent, error) {
	// 先に設定を読み込んでおく
	cfg, err := LoadAgentConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // TODO: ログレベルを設定ファイルから反映させる

	ctx, cancel := context.WithCancel(context.Background())

	gyouNoToSymbol := make(map[string]string)
	for i, symbol := range cfg.StrategySettings.Swingtrade.TargetSymbols {
		gyouNoToSymbol[fmt.Sprintf("%d", i+1)] = symbol // 1-indexed
	}

	return &Agent{
		configPath:       configPath,
		config:           cfg,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		signalPattern:    cfg.StrategySettings.Swingtrade.SignalFilePattern, // とりあえずスイングトレードに固定
		state:            NewState(),
		tradeService:     tradeService,
		eventClient:      eventClient,
		positionRepo:     positionRepo,
		executionUseCase: executionUseCase, // Initialize new field
		gyouNoToSymbol:   gyouNoToSymbol,   // マッピングを初期化
	}, nil
}

// Start はエージェントの実行ループを開始する
func (a *Agent) Start() {
	a.logger.Info("starting agent...")
	// WebSocketイベント監視をゴルーチンで開始
	go a.watchEvents()

	// 定期実行のTicker
	ticker := time.NewTicker(a.config.Agent.ExecutionInterval)
	defer ticker.Stop()
	// WebSocketクライアントのクリーンアップ
	defer a.eventClient.Close()

	// 起動時に初期状態を同期
	a.syncInitialState()

	// 起動時に一度実行 (初期状態同期後にtickを実行)
	a.tick()

	for {
		select {
		case <-ticker.C:
			a.tick()
		case <-a.ctx.Done():
			a.logger.Info("agent stopping...")
			return
		}
	}
}

// watchEvents はWebSocketイベントを監視し、ログに出力する
func (a *Agent) watchEvents() {
	a.logger.Info("starting event watcher...")

	session := a.tradeService.GetSession()
	if session == nil {
		a.logger.Error("failed to get session for event watcher")
		return
	}

	targetSymbols := a.config.StrategySettings.Swingtrade.TargetSymbols
	if len(targetSymbols) == 0 {
		a.logger.Warn("no target symbols defined in config for WebSocket subscription, connecting without symbols")
	}

	messages, errs, err := a.eventClient.Connect(a.ctx, session, targetSymbols)
	if err != nil {
		a.logger.Error("failed to connect to event stream", "error", err)
		return
	}
	a.logger.Info("event watcher connected to WebSocket")

	for {
		select {
		case msgBytes, ok := <-messages:
			if !ok {
				a.logger.Info("message channel closed, stopping event watcher")
				return
			}
			// 1. パース処理
			parsedMsg, err := a.parseEventMessage(msgBytes)
			if err != nil {
				a.logger.Error("failed to parse websocket event", "error", err, "raw_message", string(msgBytes))
				continue
			}

			// 2. イベント種別を取得
			cmd, ok := parsedMsg["p_cmd"]
			if !ok {
				a.logger.Warn("p_cmd not found in websocket event", "message", parsedMsg)
				continue
			}

			a.logger.Info("received websocket event", "command", cmd)

			// 3. 振り分け
			switch cmd {
			case "FD": // 時価配信データ (Feed Data)
				a.handlePriceData(parsedMsg)
			case "ST": // ステータス通知 (Status)
				a.handleStatus(parsedMsg)
			case "EC": // 約定通知 (Execution)
				// handleExecution は executionUseCase.Execute を呼び出す
				// executionUseCase.Execute は orderRepo.UpdateOrderStatusByExecution を呼び出す
				// orderRepo.UpdateOrderStatusByExecution が order with ID not found エラーを返す場合がある

				err := a.handleExecution(parsedMsg) // handleExecution がエラーを返すように変更した
				if err != nil {
					// "order with ID ... not found" エラーの場合は Warn レベルでログを出力
					if strings.Contains(err.Error(), "order with ID") && strings.Contains(err.Error(), "not found") {
						a.logger.Warn("execution event for non-existent order received", "error", err)
					} else {
						a.logger.Error("failed to handle execution event", "error", err)
					}
				}
			case "KP": // キープアライブ (Keep Alive)
				// 特段の処理は不要だが、ログの警告を避けるためここに記述
				a.logger.Debug("received keep alive event", "details", parsedMsg) // Debugレベルでログ出力
			default:
				a.logger.Warn("unhandled websocket event command", "command", cmd, "details", parsedMsg)
			}
		case err, ok := <-errs:
			if !ok {
				a.logger.Info("error channel closed, stopping event watcher")
				return
			}
			a.logger.Error("received error from event stream", "error", err)
			// TODO: エラー内容に応じた再接続処理などを検討
			return // エラーが発生したら一旦ウォッチャーを終了
		case <-a.ctx.Done():
			a.logger.Info("agent context done, stopping event watcher")
			return
		}
	}
}

// parseEventMessage はWebSocketのカスタムフォーマットメッセージをパースする
func (a *Agent) parseEventMessage(msg []byte) (map[string]string, error) {
	result := make(map[string]string)
	// メッセージが空、または改行コードのみの場合を除外
	if len(bytes.TrimSpace(msg)) == 0 {
		return result, nil
	}
	pairs := bytes.Split(msg, []byte{0x01}) // ^A で分割
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}
		kv := bytes.SplitN(pair, []byte{0x02}, 2) // ^B でキーと値に分割
		if len(kv) != 2 {
			// キーだけのペア（例: `p_no^B1`の後ろに`^A`がない場合など）も許容する
			if len(kv) == 1 && len(kv[0]) > 0 {
				result[string(kv[0])] = ""
				continue
			}
			// 不正な形式のペアは無視する
			a.logger.Warn("invalid key-value pair format in websocket message", "pair", string(pair))
			continue
		}

		key := string(kv[0])
		value := string(kv[1])
		result[key] = value
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("message parsing resulted in no key-value pairs: %s", string(msg))
	}
	return result, nil
}

// handleExecution は約定通知イベントを処理する
func (a *Agent) handleExecution(data map[string]string) error {
	// a.logger.Info("handling execution event", "data", data) // 詳細ログを抑制

	// map[string]string を model.Execution に変換
	execution := &model.Execution{}

	// 約定通知 (EC) イベントのキー名に合わせてマッピング
	// ExecutionID は約定ごとにユニークなIDが必要だが、ECイベントには直接存在しないため、p_ON (注文ID) と p_ENO (連番) を組み合わせる
	if orderID, ok := data["p_ON"]; ok {
		execution.OrderID = orderID
		if executionNo, ok := data["p_ENO"]; ok {
			execution.ExecutionID = fmt.Sprintf("%s-%s", orderID, executionNo) // 注文ID-約定番号
		} else {
			execution.ExecutionID = orderID // 約定番号がない場合は注文IDのみ
		}
	} else {
		a.logger.Error("EC event missing p_ON (OrderID)", "data", data)
		return errors.New("EC event missing p_ON (OrderID)")
	}

	// Symbol
	if val, ok := data["p_IC"]; ok { // p_IC は銘柄コード
		execution.Symbol = val
	} else {
		a.logger.Error("EC event missing p_IC (Symbol)", "data", data)
		return errors.New("EC event missing p_IC (Symbol)")
	}

	// TradeType
	if val, ok := data["p_ST"]; ok { // p_ST は売買区分 (1:買, 2:売)
		if val == "1" {
			execution.TradeType = model.TradeTypeBuy
		} else if val == "2" {
			execution.TradeType = model.TradeTypeSell
		} else {
			a.logger.Error("invalid p_ST (TradeType) in EC event", "p_ST", val, "data", data)
			return errors.Errorf("invalid p_ST (TradeType) in EC event: %s", val)
		}
	} else {
		a.logger.Error("EC event missing p_ST (TradeType)", "data", data)
		return errors.New("EC event missing p_ST (TradeType)")
	}

	// Quantity
	if val, ok := data["p_EXSR"]; ok { // p_EXSR は約定数量
		qty, err := parseInt(val)
		if err != nil {
			a.logger.Error("invalid p_EXSR (Quantity) in EC event", "quantity", val, "error", err, "data", data)
			return errors.Wrapf(err, "invalid p_EXSR (Quantity) in EC event: %s", val)
		}
		execution.Quantity = qty
	} else {
		a.logger.Error("EC event missing p_EXSR (Quantity)", "data", data)
		return errors.New("EC event missing p_EXSR (Quantity)")
	}

	// Price
	if val, ok := data["p_EXPR"]; ok { // p_EXPR は約定単価
		price, err := parseFloat(val)
		if err != nil {
			a.logger.Error("invalid p_EXPR (Price) in EC event", "price", val, "error", err, "data", data)
			return errors.Wrapf(err, "invalid p_EXPR (Price) in EC event: %s", val)
		}
		execution.Price = price
	} else {
		a.logger.Error("EC event missing p_EXPR (Price)", "data", data)
		return errors.New("EC event missing p_EXPR (Price)")
	}

	// ExecutedAt
	if val, ok := data["p_EXDT"]; ok { // p_EXDT は約定日時 YYYYMMDDhhmmss
		executedAt, err := parseTime(val) // parseTime を使用
		if err != nil {
			a.logger.Error("invalid p_EXDT (ExecutedAt) in EC event", "executed_at", val, "error", err, "data", data)
			return errors.Wrapf(err, "invalid p_EXDT (ExecutedAt) in EC event: %s", val)
		}
		execution.ExecutedAt = executedAt
	} else {
		a.logger.Warn("EC event missing p_EXDT (ExecutedAt), using current time", "data", data)
		execution.ExecutedAt = time.Now()
	}

	// Commission (optional): ECイベントログには見当たらないため、0とする
	execution.Commission = 0

	// ExecutionUseCase を実行
	if err := a.executionUseCase.Execute(a.ctx, execution); err != nil {
		a.logger.Error("failed to execute execution use case", "execution_id", execution.ExecutionID, "error", err)
		return errors.Wrapf(err, "failed to execute execution use case for execution_id %s", execution.ExecutionID)
	} else {
		a.logger.Debug("successfully processed execution event", "execution_id", execution.ExecutionID, "order_id", execution.OrderID)
	}
	return nil
}

// parseInt は文字列をintにパースするヘルパー関数
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat は文字列をfloat64にパースするヘルパー関数
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseTime は文字列をtime.Timeにパースするヘルパー関数
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
		"20060102150405",            // YYYYMMDDhhmmss (今回追加)
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse time string: %s", s)
}

// handlePriceData は価格情報（時価配信）イベントを処理する
func (a *Agent) handlePriceData(data map[string]string) {
	// a.logger.Info("Handling Price Data", "data", data) // 詳細ログを抑制

	// FDイベントのデータは p_行番号_項目名 という形式で来る
	// 例: p_1_DPP -> 行番号1の銘柄の現在値
	// 行番号を特定し、gyouNoToSymbolマップから銘柄コードを取得する

	for key, value := range data {
		// p_N_DPP の形式を想定 (Nは行番号)
		if strings.HasPrefix(key, "p_") && strings.HasSuffix(key, "_DPP") {
			parts := strings.Split(key, "_")
			if len(parts) != 3 {
				continue // 予期しない形式
			}
			gyouNo := parts[1] // 行番号文字列

			symbol, ok := a.gyouNoToSymbol[gyouNo]
			if !ok {
				a.logger.Warn("unknown gyouNo in price data", "gyouNo", gyouNo, "key", key)
				continue
			}

			price, err := parseFloat(value) // ヘルパー関数を再利用
			if err != nil {
				a.logger.Error("failed to parse price from FD event", "symbol", symbol, "key", key, "value", value, "error", err)
				continue
			}

			a.state.UpdatePrice(symbol, price)
			a.logger.Debug("updated price from FD event", "symbol", symbol, "price", price) // Debugレベルでログ出力
		}
	}
}

// handleStatus はステータス通知イベントを処理する
func (a *Agent) handleStatus(data map[string]string) {
	a.logger.Warn("Received unhandled Status Notification (ST) event", "data", data)
	// TODO: セッション状態などを確認し、必要に応じて再接続などの処理を行う
	// 例: "session inactive" を検知してエージェントを安全に停止させるなど
}

// SetLogger は外部からロガーを注入するために使用します。
func (a *Agent) SetLogger(logger *slog.Logger) {
	a.logger = logger
}

// Tick はエージェントの単一の評価サイクルを実行します。バックテスト用に公開されています。
func (a *Agent) Tick() {
	a.tick()
}

// Stop はエージェントの実行ループを停止する
func (a *Agent) Stop() {
	a.logger.Info("sending stop signal to agent...")
	a.cancel()
}

// SyncInitialState はエージェントの初期状態を同期します。バックテスト用に公開されています。
func (a *Agent) SyncInitialState() {
	a.syncInitialState()
}

// syncInitialState はエージェント起動時にトレードサービスから状態を取得し、内部状態を同期する
func (a *Agent) syncInitialState() { // <<<<<<<<<<<<<<<< 追加
	a.logger.Info("synchronizing initial state...")
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second) // タイムアウトを設定
	defer cancel()

	// 残高の同期
	balance, err := a.tradeService.GetBalance(ctx)
	if err != nil {
		a.logger.Error("failed to get initial balance", "error", err)
		// エージェントの起動は継続するが、残高情報がない状態で動作することになるため、適切にハンドリングする必要がある
	} else {
		a.state.UpdateBalance(balance)
		a.logger.Info("initial balance synchronized", "cash", balance.Cash, "buying_power", balance.BuyingPower)
	}

	// ポジションの同期
	positions, err := a.tradeService.GetPositions(ctx)
	if err != nil {
		a.logger.Error("failed to get initial positions", "error", err)
	} else {
		a.state.UpdatePositions(positions)
		a.logger.Info("initial positions synchronized", "count", len(positions))
	}

	// 注文の同期
	orders, err := a.tradeService.GetOrders(ctx)
	if err != nil {
		a.logger.Error("failed to get initial orders", "error", err)
	} else {
		a.state.UpdateOrders(orders)
		a.logger.Info("initial orders synchronized", "count", len(orders))
	}

	a.logger.Info("initial state synchronization completed.")
}

// tick はループごとに実行される処理
func (a *Agent) tick() {
	a.logger.Info("agent tick")

	// tickの開始時に状態を完全に同期する
	syncCtx, syncCancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer syncCancel()

	// 残高の同期
	balance, err := a.tradeService.GetBalance(syncCtx)
	if err != nil {
		a.logger.Error("failed to sync balance in tick", "error", err)
		// 致命的なエラーではないため、処理は続行
	} else {
		a.state.UpdateBalance(balance)
		a.logger.Info("balance synchronized in tick", "buying_power", balance.BuyingPower)
	}

	// ポジションの同期
	positions, err := a.tradeService.GetPositions(syncCtx)
	if err != nil {
		a.logger.Error("failed to sync positions in tick", "error", err)
	} else {
		a.state.UpdatePositions(positions)
		a.logger.Info("positions synchronized in tick", "count", len(positions))
	}

	// 注文の同期
	orders, err := a.tradeService.GetOrders(syncCtx)
	if err != nil {
		a.logger.Error("failed to sync orders in tick", "error", err)
	} else {
		a.state.UpdateOrders(orders)
		a.logger.Info("orders synchronized in tick", "count", len(orders))
	}

	// 注文処理用のコンテキストをここで一度生成する
	orderCtx, orderCancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer orderCancel()

	// 1. 保有ポジションの決済チェック (利確・損切り)
	a.checkPositionsForExit(orderCtx)

	// 2. 新規エントリーのシグナルチェック
	a.checkSignalsForEntry(orderCtx)
}

// checkPositionsForExit は保有ポジションを監視し、利確または損切りの条件を満たしているか確認する
func (a *Agent) checkPositionsForExit(ctx context.Context) {
	a.logger.Info("checking positions for exit...")
	positions := a.state.GetPositions()
	if len(positions) == 0 {
		return
	}

	for _, pos := range positions {
		// すでにこの銘柄に対する決済注文（売り注文）が出ていないか確認
		hasOpenSellOrder := false
		for _, order := range a.state.GetOrders() {
			if order.Symbol == pos.Symbol && order.TradeType == model.TradeTypeSell && order.IsUnexecuted() {
				hasOpenSellOrder = true
				break
			}
		}
		if hasOpenSellOrder {
			a.logger.Info("skipping exit check due to existing open sell order", "symbol", pos.Symbol)
			continue
		}

		currentPrice, ok := a.state.GetPrice(pos.Symbol)
		if !ok {
			a.logger.Error("failed to get price for exit check from state", "symbol", pos.Symbol)
			continue
		}
		if currentPrice == 0 {
			a.logger.Warn("skipping exit check because current price is zero", "symbol", pos.Symbol)
			continue
		}

		a.updateTrailingStopState(ctx, pos, currentPrice)

		// 決済条件の判定 (優先順位: ATRベース損切り -> トレーリングストップ -> 固定利確)

		// 1. ATRベース損切り
		atr, err := a.getATRForPosition(ctx, pos.Symbol)
		if err == nil {
			if a.shouldStopLossFixed(pos, currentPrice, atr) {
				a.placeExitOrder(ctx, pos, "STOP_LOSS_FIXED")
				continue // 注文発行後はこのポジションの他のチェックは不要
			}
		} else {
			// ATRが取得できなくても他の決済ロジックは継続する
			a.logger.Warn("could not get ATR for stop-loss check, skipping ATR-based stop", "symbol", pos.Symbol, "error", err)
		}

		// 2. トレーリングストップ
		if a.shouldStopLossTrailing(pos, currentPrice) {
			a.placeExitOrder(ctx, pos, "STOP_LOSS_TRAILING")
			continue
		}

		// 3. 固定利確
		if a.shouldProfitTakeFixed(pos, currentPrice) {
			a.placeExitOrder(ctx, pos, "PROFIT_TAKE_FIXED")
			continue
		}
	}
}

// getATRForPosition は指定された銘柄のATRを計算し、エラーハンドリングを行う
func (a *Agent) getATRForPosition(ctx context.Context, symbol string) (float64, error) {
	atrPeriod := a.config.StrategySettings.Swingtrade.ATRPeriod
	history, err := a.tradeService.GetPriceHistory(ctx, symbol, atrPeriod+1)
	if err != nil {
		return 0, fmt.Errorf("failed to get price history for ATR: %w", err)
	}
	if len(history) < atrPeriod+1 {
		return 0, fmt.Errorf("not enough historical data for ATR calculation (required: %d, got: %d)", atrPeriod+1, len(history))
	}
	atr, err := calculateATR(history, atrPeriod)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate ATR: %w", err)
	}
	if atr == 0 {
		return 0, fmt.Errorf("calculated ATR is zero")
	}
	return atr, nil
}

// updateTrailingStopState はトレーリングストップの状態を更新する
func (a *Agent) updateTrailingStopState(ctx context.Context, pos *model.Position, currentPrice float64) {
	trailingStopTriggerRate := a.config.StrategySettings.Swingtrade.TrailingStopTriggerRate
	trailingStopRate := a.config.StrategySettings.Swingtrade.TrailingStopRate

	// Positional data initialization/update for trailing stop
	if pos.HighestPrice == 0 || currentPrice > pos.HighestPrice {
		pos.HighestPrice = currentPrice
		a.state.UpdatePositionHighestPrice(pos.Symbol, currentPrice)
		if err := a.positionRepo.UpdateHighestPrice(ctx, pos.Symbol, currentPrice); err != nil {
			a.logger.Error("failed to update highest price in db", "symbol", pos.Symbol, "error", err)
			// ここではエラーをログに出力するだけで、処理は続行する
		}
	}

	trailingStopTriggerPrice := pos.AveragePrice * (1 + trailingStopTriggerRate/100)

	if pos.HighestPrice > 0 { // HighestPriceが記録されている場合のみ
		calculatedTrailingStopPrice := pos.HighestPrice * (1 - trailingStopRate/100)
		if pos.TrailingStopPrice == 0 && currentPrice >= trailingStopTriggerPrice {
			// トレーリングストップがまだトリガーされておらず、トリガー条件を満たした場合
			pos.TrailingStopPrice = calculatedTrailingStopPrice
			a.state.UpdatePositionTrailingStopPrice(pos.Symbol, pos.TrailingStopPrice)
			a.logger.Info("trailing stop activated", "symbol", pos.Symbol, "trigger_price", trailingStopTriggerPrice, "initial_stop_price", pos.TrailingStopPrice)
		} else if pos.TrailingStopPrice > 0 && calculatedTrailingStopPrice > pos.TrailingStopPrice {
			// トレーリングストップが既に有効で、損切りラインが切り上がった場合
			pos.TrailingStopPrice = calculatedTrailingStopPrice
			a.state.UpdatePositionTrailingStopPrice(pos.Symbol, pos.TrailingStopPrice)
			a.logger.Info("trailing stop price updated", "symbol", pos.Symbol, "new_stop_price", pos.TrailingStopPrice)
		}
	}
}

// shouldProfitTakeFixed は固定利確条件を満たしているか判定する
func (a *Agent) shouldProfitTakeFixed(pos *model.Position, currentPrice float64) bool {
	profitTakeRate := a.config.StrategySettings.Swingtrade.ProfitTakeRate
	profitTakePrice := pos.AveragePrice * (1 + profitTakeRate/100)
	if currentPrice >= profitTakePrice {
		a.logger.Info("profit take condition met (fixed)", "symbol", pos.Symbol, "average_price", pos.AveragePrice, "current_price", currentPrice, "target_price", profitTakePrice)
		return true
	}
	return false
}

// shouldStopLossFixed はATRベースの固定損切り条件を満たしているか判定する
func (a *Agent) shouldStopLossFixed(pos *model.Position, currentPrice, atr float64) bool {
	stopLossATRMultiplier := a.config.StrategySettings.Swingtrade.StopLossATRMultiplier
	stopLossPrice := pos.AveragePrice - (atr * stopLossATRMultiplier) // ATRベースの損切り
	if currentPrice <= stopLossPrice {
		a.logger.Info("stop loss condition met (fixed)", "symbol", pos.Symbol, "average_price", pos.AveragePrice, "current_price", currentPrice, "target_price", stopLossPrice)
		return true
	}
	return false
}

// shouldStopLossTrailing はトレーリングストップ条件を満たしているか判定する
func (a *Agent) shouldStopLossTrailing(pos *model.Position, currentPrice float64) bool {
	if pos.TrailingStopPrice > 0 && currentPrice <= pos.TrailingStopPrice {
		a.logger.Info("stop loss condition met (trailing)", "symbol", pos.Symbol, "highest_price", pos.HighestPrice, "current_price", currentPrice, "trailing_stop_price", pos.TrailingStopPrice)
		return true
	}
	return false
}

// placeExitOrder は決済注文を生成し、発行するヘルパー関数
func (a *Agent) placeExitOrder(ctx context.Context, pos *model.Position, reason string) {
	req := &PlaceOrderRequest{
		Symbol:              pos.Symbol,
		TradeType:           model.TradeTypeSell,
		OrderType:           model.OrderTypeMarket,
		Quantity:            pos.Quantity,
		PositionAccountType: pos.PositionAccountType,
	}
	order, err := a.tradeService.PlaceOrder(ctx, req)
	if err != nil {
		a.logger.Error("failed to place exit order", "symbol", pos.Symbol, "reason", reason, "error", err)
		return
	}
	a.state.AddOrder(order)
	a.logger.Info("successfully placed exit order", "symbol", pos.Symbol, "order_id", order.OrderID, "reason", reason)
}

// checkSignalsForEntry はシグナルファイルをチェックし、新規エントリー注文を行う
func (a *Agent) checkSignalsForEntry(ctx context.Context) {
	a.logger.Info("checking signals for entry...")
	// 状態の確認（ログ出力のみ）
	balance := a.state.GetBalance()
	a.logger.Info("current balance", "cash", balance.Cash, "buying_power", balance.BuyingPower)
	positions := a.state.GetPositions()
	a.logger.Info("current positions", "count", len(positions))
	for _, p := range positions {
		a.logger.Info("  position detail", "symbol", p.Symbol, "quantity", p.Quantity, "average_price", p.AveragePrice)
	}
	orders := a.state.GetOrders()
	a.logger.Info("current orders", "count", len(orders))
	for _, o := range orders {
		a.logger.Info("  order detail", "order_id", o.OrderID, "symbol", o.Symbol, "trade_type", o.TradeType, "status", o.OrderStatus)
	}
	signalFilePath, err := FindSignalFile(a.signalPattern)
	if err != nil {
		a.logger.Error("failed to find signal file", "error", err)
		return
	}
	if signalFilePath == "" {
		a.logger.Info("no signal file found, skipping entry check")
		return
	}

	a.logger.Info("found signal file", "path", signalFilePath)

	signals, err := ReadSignalFile(signalFilePath)
	if err != nil {
		a.logger.Error("failed to read signal file", "path", signalFilePath, "error", err)
		return
	}

	a.logger.Info("signals loaded", "count", len(signals))

	// 購読対象の銘柄コードのセットを作成
	targetSymbolsMap := make(map[string]struct{})
	for _, s := range a.config.StrategySettings.Swingtrade.TargetSymbols {
		targetSymbolsMap[s] = struct{}{}
	}

	for _, s := range signals {
		symbolStr := fmt.Sprintf("%d", s.Symbol)

		// 購読対象の銘柄でなければスキップ
		if _, ok := targetSymbolsMap[symbolStr]; !ok {
			a.logger.Info("skipping signal for non-target symbol", "symbol", symbolStr)
			continue
		}

		a.logger.Info("processing signal for target symbol", "symbol", symbolStr, "signal", s.Signal)

		// Check for existing open orders for this symbol before processing the signal
		hasOpenOrder := false
		for _, o := range a.state.GetOrders() {
			if o.Symbol == symbolStr && o.OrderStatus.IsUnexecuted() {
				hasOpenOrder = true
				break
			}
		}

		if hasOpenOrder {
			a.logger.Info("skipping signal due to existing open order", "symbol", symbolStr)
			continue
		}

		// 意思決定ロジック
		if s.Signal == BuySignal {
			if _, ok := a.state.GetPosition(symbolStr); ok {
				a.logger.Info("skipping buy signal for already held position", "symbol", symbolStr)
				continue
			}
			a.logger.Info("preparing to place buy order", "symbol", symbolStr)

			// 買付余力と現在価格を取得
			balance := a.state.GetBalance()

			// WebSocketから取得した現在価格を使用
			currentPrice, ok := a.state.GetPrice(symbolStr)
			if !ok {
				a.logger.Warn("current price not available in state for sizing", "symbol", symbolStr)
				continue
			}
			if currentPrice == 0 {
				a.logger.Warn("skipping buy signal because current price in state is zero", "symbol", symbolStr)
				continue
			}

			atrPeriod := a.config.StrategySettings.Swingtrade.ATRPeriod
			stopLossATRMultiplier := a.config.StrategySettings.Swingtrade.StopLossATRMultiplier
			unitSize := float64(a.config.StrategySettings.Swingtrade.UnitSize)
			maxPositionSizePercentage := a.config.StrategySettings.Swingtrade.MaxPositionSizePercentage

			// 履歴価格データを取得
			history, err := a.tradeService.GetPriceHistory(ctx, symbolStr, atrPeriod+1) // ATR計算に必要な期間 + 1
			if err != nil {
				a.logger.Error("failed to get price history for ATR calculation", "symbol", symbolStr, "error", err)
				continue
			}
			if len(history) < atrPeriod+1 {
				a.logger.Warn("not enough historical data for ATR calculation, skipping sizing", "symbol", symbolStr, "required", atrPeriod+1, "got", len(history))
				continue
			}

			atr, err := calculateATR(history, atrPeriod)
			if err != nil {
				a.logger.Error("failed to calculate ATR", "symbol", symbolStr, "error", err)
				continue
			}
			if atr == 0 {
				a.logger.Warn("calculated ATR is zero, skipping sizing", "symbol", symbolStr)
				continue
			}

			// ポジションサイズをATRに基づいて計算
			// totalRiskAmount はポートフォリオ全体のリスク許容額（許容損失額）
			totalRiskAmount := balance.BuyingPower * a.config.StrategySettings.Swingtrade.TradeRiskPercentage
			// 1株あたりのボラティリティリスク (ATRに基づく)
			riskPerShare := atr * stopLossATRMultiplier
			if riskPerShare == 0 {
				a.logger.Warn("risk per share is zero, skipping sizing", "symbol", symbolStr)
				continue
			}

			// ATRベースで計算される最大許容株数
			maxSharesByATR := totalRiskAmount / riskPerShare

			// 買付余力から計算される最大購入株数
			maxSharesByBuyingPower := balance.BuyingPower / currentPrice
			if maxSharesByBuyingPower <= 0 {
				a.logger.Warn("insufficient buying power to purchase even one unit", "symbol", symbolStr)
				continue
			}

			// ポジションサイズ上限から計算される株数
			maxPositionValue := balance.BuyingPower * maxPositionSizePercentage
			maxSharesByPositionLimit := maxPositionValue / currentPrice

			// 3つの条件のうち最も小さい値を採用
			maxShares := math.Min(maxSharesByATR, maxSharesByBuyingPower)
			maxShares = math.Min(maxShares, maxSharesByPositionLimit)

			// unitSizeの倍数に切り捨て
			quantity := math.Floor(maxShares/unitSize) * unitSize

			a.logger.Info("calculated order quantity (ATR-based)",
				"symbol", symbolStr,
				"buying_power", balance.BuyingPower,
				"trade_risk_percentage", a.config.StrategySettings.Swingtrade.TradeRiskPercentage,
				"max_position_size_percentage", maxPositionSizePercentage,
				"atr_period", atrPeriod,
				"stop_loss_atr_multiplier", stopLossATRMultiplier,
				"calculated_atr", atr,
				"risk_per_share", riskPerShare,
				"total_risk_amount", totalRiskAmount,
				"max_shares_by_atr", maxSharesByATR,
				"max_shares_by_buying_power", maxSharesByBuyingPower,
				"max_shares_by_position_limit", maxSharesByPositionLimit,
				"final_max_shares", maxShares,
				"calculated_quantity", quantity)

			if quantity <= 0 {
				a.logger.Info("skipping buy signal due to zero calculated quantity", "symbol", symbolStr)
				continue
			}

			// 注文リクエストを作成
			req := &PlaceOrderRequest{
				Symbol:              symbolStr,
				TradeType:           model.TradeTypeBuy,
				OrderType:           model.OrderTypeMarket,
				Quantity:            int(quantity),
				Price:               0,                             // 成行注文のため価格は0
				PositionAccountType: model.PositionAccountTypeCash, // 新規買い注文は現物と仮定
			}

			// 注文を発行
			order, err := a.tradeService.PlaceOrder(ctx, req)
			if err != nil {
				a.logger.Error("failed to place buy order", "symbol", symbolStr, "error", err)
				continue // 次のシグナルへ
			}
			a.logger.Info("successfully placed buy order", "symbol", symbolStr, "order_id", order.OrderID)
			a.state.AddOrder(order) // 発注成功後、内部状態を更新する

		} else if s.Signal == SellSignal {
			position, ok := a.state.GetPosition(symbolStr)
			if !ok {
				a.logger.Info("skipping sell signal for non-held position", "symbol", symbolStr)
				continue
			}
			a.logger.Info("preparing to place sell order", "symbol", symbolStr, "quantity", position.Quantity)

			// 注文リクエストを作成
			req := &PlaceOrderRequest{
				Symbol:              symbolStr,
				TradeType:           model.TradeTypeSell,
				OrderType:           model.OrderTypeMarket,
				Quantity:            position.Quantity, // 保有する全数量を売却
				Price:               0,                 // 成行注文のため価格は0
				PositionAccountType: position.PositionAccountType,
			}

			// 注文を発行
			order, err := a.tradeService.PlaceOrder(ctx, req)
			if err != nil {
				a.logger.Error("failed to place sell order", "symbol", symbolStr, "error", err)
				continue // 次のシグナルへ
			}
			a.logger.Info("successfully placed sell order", "symbol", symbolStr, "order_id", order.OrderID)
			a.state.AddOrder(order) // 発注成功後、内部状態を更新する
		}
	}
}
