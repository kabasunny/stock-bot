package app

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
)

// StrategyFilter は戦略検索のフィルター条件
type StrategyFilter struct {
	Type   *model.StrategyType   `json:"type,omitempty"`
	Status *model.StrategyStatus `json:"status,omitempty"`
	Limit  int                   `json:"limit,omitempty"`
}

// StrategyUseCase は戦略管理のユースケース
type StrategyUseCase interface {
	// CreateStrategy は新しい戦略を作成
	CreateStrategy(ctx context.Context, req *service.CreateStrategyRequest) (*model.Strategy, error)

	// GetStrategy は戦略を取得
	GetStrategy(ctx context.Context, id string) (*model.Strategy, error)

	// ListStrategies は戦略一覧を取得（フィルター対応）
	ListStrategies(ctx context.Context, filter StrategyFilter) ([]*model.Strategy, error)

	// GetActiveStrategies はアクティブな戦略を取得
	GetActiveStrategies(ctx context.Context) ([]*model.Strategy, error)

	// UpdateStrategy は戦略を更新
	UpdateStrategy(ctx context.Context, strategy *model.Strategy) error

	// ActivateStrategy は戦略をアクティブ化
	ActivateStrategy(ctx context.Context, id string) error

	// DeactivateStrategy は戦略を非アクティブ化
	DeactivateStrategy(ctx context.Context, id string) error

	// PauseStrategy は戦略を一時停止
	PauseStrategy(ctx context.Context, id string) error

	// DeleteStrategy は戦略を削除
	DeleteStrategy(ctx context.Context, id string) error

	// UpdateStatistics は戦略の統計を更新
	UpdateStatistics(ctx context.Context, strategyID string, pl float64, isWin bool) error

	// ExecuteStrategies はアクティブな戦略を実行
	ExecuteStrategies(ctx context.Context) error
}
