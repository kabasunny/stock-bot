package repository

import (
	"context"
	"stock-bot/domain/model"
)

// StrategyRepository は戦略のリポジトリインターフェース
type StrategyRepository interface {
	// Save は戦略を保存
	Save(ctx context.Context, strategy *model.Strategy) error

	// FindByID はIDで戦略を検索
	FindByID(ctx context.Context, id string) (*model.Strategy, error)

	// FindByType は種別で戦略を検索
	FindByType(ctx context.Context, strategyType model.StrategyType) ([]*model.Strategy, error)

	// FindByStatus はステータスで戦略を検索
	FindByStatus(ctx context.Context, status model.StrategyStatus) ([]*model.Strategy, error)

	// FindActive はアクティブな戦略を検索
	FindActive(ctx context.Context) ([]*model.Strategy, error)

	// FindAll は全ての戦略を検索
	FindAll(ctx context.Context) ([]*model.Strategy, error)

	// Update は戦略を更新
	Update(ctx context.Context, strategy *model.Strategy) error

	// Delete は戦略を削除
	Delete(ctx context.Context, id string) error

	// UpdateStatistics は戦略の統計を更新
	UpdateStatistics(ctx context.Context, id string, stats *model.StrategyStats) error
}
