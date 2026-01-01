package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stock-bot/domain/model"
	"stock-bot/internal/config"
	"stock-bot/internal/infrastructure/adapter"
	"stock-bot/internal/infrastructure/client"
	"syscall"
	"time"

	_ "stock-bot/internal/logger" // loggerパッケージをインポートし、slog.Default()を初期化
)

func main() {
	// 1. コマンドラインフラグの定義
	configPath := flag.String("config", "agent_config.yaml", "Path to agent configuration file")
	strategyAPIURL := flag.String("strategy-api", "http://localhost:8080", "Strategy Management API URL")
	tradeAPIURL := flag.String("trade-api", "http://localhost:8080", "Trade API URL")
	agentType := flag.String("type", "swing", "Agent type (swing, day, scalp)")
	strategyID := flag.String("strategy-id", "", "Strategy ID to execute (required)")
	flag.Parse()

	if *strategyID == "" {
		slog.Default().Error("strategy-id is required")
		os.Exit(1)
	}

	// 2. 設定ファイルの読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Default().Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. HTTPクライアントの初期化
	strategyClient := adapter.NewHTTPStrategyClient(*strategyAPIURL, slog.Default())
	tradeClient := adapter.NewHTTPTradeService(*tradeAPIURL)

	// 4. 戦略情報の取得
	ctx := context.Background()
	strategy, err := strategyClient.GetStrategy(ctx, *strategyID)
	if err != nil {
		slog.Default().Error("failed to get strategy",
			slog.String("strategy_id", *strategyID),
			slog.Any("error", err))
		os.Exit(1)
	}

	slog.Default().Info("strategy loaded",
		slog.String("id", strategy.ID),
		slog.String("name", strategy.Name),
		slog.String("type", string(strategy.Type)),
		slog.String("status", string(strategy.Status)))

	// 5. 戦略がアクティブでない場合は終了
	if !strategy.IsActive() {
		slog.Default().Warn("strategy is not active",
			slog.String("status", string(strategy.Status)))
		os.Exit(0)
	}

	// 6. エージェントの初期化（簡略化版）
	logger := slog.Default()

	// 基本的なクライアント初期化（実際の実装では戦略タイプに応じて分岐）
	tachibanaClient := client.NewTachibanaClient(cfg)
	eventClient := client.NewEventClient(logger)

	// 統合クライアントの初期化
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		cfg.TachibanaUserID,
		cfg.TachibanaPassword,
		cfg.TachibanaSecondPassword,
		logger,
	)

	// ログイン処理
	logger.Info("logging in to Tachibana API...")
	session, err := unifiedClient.GetSession(ctx)
	if err != nil {
		logger.Error("failed to login", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("login successful")

	// 7. エージェントの実行ループ
	logger.Info("starting strategy execution agent",
		slog.String("agent_type", *agentType),
		slog.String("strategy_id", *strategyID))

	// シグナル処理の設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 実行ループ
	ticker := time.NewTicker(time.Duration(strategy.Config.ExecutionInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 戦略実行
			if err := executeStrategy(ctx, strategy, strategyClient, unifiedClient, logger); err != nil {
				logger.Error("strategy execution failed",
					slog.String("strategy_id", *strategyID),
					slog.Any("error", err))
			}

		case sig := <-sigChan:
			logger.Info("received signal, shutting down", slog.String("signal", sig.String()))

			// 戦略を非アクティブ化
			if err := strategyClient.DeactivateStrategy(ctx, *strategyID); err != nil {
				logger.Error("failed to deactivate strategy", slog.Any("error", err))
			}

			return
		}
	}
}

// executeStrategy は戦略を実行する（簡略化版）
func executeStrategy(
	ctx context.Context,
	strategy *model.Strategy,
	strategyClient *adapter.HTTPStrategyClient,
	unifiedClient *client.TachibanaUnifiedClient,
	logger *slog.Logger,
) error {
	logger.Debug("executing strategy",
		slog.String("strategy_id", strategy.ID),
		slog.String("strategy_name", strategy.Name))

	// 戦略タイプに応じた実行ロジック
	switch strategy.Type {
	case model.StrategyTypeSwing:
		return executeSwingStrategy(ctx, strategy, strategyClient, unifiedClient, logger)
	case model.StrategyTypeDay:
		return executeDayStrategy(ctx, strategy, strategyClient, unifiedClient, logger)
	case model.StrategyTypeScalp:
		return executeScalpStrategy(ctx, strategy, strategyClient, unifiedClient, logger)
	default:
		return fmt.Errorf("unsupported strategy type: %s", strategy.Type)
	}
}

// executeSwingStrategy はスイング戦略を実行
func executeSwingStrategy(
	ctx context.Context,
	strategy *model.Strategy,
	strategyClient *adapter.HTTPStrategyClient,
	unifiedClient *client.TachibanaUnifiedClient,
	logger *slog.Logger,
) error {
	logger.Debug("executing swing strategy", slog.String("strategy_id", strategy.ID))

	// TODO: 実際のスイング戦略ロジックを実装
	// 1. シグナルファイルの読み込み
	// 2. エントリー/エグジット判断
	// 3. 注文の発行
	// 4. 統計情報の更新

	// プレースホルダー: 統計情報の更新
	if err := strategyClient.UpdateStatistics(ctx, strategy.ID, 0.0, true); err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}

	return nil
}

// executeDayStrategy はデイトレード戦略を実行
func executeDayStrategy(
	ctx context.Context,
	strategy *model.Strategy,
	strategyClient *adapter.HTTPStrategyClient,
	unifiedClient *client.TachibanaUnifiedClient,
	logger *slog.Logger,
) error {
	logger.Debug("executing day strategy", slog.String("strategy_id", strategy.ID))

	// TODO: 実際のデイトレード戦略ロジックを実装
	return nil
}

// executeScalpStrategy はスキャルピング戦略を実行
func executeScalpStrategy(
	ctx context.Context,
	strategy *model.Strategy,
	strategyClient *adapter.HTTPStrategyClient,
	unifiedClient *client.TachibanaUnifiedClient,
	logger *slog.Logger,
) error {
	logger.Debug("executing scalp strategy", slog.String("strategy_id", strategy.ID))

	// TODO: 実際のスキャルピング戦略ロジックを実装
	return nil
}
