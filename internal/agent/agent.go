package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"stock-bot/internal/agent/state"
	"stock-bot/internal/app"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/adapter"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// Agent はイベント処理が分離されたエージェント
type Agent struct {
	configPath       string
	config           *AgentConfig
	logger           *slog.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	signalPattern    string
	state            *state.State
	tradeService     service.TradeService
	positionRepo     repository.PositionRepository
	executionUseCase app.ExecutionUseCase
	gyouNoToSymbol   map[string]string // 行番号(文字列)から銘柄コードへのマッピング

	// イベント処理関連
	webSocketEventService *eventprocessing.WebSocketEventService
	eventDispatcher       service.EventDispatcher
}

// NewAgent は新しいエージェントのインスタンスを作成する
func NewAgent(
	configPath string,
	tradeService service.TradeService,
	eventClient client.EventClient,
	positionRepo repository.PositionRepository,
	executionUseCase app.ExecutionUseCase,
) (*Agent, error) {
	// 先に設定を読み込んでおく
	cfg, err := LoadAgentConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := context.WithCancel(context.Background())

	gyouNoToSymbol := make(map[string]string)
	for i, symbol := range cfg.StrategySettings.Swingtrade.TargetSymbols {
		gyouNoToSymbol[fmt.Sprintf("%d", i+1)] = symbol // 1-indexed
	}

	agentState := state.NewState()

	// イベントディスパッチャーを作成
	eventDispatcher := eventprocessing.NewEventDispatcher(logger)

	// イベントハンドラーを作成・登録
	executionHandler := eventprocessing.NewExecutionEventHandler(executionUseCase, logger)
	priceHandler := eventprocessing.NewPriceEventHandler(agentState, gyouNoToSymbol, logger)
	statusHandler := eventprocessing.NewStatusEventHandler(logger)

	eventDispatcher.RegisterHandler("EC", executionHandler)
	eventDispatcher.RegisterHandler("FD", priceHandler)
	eventDispatcher.RegisterHandler("ST", statusHandler)

	// WebSocketイベントサービスを作成
	webSocketEventService := eventprocessing.NewWebSocketEventService(
		eventClient,
		eventDispatcher,
		logger,
	)

	return &Agent{
		configPath:            configPath,
		config:                cfg,
		logger:                logger,
		ctx:                   ctx,
		cancel:                cancel,
		signalPattern:         cfg.StrategySettings.Swingtrade.SignalFilePattern,
		state:                 agentState,
		tradeService:          tradeService,
		positionRepo:          positionRepo,
		executionUseCase:      executionUseCase,
		gyouNoToSymbol:        gyouNoToSymbol,
		webSocketEventService: webSocketEventService,
		eventDispatcher:       eventDispatcher,
	}, nil
}

// Start はエージェントの実行ループを開始する
func (a *Agent) Start() {
	a.logger.Info("starting agent...")

	// WebSocketイベント監視を開始
	domainSession := a.tradeService.GetSession()
	if domainSession == nil {
		a.logger.Error("failed to get session for event watcher")
		return
	}

	// ドメインセッションをクライアントセッションに変換
	sessionAdapter := adapter.NewSessionAdapter()
	clientSession := sessionAdapter.ToClientSession(domainSession)

	targetSymbols := a.config.StrategySettings.Swingtrade.TargetSymbols
	if err := a.webSocketEventService.StartEventWatcher(a.ctx, clientSession, targetSymbols); err != nil {
		a.logger.Error("failed to start event watcher", "error", err)
		return
	}

	// 定期実行のTicker
	ticker := time.NewTicker(a.config.Agent.ExecutionInterval)
	defer ticker.Stop()

	// 起動時に初期状態を同期
	a.syncInitialState()

	// 起動時に一度実行
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

// Stop はエージェントの実行ループを停止する
func (a *Agent) Stop() {
	a.logger.Info("sending stop signal to agent...")
	a.cancel()
}

// Tick はエージェントの単一の評価サイクルを実行します。バックテスト用に公開されています。
func (a *Agent) Tick() {
	a.tick()
}

// syncInitialState はエージェント起動時にトレードサービスから状態を取得し、内部状態を同期する
func (a *Agent) syncInitialState() {
	a.logger.Info("synchronizing initial state...")
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	// 残高の同期
	balance, err := a.tradeService.GetBalance(ctx)
	if err != nil {
		a.logger.Error("failed to get initial balance", "error", err)
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

	for _, position := range positions {
		a.logger.Info("evaluating position for exit",
			"symbol", position.Symbol,
			"quantity", position.Quantity,
			"average_price", position.AveragePrice)

		// 利確・損切り判定ロジック（簡易実装）
		// 実際の実装では、現在価格を取得して利確・損切り条件をチェック
		// ここでは基本的な構造のみ実装

		// 現在価格の取得（実装例）
		// currentPrice := a.getCurrentPrice(position.Symbol)
		// profitLossPercent := (currentPrice - position.AveragePrice) / position.AveragePrice * 100

		// 利確条件（例：+5%）
		// if profitLossPercent >= 5.0 {
		//     a.placeSellOrder(ctx, position)
		// }

		// 損切り条件（例：-3%）
		// if profitLossPercent <= -3.0 {
		//     a.placeSellOrder(ctx, position)
		// }

		a.logger.Debug("position evaluation completed", "symbol", position.Symbol)
	}
}

// checkSignalsForEntry はシグナルファイルをチェックし、新規エントリー注文を行う
func (a *Agent) checkSignalsForEntry(ctx context.Context) {
	a.logger.Info("checking signals for entry...")

	// 状態の確認（ログ出力）
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

	// シグナルファイルの検索
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

	// シグナルファイルの読み込み
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

	// シグナル処理
	for _, signal := range signals {
		// 銘柄コードを文字列に変換
		symbolStr := fmt.Sprintf("%d", signal.Symbol)

		// 対象銘柄かチェック
		if _, exists := targetSymbolsMap[symbolStr]; !exists {
			a.logger.Debug("signal for non-target symbol, skipping", "symbol", symbolStr)
			continue
		}

		// 既存ポジション・注文のチェック
		if a.hasExistingPositionOrOrder(symbolStr) {
			a.logger.Info("already have position or order for symbol, skipping", "symbol", symbolStr)
			continue
		}

		// 注文発行の判定・実行
		a.processEntrySignal(ctx, signal)
	}
}

// hasExistingPositionOrOrder は指定銘柄の既存ポジションまたは注文があるかチェック
func (a *Agent) hasExistingPositionOrOrder(symbol string) bool {
	// ポジションチェック
	positions := a.state.GetPositions()
	for _, p := range positions {
		if p.Symbol == symbol && p.Quantity > 0 {
			return true
		}
	}

	// 注文チェック
	orders := a.state.GetOrders()
	for _, o := range orders {
		if o.Symbol == symbol && o.IsUnexecuted() {
			return true
		}
	}

	return false
}

// processEntrySignal はエントリーシグナルを処理し、注文を発行する
func (a *Agent) processEntrySignal(ctx context.Context, signal *SignalRecord) {
	// 銘柄コードを文字列に変換
	symbolStr := fmt.Sprintf("%d", signal.Symbol)

	a.logger.Info("processing entry signal",
		"symbol", symbolStr,
		"signal_type", signal.Signal)

	// 注文数量の計算（簡易実装）
	balance := a.state.GetBalance()
	// 仮の価格（実際の実装では現在価格を取得）
	estimatedPrice := 1000.0                                 // プレースホルダー価格
	maxOrderAmount := balance.BuyingPower * 0.1              // 買付余力の10%を上限
	quantity := int(maxOrderAmount/estimatedPrice/100) * 100 // 100株単位

	if quantity < 100 {
		a.logger.Warn("insufficient buying power for minimum order",
			"symbol", symbolStr,
			"required", estimatedPrice*100,
			"available", balance.BuyingPower)
		return
	}

	// 注文発行（簡易実装 - 成行注文）
	req := &service.PlaceOrderRequest{
		Symbol:              symbolStr,
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeMarket,
		Quantity:            quantity,
		Price:               0, // 成行注文
		PositionAccountType: model.PositionAccountTypeCash,
	}

	order, err := a.tradeService.PlaceOrder(ctx, req)
	if err != nil {
		a.logger.Error("failed to place entry order",
			"symbol", symbolStr,
			"error", err)
		return
	}

	a.logger.Info("entry order placed successfully",
		"order_id", order.OrderID,
		"symbol", order.Symbol,
		"quantity", order.Quantity)
}
