package service

import (
	"context"
	"stock-bot/domain/model"
)

// StrategyService は戦略管理のドメインサービス
type StrategyService interface {
	// CreateStrategy は新しい戦略を作成
	CreateStrategy(ctx context.Context, req *CreateStrategyRequest) (*model.Strategy, error)

	// GetStrategy は戦略を取得
	GetStrategy(ctx context.Context, id string) (*model.Strategy, error)

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

	// ValidateStrategy は戦略の設定を検証
	ValidateStrategy(ctx context.Context, strategy *model.Strategy) error

	// CheckRiskLimits は戦略のリスク制限をチェック
	CheckRiskLimits(ctx context.Context, strategy *model.Strategy) ([]string, error)

	// UpdateStatistics は戦略の統計を更新
	UpdateStatistics(ctx context.Context, strategyID string, pl float64, isWin bool) error
}

// CreateStrategyRequest は戦略作成リクエスト
type CreateStrategyRequest struct {
	Name        string               `json:"name"`
	Type        model.StrategyType   `json:"type"`
	Description string               `json:"description"`
	Config      model.StrategyConfig `json:"config"`
	RiskLimits  model.RiskLimits     `json:"risk_limits"`
	CreatedBy   string               `json:"created_by"`
}

// StrategyExecutor は戦略実行のインターフェース
type StrategyExecutor interface {
	// Execute は戦略を実行
	Execute(ctx context.Context, strategy *model.Strategy) error

	// CanExecute は戦略が実行可能かどうかを判定
	CanExecute(ctx context.Context, strategy *model.Strategy) (bool, error)

	// GetSupportedTypes はサポートされている戦略タイプを取得
	GetSupportedTypes() []model.StrategyType
}
