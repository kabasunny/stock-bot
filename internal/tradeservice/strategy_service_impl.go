package tradeservice

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/event"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"time"

	"github.com/google/uuid"
)

// StrategyServiceImpl はStrategyServiceの実装
type StrategyServiceImpl struct {
	strategyRepo   repository.StrategyRepository
	eventPublisher event.EventPublisher
	logger         *slog.Logger
}

// NewStrategyService は新しいStrategyServiceを作成
func NewStrategyService(
	strategyRepo repository.StrategyRepository,
	eventPublisher event.EventPublisher,
	logger *slog.Logger,
) service.StrategyService {
	return &StrategyServiceImpl{
		strategyRepo:   strategyRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateStrategy は新しい戦略を作成
func (s *StrategyServiceImpl) CreateStrategy(ctx context.Context, req *service.CreateStrategyRequest) (*model.Strategy, error) {
	// 戦略IDを生成
	strategyID := uuid.New().String()

	// デフォルト統計を初期化
	stats := model.StrategyStats{
		ExecutionCount:   0,
		TotalOrders:      0,
		SuccessfulOrders: 0,
		FailedOrders:     0,
		TotalPL:          0.0,
		RealizedPL:       0.0,
		UnrealizedPL:     0.0,
		WinCount:         0,
		LossCount:        0,
		WinRate:          0.0,
		MaxDrawdown:      0.0,
		CurrentDrawdown:  0.0,
		LastExecutedAt:   time.Time{},
	}

	// 戦略を作成
	strategy := &model.Strategy{
		ID:          strategyID,
		Name:        req.Name,
		Type:        req.Type,
		Status:      model.StrategyStatusInactive, // 初期状態は非アクティブ
		Description: req.Description,
		Config:      req.Config,
		RiskLimits:  req.RiskLimits,
		Statistics:  stats,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   req.CreatedBy,
	}

	// 戦略の検証
	if err := s.ValidateStrategy(ctx, strategy); err != nil {
		return nil, fmt.Errorf("strategy validation failed: %w", err)
	}

	s.logger.Info("strategy created",
		slog.String("id", strategy.ID),
		slog.String("name", strategy.Name),
		slog.String("type", string(strategy.Type)))

	return strategy, nil
}

// GetStrategy は戦略を取得
func (s *StrategyServiceImpl) GetStrategy(ctx context.Context, id string) (*model.Strategy, error) {
	strategy, err := s.strategyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy: %w", err)
	}
	if strategy == nil {
		return nil, fmt.Errorf("strategy not found: %s", id)
	}
	return strategy, nil
}

// GetActiveStrategies はアクティブな戦略を取得
func (s *StrategyServiceImpl) GetActiveStrategies(ctx context.Context) ([]*model.Strategy, error) {
	strategies, err := s.strategyRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active strategies: %w", err)
	}
	return strategies, nil
}

// UpdateStrategy は戦略を更新
func (s *StrategyServiceImpl) UpdateStrategy(ctx context.Context, strategy *model.Strategy) error {
	// 戦略の検証
	if err := s.ValidateStrategy(ctx, strategy); err != nil {
		return fmt.Errorf("strategy validation failed: %w", err)
	}

	strategy.UpdatedAt = time.Now()

	if err := s.strategyRepo.Update(ctx, strategy); err != nil {
		return fmt.Errorf("failed to update strategy: %w", err)
	}

	s.logger.Info("strategy updated",
		slog.String("id", strategy.ID),
		slog.String("name", strategy.Name))

	return nil
}

// ActivateStrategy は戦略をアクティブ化
func (s *StrategyServiceImpl) ActivateStrategy(ctx context.Context, id string) error {
	strategy, err := s.GetStrategy(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get strategy: %w", err)
	}

	// リスク制限チェック
	violations, err := s.CheckRiskLimits(ctx, strategy)
	if err != nil {
		return fmt.Errorf("failed to check risk limits: %w", err)
	}
	if len(violations) > 0 {
		return fmt.Errorf("cannot activate strategy due to risk limit violations: %v", violations)
	}

	strategy.Activate()

	if err := s.strategyRepo.Update(ctx, strategy); err != nil {
		return fmt.Errorf("failed to activate strategy: %w", err)
	}

	// ドメインイベントを発行
	if s.eventPublisher != nil {
		strategyEvent := event.NewStrategyActivatedEvent(strategy)
		if err := s.eventPublisher.Publish(ctx, strategyEvent); err != nil {
			s.logger.Warn("failed to publish strategy activated event", slog.Any("error", err))
		}
	}

	s.logger.Info("strategy activated", slog.String("id", id))
	return nil
}

// DeactivateStrategy は戦略を非アクティブ化
func (s *StrategyServiceImpl) DeactivateStrategy(ctx context.Context, id string) error {
	strategy, err := s.GetStrategy(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get strategy: %w", err)
	}

	strategy.Deactivate()

	if err := s.strategyRepo.Update(ctx, strategy); err != nil {
		return fmt.Errorf("failed to deactivate strategy: %w", err)
	}

	// ドメインイベントを発行
	if s.eventPublisher != nil {
		strategyEvent := event.NewStrategyDeactivatedEvent(strategy)
		if err := s.eventPublisher.Publish(ctx, strategyEvent); err != nil {
			s.logger.Warn("failed to publish strategy deactivated event", slog.Any("error", err))
		}
	}

	s.logger.Info("strategy deactivated", slog.String("id", id))
	return nil
}

// PauseStrategy は戦略を一時停止
func (s *StrategyServiceImpl) PauseStrategy(ctx context.Context, id string) error {
	strategy, err := s.GetStrategy(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get strategy: %w", err)
	}

	strategy.Pause()

	if err := s.strategyRepo.Update(ctx, strategy); err != nil {
		return fmt.Errorf("failed to pause strategy: %w", err)
	}

	s.logger.Info("strategy paused", slog.String("id", id))
	return nil
}

// ValidateStrategy は戦略の設定を検証
func (s *StrategyServiceImpl) ValidateStrategy(ctx context.Context, strategy *model.Strategy) error {
	// 基本的な検証
	if strategy.Name == "" {
		return fmt.Errorf("strategy name is required")
	}
	if strategy.Type == "" {
		return fmt.Errorf("strategy type is required")
	}

	// 設定の検証
	config := strategy.Config
	if len(config.TargetSymbols) == 0 {
		return fmt.Errorf("at least one target symbol is required")
	}
	if config.ExecutionInterval <= 0 {
		return fmt.Errorf("execution interval must be positive")
	}

	// リスク制限の検証
	limits := strategy.RiskLimits
	if limits.MaxLossAmount < 0 {
		return fmt.Errorf("max loss amount cannot be negative")
	}
	if limits.MaxLossPercent < 0 || limits.MaxLossPercent > 100 {
		return fmt.Errorf("max loss percent must be between 0 and 100")
	}
	if limits.MaxPositions < 0 {
		return fmt.Errorf("max positions cannot be negative")
	}

	return nil
}

// CheckRiskLimits は戦略のリスク制限をチェック
func (s *StrategyServiceImpl) CheckRiskLimits(ctx context.Context, strategy *model.Strategy) ([]string, error) {
	violations := strategy.CheckRiskLimits()

	// 追加のリスクチェックロジックをここに実装
	// 例：市場状況、ボラティリティ、相関リスクなど

	if len(violations) > 0 {
		s.logger.Warn("risk limit violations detected",
			slog.String("strategy_id", strategy.ID),
			slog.Any("violations", violations))

		// リスク制限超過イベントを発行
		if s.eventPublisher != nil {
			for _, violation := range violations {
				riskEvent := event.NewRiskLimitExceededEvent(
					violation,
					strategy.Statistics.TotalPL,
					strategy.RiskLimits.MaxLossAmount,
					strategy.ID,
				)
				if err := s.eventPublisher.Publish(ctx, riskEvent); err != nil {
					s.logger.Warn("failed to publish risk limit exceeded event", slog.Any("error", err))
				}
			}
		}
	}

	return violations, nil
}

// UpdateStatistics は戦略の統計を更新
func (s *StrategyServiceImpl) UpdateStatistics(ctx context.Context, strategyID string, pl float64, isWin bool) error {
	strategy, err := s.GetStrategy(ctx, strategyID)
	if err != nil {
		return fmt.Errorf("failed to get strategy: %w", err)
	}

	strategy.UpdateStatistics(pl, isWin)

	if err := s.strategyRepo.UpdateStatistics(ctx, strategyID, &strategy.Statistics); err != nil {
		return fmt.Errorf("failed to update strategy statistics: %w", err)
	}

	s.logger.Debug("strategy statistics updated",
		slog.String("strategy_id", strategyID),
		slog.Float64("pl", pl),
		slog.Bool("is_win", isWin))

	return nil
}
