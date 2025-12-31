package refactoredagent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"stock-bot/internal/agent"
	"stock-bot/internal/app"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/state"
	"time"
)

// Agent はイベント処理が分離されたエージェント
type Agent struct {
	configPath       string
	config           *agent.AgentConfig
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

// NewAgent は新しいリファクタリング済みエージェントのインスタンスを作成する
func NewAgent(
	configPath string,
	tradeService service.TradeService,
	eventClient client.EventClient,
	positionRepo repository.PositionRepository,
	executionUseCase app.ExecutionUseCase,
) (*Agent, error) {
	// 先に設定を読み込んでおく
	cfg, err := agent.LoadAgentConfig(configPath)
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
	a.logger.Info("starting refactored agent...")

	// WebSocketイベント監視を開始
	session := a.tradeService.GetSession()
	if session == nil {
		a.logger.Error("failed to get session for event watcher")
		return
	}

	targetSymbols := a.config.StrategySettings.Swingtrade.TargetSymbols
	if err := a.webSocketEventService.StartEventWatcher(a.ctx, session, targetSymbols); err != nil {
		a.logger.Error("failed to start event watcher", "error", err)
		return
	}

	// 定期実行のTicker
	ticker := time.NewTicker(a.config.Agent.ExecutionInterval)
	defer ticker.Stop()

	// 起動時に初期状態を同期
	a.syncInitialState()

	// 起動時に一度実行 (初期状態同期後にtickを実行)
	a.tick()

	for {
		select {
		case <-ticker.C:
			a.tick()
		case <-a.ctx.Done():
			a.logger.Info("refactored agent stopping...")
			return
		}
	}
}

// Stop はエージェントの実行ループを停止する
func (a *Agent) Stop() {
	a.logger.Info("sending stop signal to refactored agent...")
	a.cancel()
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

// tick はループごとに実行される処理（戦略実行のみ）
func (a *Agent) tick() {
	a.logger.Info("refactored agent tick")

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
	// TODO: 元のagent.goから移植する
	a.logger.Info("checking positions for exit (placeholder)")
}

// checkSignalsForEntry はシグナルファイルをチェックし、新規エントリー注文を行う
func (a *Agent) checkSignalsForEntry(ctx context.Context) {
	// TODO: 元のagent.goから移植する
	a.logger.Info("checking signals for entry (placeholder)")
}
