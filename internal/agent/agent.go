package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"stock-bot/domain/model"
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
	state         *State         // <<<<<<<<<<<<<<<< 追加
	tradeService  TradeService // <<<<<<<<<<<<<<<< 追加
}

// NewAgent は新しいエージェントのインスタンスを作成する
// tradeService はトレードサービス（Go APIラッパー）の実装
func NewAgent(configPath string, tradeService TradeService) (*Agent, error) { // <<<<<<<<<<<<<<<< 引数追加
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
		state:         NewState(), // <<<<<<<<<<<<<<<< 追加
		tradeService:  tradeService, // <<<<<<<<<<<<<<<< 追加
	}, nil
}

// Start はエージェントの実行ループを開始する
func (a *Agent) Start() {
	a.logger.Info("starting agent...")
	ticker := time.NewTicker(a.config.Agent.ExecutionInterval)
	defer ticker.Stop()

	// 起動時に初期状態を同期
	a.syncInitialState() // <<<<<<<<<<<<<<<< 追加

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

// Stop はエージェントの実行ループを停止する
func (a *Agent) Stop() {
	a.logger.Info("sending stop signal to agent...")
	a.cancel()
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

	// 状態の確認（ログ出力のみ）
	balance := a.state.GetBalance()
	a.logger.Info("current balance", "cash", balance.Cash, "buying_power", balance.BuyingPower)
	// TODO: ポジションや注文のログ出力も追加

	// TODO: シグナルファイルが複数見つかった場合の処理 (最新のものを一つ選ぶなど)
	// 現状はFindSignalFileが一つだけ返すことを期待
	signalFilePath, err := FindSignalFile(a.signalPattern)
	if err != nil {
		a.logger.Error("failed to find signal file", "error", err)
		return
	}
	if signalFilePath == "" {
		a.logger.Info("no signal file found, skipping this tick")
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

		// 意思決定ロジック
		if s.Signal == BuySignal {
			if _, ok := a.state.GetPosition(symbolStr); ok {
				a.logger.Info("skipping buy signal for already held position", "symbol", symbolStr)
				continue
			}
			a.logger.Info("preparing to place buy order", "symbol", symbolStr)

			// 注文リクエストを作成
			req := &PlaceOrderRequest{
				Symbol:    symbolStr,
				TradeType: model.TradeTypeBuy,
				OrderType: model.OrderTypeMarket,
				Quantity:  a.config.StrategySettings.Swingtrade.LotSize,
				Price:     0, // 成行注文のため価格は0
			}

			// 注文を発行
			// TODO: このコンテキストはループの外でタイムアウト付きで生成した方が良いかもしれない
			ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
			defer cancel()

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
			ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
			defer cancel()

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

// FindSignalFile は指定されたパターンに一致するシグナルファイルを探す
// 今は単純に最初に見つかったファイルを返す
func FindSignalFile(pattern string) (string, error) {
	// globはパターンに一致するファイル名のスライスを返す
	// ファイルは辞書順でソートされる
	files, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", nil // ファイルが見つからなくてもエラーではない
	}
	// TODO: 複数のファイルが見つかった場合に、どれを使うべきか決定するロジックが必要
	// (例: 最も新しいタイムスタンプを持つファイルを選ぶ)
	// 今回は簡単のため、最初に見つかったファイルを返す
	return files[0], nil
}
