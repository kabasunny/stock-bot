package tradeservice

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/domain/event"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	repository_impl "stock-bot/internal/infrastructure/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&repository_impl.StrategyEntity{})
	require.NoError(t, err)

	return db
}

func TestStrategyService_CreateStrategy(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()
	req := &service.CreateStrategyRequest{
		Name:        "Test Strategy",
		Type:        model.StrategyTypeSwing,
		Description: "Test strategy for unit testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301", "1332"},
			ExecutionInterval: 5 * time.Minute,
		},
		RiskLimits: model.RiskLimits{
			MaxLossAmount:         50000,
			MaxPositions:          5,
			MaxPositionsPerSymbol: 2,
			MaxLeverage:           2.0,
		},
		CreatedBy: "test-user",
	}

	// Create strategy
	strategy, err := strategyService.CreateStrategy(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, strategy)
	assert.NotEmpty(t, strategy.ID)
	assert.Equal(t, req.Name, strategy.Name)
	assert.Equal(t, req.Type, strategy.Type)
	assert.Equal(t, model.StrategyStatusInactive, strategy.Status)
	assert.Equal(t, len(req.Config.TargetSymbols), len(strategy.Config.TargetSymbols))
}

func TestStrategyService_GetStrategy(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()

	// First create a strategy
	req := &service.CreateStrategyRequest{
		Name:        "Test Strategy for Get",
		Type:        model.StrategyTypeDay,
		Description: "Test strategy for get testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301"},
			ExecutionInterval: 3 * time.Minute,
		},
		RiskLimits: model.RiskLimits{
			MaxLossAmount: 30000,
			MaxPositions:  3,
		},
		CreatedBy: "test-user",
	}

	created, err := strategyService.CreateStrategy(ctx, req)
	require.NoError(t, err)

	// Save to repository
	err = strategyRepo.Save(ctx, created)
	require.NoError(t, err)

	// Now get the strategy
	retrieved, err := strategyService.GetStrategy(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Type, retrieved.Type)
}

func TestStrategyService_UpdateStrategy(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()

	// Create and save a strategy
	req := &service.CreateStrategyRequest{
		Name:        "Test Strategy for Update",
		Type:        model.StrategyTypeSwing,
		Description: "Test strategy for update testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301"},
			ExecutionInterval: 5 * time.Minute,
		},
		RiskLimits: model.RiskLimits{
			MaxLossAmount: 40000,
			MaxPositions:  4,
		},
		CreatedBy: "test-user",
	}

	strategy, err := strategyService.CreateStrategy(ctx, req)
	require.NoError(t, err)

	err = strategyRepo.Save(ctx, strategy)
	require.NoError(t, err)

	// Update strategy
	strategy.Name = "Updated Test Strategy"
	strategy.RiskLimits.MaxPositions = 10
	strategy.RiskLimits.MaxLossAmount = 60000

	err = strategyService.UpdateStrategy(ctx, strategy)
	assert.NoError(t, err)

	// Verify strategy was updated
	retrieved, err := strategyService.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Strategy", retrieved.Name)
	assert.Equal(t, 10, retrieved.RiskLimits.MaxPositions)
	assert.Equal(t, 60000.0, retrieved.RiskLimits.MaxLossAmount)
}

func TestStrategyService_ActivateStrategy(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()

	// Create and save a strategy
	req := &service.CreateStrategyRequest{
		Name:        "Test Strategy for Activation",
		Type:        model.StrategyTypeDay,
		Description: "Test strategy for activation testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301"},
			ExecutionInterval: 3 * time.Minute,
		},
		RiskLimits: model.RiskLimits{
			MaxLossAmount: 30000,
			MaxPositions:  3,
		},
		CreatedBy: "test-user",
	}

	strategy, err := strategyService.CreateStrategy(ctx, req)
	require.NoError(t, err)

	err = strategyRepo.Save(ctx, strategy)
	require.NoError(t, err)

	// Activate strategy
	err = strategyService.ActivateStrategy(ctx, strategy.ID)
	assert.NoError(t, err)

	// Verify strategy was activated
	retrieved, err := strategyService.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.StrategyStatusActive, retrieved.Status)
}

func TestStrategyService_GetActiveStrategies(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()

	// Create multiple strategies
	strategies := []*service.CreateStrategyRequest{
		{
			Name:        "Active Strategy 1",
			Type:        model.StrategyTypeSwing,
			Description: "First active test strategy",
			Config: model.StrategyConfig{
				TargetSymbols:     []string{"1301"},
				ExecutionInterval: 5 * time.Minute,
			},
			CreatedBy: "test-user",
		},
		{
			Name:        "Inactive Strategy",
			Type:        model.StrategyTypeDay,
			Description: "Inactive test strategy",
			Config: model.StrategyConfig{
				TargetSymbols:     []string{"1332"},
				ExecutionInterval: 3 * time.Minute,
			},
			CreatedBy: "test-user",
		},
		{
			Name:        "Active Strategy 2",
			Type:        model.StrategyTypeScalp,
			Description: "Second active test strategy",
			Config: model.StrategyConfig{
				TargetSymbols:     []string{"1333"},
				ExecutionInterval: 1 * time.Minute,
			},
			CreatedBy: "test-user",
		},
	}

	var createdStrategies []*model.Strategy
	for _, req := range strategies {
		strategy, err := strategyService.CreateStrategy(ctx, req)
		require.NoError(t, err)

		err = strategyRepo.Save(ctx, strategy)
		require.NoError(t, err)

		createdStrategies = append(createdStrategies, strategy)
	}

	// Activate first and third strategies
	err := strategyService.ActivateStrategy(ctx, createdStrategies[0].ID)
	require.NoError(t, err)
	err = strategyService.ActivateStrategy(ctx, createdStrategies[2].ID)
	require.NoError(t, err)

	// Get active strategies
	activeStrategies, err := strategyService.GetActiveStrategies(ctx)
	assert.NoError(t, err)
	assert.Len(t, activeStrategies, 2)
}

func TestStrategyService_UpdateStatistics(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := NewStrategyService(strategyRepo, eventPublisher, logger)

	ctx := context.Background()

	// Create and save a strategy
	req := &service.CreateStrategyRequest{
		Name:        "Test Strategy for Statistics",
		Type:        model.StrategyTypeSwing,
		Description: "Test strategy for statistics testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301"},
			ExecutionInterval: 5 * time.Minute,
		},
		CreatedBy: "test-user",
	}

	strategy, err := strategyService.CreateStrategy(ctx, req)
	require.NoError(t, err)

	err = strategyRepo.Save(ctx, strategy)
	require.NoError(t, err)

	// Update statistics with a win
	err = strategyService.UpdateStatistics(ctx, strategy.ID, 5000.0, true)
	assert.NoError(t, err)

	// Update statistics with a loss
	err = strategyService.UpdateStatistics(ctx, strategy.ID, -2000.0, false)
	assert.NoError(t, err)

	// Verify statistics were updated
	retrieved, err := strategyService.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, retrieved.Statistics.ExecutionCount)
	assert.Equal(t, 1, retrieved.Statistics.WinCount)
	assert.Equal(t, 1, retrieved.Statistics.LossCount)
	assert.Equal(t, 3000.0, retrieved.Statistics.TotalPL)
	assert.Equal(t, 0.5, retrieved.Statistics.WinRate)
}
