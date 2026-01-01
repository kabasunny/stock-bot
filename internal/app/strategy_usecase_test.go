package app

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/domain/event"
	"stock-bot/domain/model"
	repository_impl "stock-bot/internal/infrastructure/repository"
	"stock-bot/internal/tradeservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStrategyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&repository_impl.StrategyEntity{})
	require.NoError(t, err)

	return db
}

func TestStrategyUseCaseImpl_CreateStrategy(t *testing.T) {
	db := setupStrategyTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)
	unitOfWork := repository_impl.NewUnitOfWork(db, eventPublisher, logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := tradeservice.NewStrategyService(strategyRepo, eventPublisher, logger)
	strategyUseCase := NewStrategyUseCaseImpl(strategyRepo, strategyService, unitOfWork, logger)

	ctx := context.Background()
	strategy := &model.Strategy{
		ID:          "usecase-test-strategy-1",
		Name:        "UseCase Test Strategy",
		Type:        model.StrategyTypeSwing,
		Status:      model.StrategyStatusActive,
		Description: "Test strategy for use case testing",
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create strategy through use case
	err := strategyUseCase.CreateStrategy(ctx, strategy)
	assert.NoError(t, err)

	// Verify strategy was created
	retrieved, err := strategyUseCase.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, strategy.ID, retrieved.ID)
	assert.Equal(t, strategy.Name, retrieved.Name)
	assert.Equal(t, strategy.Type, retrieved.Type)
}

func TestStrategyUseCaseImpl_UpdateStrategy(t *testing.T) {
	db := setupStrategyTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)
	unitOfWork := repository_impl.NewUnitOfWork(db, eventPublisher, logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := tradeservice.NewStrategyService(strategyRepo, eventPublisher, logger)
	strategyUseCase := NewStrategyUseCaseImpl(strategyRepo, strategyService, unitOfWork, logger)

	ctx := context.Background()
	strategy := &model.Strategy{
		ID:          "usecase-test-strategy-2",
		Name:        "UseCase Test Strategy 2",
		Type:        model.StrategyTypeDay,
		Status:      model.StrategyStatusActive,
		Description: "Test strategy for use case update testing",
		Config: model.StrategyConfig{
			TargetSymbols:     []string{"1301"},
			ExecutionInterval: 3 * time.Minute,
		},
		RiskLimits: model.RiskLimits{
			MaxLossAmount: 30000,
			MaxPositions:  3,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create strategy
	err := strategyUseCase.CreateStrategy(ctx, strategy)
	require.NoError(t, err)

	// Update strategy
	strategy.Name = "Updated UseCase Test Strategy 2"
	strategy.RiskLimits.MaxPositions = 10
	strategy.RiskLimits.MaxLossAmount = 40000

	err = strategyUseCase.UpdateStrategy(ctx, strategy)
	assert.NoError(t, err)

	// Verify strategy was updated
	retrieved, err := strategyUseCase.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated UseCase Test Strategy 2", retrieved.Name)
	assert.Equal(t, 10, retrieved.RiskLimits.MaxPositions)
	assert.Equal(t, 40000.0, retrieved.RiskLimits.MaxLossAmount)
}

func TestStrategyUseCaseImpl_ListStrategies(t *testing.T) {
	db := setupStrategyTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)
	unitOfWork := repository_impl.NewUnitOfWork(db, eventPublisher, logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := tradeservice.NewStrategyService(strategyRepo, eventPublisher, logger)
	strategyUseCase := NewStrategyUseCaseImpl(strategyRepo, strategyService, unitOfWork, logger)

	ctx := context.Background()

	// Create multiple strategies
	strategies := []*model.Strategy{
		{
			ID:          "usecase-strategy-1",
			Name:        "UseCase Strategy 1",
			Type:        model.StrategyTypeSwing,
			Status:      model.StrategyStatusActive,
			Description: "First use case test strategy",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "usecase-strategy-2",
			Name:        "UseCase Strategy 2",
			Type:        model.StrategyTypeDay,
			Status:      model.StrategyStatusInactive,
			Description: "Second use case test strategy",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "usecase-strategy-3",
			Name:        "UseCase Strategy 3",
			Type:        model.StrategyTypeScalp,
			Status:      model.StrategyStatusActive,
			Description: "Third use case test strategy",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, strategy := range strategies {
		err := strategyUseCase.CreateStrategy(ctx, strategy)
		require.NoError(t, err)
	}

	// List all strategies
	allStrategies, err := strategyUseCase.ListStrategies(ctx)
	assert.NoError(t, err)
	assert.Len(t, allStrategies, 3)
}

func TestStrategyUseCaseImpl_DeleteStrategy(t *testing.T) {
	db := setupStrategyTestDB(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	eventPublisher := event.NewInMemoryEventPublisher(logger)
	unitOfWork := repository_impl.NewUnitOfWork(db, eventPublisher, logger)

	strategyRepo := repository_impl.NewStrategyRepository(db)
	strategyService := tradeservice.NewStrategyService(strategyRepo, eventPublisher, logger)
	strategyUseCase := NewStrategyUseCaseImpl(strategyRepo, strategyService, unitOfWork, logger)

	ctx := context.Background()
	strategy := &model.Strategy{
		ID:          "usecase-test-strategy-delete",
		Name:        "UseCase Strategy to Delete",
		Type:        model.StrategyTypeDay,
		Status:      model.StrategyStatusActive,
		Description: "This use case strategy will be deleted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create strategy
	err := strategyUseCase.CreateStrategy(ctx, strategy)
	require.NoError(t, err)

	// Verify strategy exists
	_, err = strategyUseCase.GetStrategy(ctx, strategy.ID)
	assert.NoError(t, err)

	// Delete strategy
	err = strategyUseCase.DeleteStrategy(ctx, strategy.ID)
	assert.NoError(t, err)

	// Verify strategy was deleted
	_, err = strategyUseCase.GetStrategy(ctx, strategy.ID)
	assert.Error(t, err)
}
