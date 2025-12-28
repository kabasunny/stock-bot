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
	"stock-bot/internal/infrastructure/client"
	"time"
)

// Agent は取引エージェントのメイン構造体
type Agent struct {
	configPath    string
	config        *AgentConfig
	logger        *slog.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	signalPattern string
	state         *State
	tradeService  TradeService
	eventClient   client.EventClient
	positionRepo  repository.PositionRepository
}

// NewAgent は新しいエージェントのインスタンスを作成する
// tradeService はトレードサービス（Go APIラッパー）の実装
func NewAgent(configPath string, tradeService TradeService, eventClient client.EventClient, positionRepo repository.PositionRepository) (*Agent, error) {
	// 先に設定を読み込んでおく
	cfg, err := LoadAgentConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // TODO: ログレベルを設定ファイルから反映させる

	ctx, cancel := context.WithCancel(context.Background())

	return &Agent{
		configPath:    configPath,
		config:        cfg,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		signalPattern: cfg.StrategySettings.Swingtrade.SignalFilePattern, // とりあえずスイングトレードに固定
		state:         NewState(),
		tradeService:  tradeService,
		eventClient:   eventClient,
		positionRepo:  positionRepo,
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

	messages, errs, err := a.eventClient.Connect(a.ctx, session)
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
			// TODO: 約定通知など、他のコマンドもここに追加していく
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

// --- プレースホルダー関数 ---

// handlePriceData は価格情報（時価配信）イベントを処理する
func (a *Agent) handlePriceData(data map[string]string) {
	a.logger.Info("[Placeholder] Handling Price Data", "data", data)
	// TODO: 価格データをパースし、内部状態を更新する
	// 例: 銘柄コード、現在値などを抽出し、a.stateを更新
}

// handleStatus はステータス通知イベントを処理する
func (a *Agent) handleStatus(data map[string]string) {
	a.logger.Info("[Placeholder] Handling Status Notification", "data", data)
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

		currentPrice, err := a.tradeService.GetPrice(ctx, pos.Symbol)
		if err != nil {
			a.logger.Error("failed to get price for exit check", "symbol", pos.Symbol, "error", err)
			continue
		}
		if currentPrice == 0 {
			a.logger.Warn("skipping exit check because current price is zero", "symbol", pos.Symbol)
			continue
		}

		atr, err := a.getATRForPosition(ctx, pos.Symbol)
		if err != nil {
			a.logger.Error("failed to get ATR for exit check", "symbol", pos.Symbol, "error", err)
			continue
		}

		a.updateTrailingStopState(ctx, pos, currentPrice)

		// 決済条件の判定 (優先順位: ATRベース損切り -> トレーリングストップ -> 固定利確)
		if a.shouldStopLossFixed(pos, currentPrice, atr) {
			a.placeExitOrder(ctx, pos, "STOP_LOSS_FIXED")
		} else if a.shouldStopLossTrailing(pos, currentPrice) {
			a.placeExitOrder(ctx, pos, "STOP_LOSS_TRAILING")
		} else if a.shouldProfitTakeFixed(pos, currentPrice) {
			a.placeExitOrder(ctx, pos, "PROFIT_TAKE_FIXED")
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
		Symbol:    pos.Symbol,
		TradeType: model.TradeTypeSell,
		OrderType: model.OrderTypeMarket,
		Quantity:  pos.Quantity,
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
	for _, s := range signals {
		a.logger.Info("signal detail", "symbol", s.Symbol, "signal", s.Signal)
		symbolStr := fmt.Sprintf("%d", s.Symbol)

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
			currentPrice, err := a.tradeService.GetPrice(ctx, symbolStr)
			if err != nil {
				a.logger.Error("failed to get price for sizing", "symbol", symbolStr, "error", err)
				continue
			}
			if currentPrice == 0 {
				a.logger.Warn("skipping buy signal because current price is zero", "symbol", symbolStr)
				continue
			}

			atrPeriod := a.config.StrategySettings.Swingtrade.ATRPeriod
			stopLossATRMultiplier := a.config.StrategySettings.Swingtrade.StopLossATRMultiplier
			unitSize := float64(a.config.StrategySettings.Swingtrade.UnitSize)

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

			// 両方の条件のうち小さい方を採用
			maxShares := math.Min(maxSharesByATR, maxSharesByBuyingPower)
			
			// unitSizeの倍数に切り捨て
			quantity := math.Floor(maxShares / unitSize) * unitSize

			a.logger.Info("calculated order quantity (ATR-based)",
				"symbol", symbolStr,
				"buying_power", balance.BuyingPower,
				"trade_risk_percentage", a.config.StrategySettings.Swingtrade.TradeRiskPercentage,
				"atr_period", atrPeriod,
				"stop_loss_atr_multiplier", stopLossATRMultiplier,
				"calculated_atr", atr,
				"risk_per_share", riskPerShare,
				"total_risk_amount", totalRiskAmount,
				"calculated_quantity", quantity)

			if quantity <= 0 {
				a.logger.Info("skipping buy signal due to zero calculated quantity", "symbol", symbolStr)
				continue
			}

			// 注文リクエストを作成
			req := &PlaceOrderRequest{
				Symbol:    symbolStr,
				TradeType: model.TradeTypeBuy,
				OrderType: model.OrderTypeMarket,
				Quantity:  int(quantity),
				Price:     0, // 成行注文のため価格は0
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
				Symbol:    symbolStr,
				TradeType: model.TradeTypeSell,
				OrderType: model.OrderTypeMarket,
				Quantity:  position.Quantity, // 保有する全数量を売却
				Price:     0,                 // 成行注文のため価格は0
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




