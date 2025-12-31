package app

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"time"
)

// StrategyUseCaseImpl はStrategyUseCaseの実装
type StrategyUseCaseImpl struct {
	strategyRepo    repository.StrategyRepository
	strategyService service.StrategyService
	unitOfWork      repository.UnitOfWork
	logger          *slog.Logger
}

// NewStrategyUseCaseImpl は新しいStrategyUseCaseImplを作成
func NewStrategyUseCaseImpl(
	strategyRepo repository.StrategyRepository,
	strategyService service.StrategyService,
	unitOfWork repository.UnitOfWork,
	logger *slog.Logger,
) *StrategyUseCaseImpl {
	return &StrategyUseCaseImpl{
		strategyRepo:    strategyRepo,
		strategyService: strategyService,
		unitOfWork:      unitOfWork,
		logger:          logger,
	}
}

// CreateStrategy は新しい戦略を作成
func (uc *StrategyUseCaseImpl) CreateStrategy(ctx context.Context, req *service.CreateStrategyRequest) (*model.Strategy, error) {
	uc.logger.Info("creating new strategy", slog.String("name", req.Name), slog.String("type", string(req.Type)))

	strategy, err := uc.strategyService.CreateStrategy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create strategy: %w", err)
	}

	if err := uc.strategyRepo.Save(ctx, strategy); err != nil {
		return nil, fmt.Errorf("failed to save strategy: %w", err)
	}

	uc.logger.Info("strategy created successfully", slog.String("id", strategy.ID), slog.String("name", strategy.Name))
	return strategy, nil
}

// GetStrategy は戦略を取得
func (uc *StrategyUseCaseImpl) GetStrategy(ctx context.Context, id string) (*model.Strategy, error) {
	strategy, err := uc.strategyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy: %w", err)
	}
	return strategy, nil
}

// ListStrategies は戦略一覧を取得
func (uc *StrategyUseCaseImpl) ListStrategies(ctx context.Context, status *model.StrategyStatus) ([]*model.Strategy, error) {
	if status != nil {
		return uc.strategyRepo.FindByStatus(ctx, *status)
	}
	return uc.strategyRepo.FindAll(ctx)
}

// UpdateStrategy は戦略を更新
func (uc *StrategyUseCaseImpl) UpdateStrategy(ctx context.Context, strategy *model.Strategy) error {
	uc.logger.Info("updating strategy", slog.String("id", strategy.ID), slog.String("name", strategy.Name))

	// 戦略の検証
	if err := uc.strategyService.ValidateStrategy(ctx, strategy); err != nil {
		return fmt.Errorf("strategy validation failed: %w", err)
	}

	strategy.UpdatedAt = time.Now()
	if err := uc.strategyRepo.Update(ctx, strategy); err != nil {
		return fmt.Errorf("failed to update strategy: %w", err)
	}

	uc.logger.Info("strategy updated successfully", slog.String("id", strategy.ID))
	return nil
}

// ActivateStrategy は戦略をアクティブ化
func (uc *StrategyUseCaseImpl) ActivateStrategy(ctx context.Context, id string) error {
	uc.logger.Info("activating strategy", slog.String("id", id))

	if err := uc.strategyService.ActivateStrategy(ctx, id); err != nil {
		return fmt.Errorf("failed to activate strategy: %w", err)
	}

	uc.logger.Info("strategy activated successfully", slog.String("id", id))
	return nil
}

// DeactivateStrategy は戦略を非アクティブ化
func (uc *StrategyUseCaseImpl) DeactivateStrategy(ctx context.Context, id string) error {
	uc.logger.Info("deactivating strategy", slog.String("id", id))

	if err := uc.strategyService.DeactivateStrategy(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate strategy: %w", err)
	}

	uc.logger.Info("strategy deactivated successfully", slog.String("id", id))
	return nil
}

// PauseStrategy は戦略を一時停止
func (uc *StrategyUseCaseImpl) PauseStrategy(ctx context.Context, id string) error {
	uc.logger.Info("pausing strategy", slog.String("id", id))

	if err := uc.strategyService.PauseStrategy(ctx, id); err != nil {
		return fmt.Errorf("failed to pause strategy: %w", err)
	}

	uc.logger.Info("strategy paused successfully", slog.String("id", id))
	return nil
}

// DeleteStrategy は戦略を削除
func (uc *StrategyUseCaseImpl) DeleteStrategy(ctx context.Context, id string) error {
	uc.logger.Info("deleting strategy", slog.String("id", id))

	if err := uc.strategyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete strategy: %w", err)
	}

	uc.logger.Info("strategy deleted successfully", slog.String("id", id))
	return nil
}

// ExecuteStrategies はアクティブな戦略を実行
func (uc *StrategyUseCaseImpl) ExecuteStrategies(ctx context.Context) error {
	strategies, err := uc.strategyRepo.FindActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active strategies: %w", err)
	}

	uc.logger.Info("executing active strategies", slog.Int("count", len(strategies)))

	for _, strategy := range strategies {
		if err := uc.executeStrategy(ctx, strategy); err != nil {
			uc.logger.Error("failed to execute strategy",
				slog.String("id", strategy.ID),
				slog.String("name", strategy.Name),
				slog.Any("error", err))
			// 一つの戦略が失敗しても他の戦略は続行
			continue
		}
	}

	return nil
}

// executeStrategy は単一の戦略を実行
func (uc *StrategyUseCaseImpl) executeStrategy(ctx context.Context, strategy *model.Strategy) error {
	uc.logger.Debug("executing strategy", slog.String("id", strategy.ID), slog.String("name", strategy.Name))

	// リスク制限チェック
	violations, err := uc.strategyService.CheckRiskLimits(ctx, strategy)
	if err != nil {
		return fmt.Errorf("failed to check risk limits: %w", err)
	}

	if len(violations) > 0 {
		uc.logger.Warn("strategy has risk limit violations, skipping execution",
			slog.String("id", strategy.ID),
			slog.Any("violations", violations))
		return nil
	}

	// 戦略実行の実装は戦略タイプに応じて分岐
	// 現在は基本的な構造のみ実装
	uc.logger.Info("strategy executed successfully", slog.String("id", strategy.ID))
	return nil
}
