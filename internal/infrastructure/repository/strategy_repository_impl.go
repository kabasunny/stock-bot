package repository

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"time"

	"gorm.io/gorm"
)

// StrategyEntity は戦略のデータベースエンティティ
type StrategyEntity struct {
	ID          string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"`
	Status      string `gorm:"type:varchar(50);not null" json:"status"`
	Description string `gorm:"type:text" json:"description"`

	// JSON形式で保存
	ConfigJSON     string `gorm:"type:text;column:config" json:"config"`
	RiskLimitsJSON string `gorm:"type:text;column:risk_limits" json:"risk_limits"`
	StatisticsJSON string `gorm:"type:text;column:statistics" json:"statistics"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy string    `gorm:"type:varchar(255)" json:"created_by"`
}

// TableName はテーブル名を指定
func (StrategyEntity) TableName() string {
	return "strategies"
}

// StrategyRepositoryImpl はStrategyRepositoryの実装
type StrategyRepositoryImpl struct {
	db *gorm.DB
}

// NewStrategyRepository は新しいStrategyRepositoryを作成
func NewStrategyRepository(db *gorm.DB) repository.StrategyRepository {
	return &StrategyRepositoryImpl{
		db: db,
	}
}

// Save は戦略を保存
func (r *StrategyRepositoryImpl) Save(ctx context.Context, strategy *model.Strategy) error {
	entity, err := r.toEntity(strategy)
	if err != nil {
		return fmt.Errorf("failed to convert strategy to entity: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to save strategy: %w", err)
	}

	return nil
}

// FindByID はIDで戦略を検索
func (r *StrategyRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Strategy, error) {
	var entity StrategyEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find strategy by id: %w", err)
	}

	strategy, err := r.toModel(&entity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entity to strategy: %w", err)
	}

	return strategy, nil
}

// FindByType は種別で戦略を検索
func (r *StrategyRepositoryImpl) FindByType(ctx context.Context, strategyType model.StrategyType) ([]*model.Strategy, error) {
	var entities []StrategyEntity
	if err := r.db.WithContext(ctx).Where("type = ?", string(strategyType)).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to find strategies by type: %w", err)
	}

	strategies := make([]*model.Strategy, 0, len(entities))
	for _, entity := range entities {
		strategy, err := r.toModel(&entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entity to strategy: %w", err)
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// FindByStatus はステータスで戦略を検索
func (r *StrategyRepositoryImpl) FindByStatus(ctx context.Context, status model.StrategyStatus) ([]*model.Strategy, error) {
	var entities []StrategyEntity
	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to find strategies by status: %w", err)
	}

	strategies := make([]*model.Strategy, 0, len(entities))
	for _, entity := range entities {
		strategy, err := r.toModel(&entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entity to strategy: %w", err)
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// FindActive はアクティブな戦略を検索
func (r *StrategyRepositoryImpl) FindActive(ctx context.Context) ([]*model.Strategy, error) {
	return r.FindByStatus(ctx, model.StrategyStatusActive)
}

// FindAll は全ての戦略を検索
func (r *StrategyRepositoryImpl) FindAll(ctx context.Context) ([]*model.Strategy, error) {
	var entities []StrategyEntity
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to find all strategies: %w", err)
	}

	strategies := make([]*model.Strategy, 0, len(entities))
	for _, entity := range entities {
		strategy, err := r.toModel(&entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entity to strategy: %w", err)
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// Update は戦略を更新
func (r *StrategyRepositoryImpl) Update(ctx context.Context, strategy *model.Strategy) error {
	entity, err := r.toEntity(strategy)
	if err != nil {
		return fmt.Errorf("failed to convert strategy to entity: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update strategy: %w", err)
	}

	return nil
}

// Delete は戦略を削除
func (r *StrategyRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&StrategyEntity{}).Error; err != nil {
		return fmt.Errorf("failed to delete strategy: %w", err)
	}

	return nil
}

// UpdateStatistics は戦略の統計を更新
func (r *StrategyRepositoryImpl) UpdateStatistics(ctx context.Context, id string, stats *model.StrategyStats) error {
	statsJSON, err := marshalJSON(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal statistics: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&StrategyEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"statistics": statsJSON,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update strategy statistics: %w", err)
	}

	return nil
}

// toEntity はドメインモデルをエンティティに変換
func (r *StrategyRepositoryImpl) toEntity(strategy *model.Strategy) (*StrategyEntity, error) {
	configJSON, err := marshalJSON(strategy.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	riskLimitsJSON, err := marshalJSON(strategy.RiskLimits)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal risk limits: %w", err)
	}

	statisticsJSON, err := marshalJSON(strategy.Statistics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal statistics: %w", err)
	}

	return &StrategyEntity{
		ID:             strategy.ID,
		Name:           strategy.Name,
		Type:           string(strategy.Type),
		Status:         string(strategy.Status),
		Description:    strategy.Description,
		ConfigJSON:     configJSON,
		RiskLimitsJSON: riskLimitsJSON,
		StatisticsJSON: statisticsJSON,
		CreatedAt:      strategy.CreatedAt,
		UpdatedAt:      strategy.UpdatedAt,
		CreatedBy:      strategy.CreatedBy,
	}, nil
}

// toModel はエンティティをドメインモデルに変換
func (r *StrategyRepositoryImpl) toModel(entity *StrategyEntity) (*model.Strategy, error) {
	var config model.StrategyConfig
	if err := unmarshalJSON(entity.ConfigJSON, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var riskLimits model.RiskLimits
	if err := unmarshalJSON(entity.RiskLimitsJSON, &riskLimits); err != nil {
		return nil, fmt.Errorf("failed to unmarshal risk limits: %w", err)
	}

	var statistics model.StrategyStats
	if err := unmarshalJSON(entity.StatisticsJSON, &statistics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal statistics: %w", err)
	}

	return &model.Strategy{
		ID:          entity.ID,
		Name:        entity.Name,
		Type:        model.StrategyType(entity.Type),
		Status:      model.StrategyStatus(entity.Status),
		Description: entity.Description,
		Config:      config,
		RiskLimits:  riskLimits,
		Statistics:  statistics,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
	}, nil
}
