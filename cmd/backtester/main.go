package main

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/agent"
	"stock-bot/internal/app" // Add this import
	"time"
)

// mainはバックテスト実行のエントリポイントです。
func main() {
	logger := slog.Default() // デフォルトロガーを使用

	// 1. バックテスト設定
	dataDir := "./data/history"
	configPath := "./agent_config.yaml" // エージェントのコンフィグファイル
	initialCash := 10_000_000.0         // 初期資金 1000万円
	targetSymbols := []string{"7203", "6758", "9984"} // バックテスト対象銘柄

	// 2. BacktestTradeServiceの初期化
	backtestService := agent.NewBacktestTradeService(dataDir, initialCash)
	err := backtestService.LoadHistory(targetSymbols)
	if err != nil {
		logger.Error("failed to load history data", "error", err)
		os.Exit(1)
	}

	// 3. Agentの初期化
	// バックテストでは実際のDBは使用しないため、ダミーのPositionRepositoryとExecutionUseCaseを渡す
	dummyPositionRepo := &dummyPositionRepository{}
	dummyExecutionUseCase := &dummyExecutionUseCase{}
	stockAgent, err := agent.NewAgent(configPath, backtestService, nil, dummyPositionRepo, dummyExecutionUseCase)
	if err != nil {
		logger.Error("failed to create agent for backtest", "error", err)
		os.Exit(1)
	}
	stockAgent.SetLogger(logger) // バックテスト用にロガーを設定
	stockAgent.SyncInitialState() // バックテスト開始前に状態を同期

	// 4. バックテストのメインループ
	// 全対象銘柄のすべての日付を収集し、ユニークでソートされた日付リストを作成する
	dateSet := make(map[time.Time]struct{})
	for _, symbol := range targetSymbols {
		if history, ok := backtestService.AllHistory[symbol]; ok {
			for _, h := range history {
				dateSet[h.Date] = struct{}{}
			}
		}
	}
	
	var dates []time.Time
	for date := range dateSet {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	if len(dates) == 0 {
		logger.Error("no history data loaded or invalid date range")
		os.Exit(1)
	}
	
	logger.Info("Starting backtest...", "startDate", dates[0].Format("2006-01-02"), "endDate", dates[len(dates)-1].Format("2006-01-02"))

	for _, currentDate := range dates {
		backtestService.SetCurrentTick(currentDate)
		stockAgent.Tick()
		
		logger.Info("Backtest tick completed", "date", currentDate.Format("2006-01-02"), "balance", backtestService.Balance.Cash, "positions", len(backtestService.Positions))
	}


	logger.Info("Backtest finished.")
	logger.Info("Final Balance", "cash", backtestService.Balance.Cash, "buying_power", backtestService.Balance.BuyingPower)
	logger.Info("Final Positions")
	for symbol, pos := range backtestService.Positions {
		logger.Info("  Position", "symbol", symbol, "quantity", pos.Quantity, "avgPrice", pos.AveragePrice, "highestPrice", pos.HighestPrice)
	}
	logger.Info("Final Orders", "count", len(backtestService.Orders))
}

// dummyPositionRepository はバックテスト時に使用するダミーのPositionRepositoryです。
type dummyPositionRepository struct{}

// Ensure dummyPositionRepository implements repository.PositionRepository
var _ repository.PositionRepository = (*dummyPositionRepository)(nil)

func (d *dummyPositionRepository) Save(ctx context.Context, position *model.Position) error { return nil }
func (d *dummyPositionRepository) FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) { return nil, nil }
func (d *dummyPositionRepository) FindAll(ctx context.Context) ([]*model.Position, error) { return nil, nil }
func (d *dummyPositionRepository) UpdateHighestPrice(ctx context.Context, symbol string, price float64) error { return nil }
func (d *dummyPositionRepository) UpsertPositionByExecution(ctx context.Context, execution *model.Execution) error { return nil }
func (d *dummyPositionRepository) DeletePosition(ctx context.Context, symbol string) error { return nil }

// dummyExecutionUseCase はバックテスト時に使用するダミーのExecutionUseCaseです。
type dummyExecutionUseCase struct{}

// Ensure dummyExecutionUseCase implements app.ExecutionUseCase
var _ app.ExecutionUseCase = (*dummyExecutionUseCase)(nil)

func (d *dummyExecutionUseCase) Execute(ctx context.Context, execution *model.Execution) error {
	return nil // No-op for backtesting
}
