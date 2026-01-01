package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"stock-bot/internal/agent"
	"syscall"
)

func main() {
	// コマンドラインフラグ
	baseURL := flag.String("base-url", "http://localhost:8080", "Base URL of the Goa service")
	strategyType := flag.String("strategy", "simple", "Strategy type (simple, swing, day)")
	targetSymbol := flag.String("symbol", "7203", "Target symbol for trading")
	quantity := flag.Uint("quantity", 100, "Order quantity")
	flag.Parse()

	// ロガーのセットアップ
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 戦略の作成
	var strategy agent.Strategy
	switch *strategyType {
	case "swing":
		strategy = agent.NewSwingStrategy("SwingStrategy", *targetSymbol, *quantity, logger)
	case "day":
		strategy = agent.NewDayTradingStrategy("DayTradingStrategy", *targetSymbol, *quantity, logger)
	default:
		strategy = agent.NewSimpleStrategy("SimpleStrategy", *targetSymbol, *quantity, logger)
	}

	// 軽量エージェントの作成
	lightAgent := agent.NewLightweightAgent(*baseURL, strategy, logger)

	// コンテキストとシグナルハンドリング
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// シグナルハンドリング
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal, shutting down", "signal", sig)
		cancel()
	}()

	// エージェント実行
	logger.Info("Starting lightweight agent",
		"base_url", *baseURL,
		"strategy", strategy.Name(),
		"target_symbol", *targetSymbol,
		"quantity", *quantity)

	if err := lightAgent.Run(ctx); err != nil && err != context.Canceled {
		logger.Error("Agent execution failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Agent shutdown complete")
}
