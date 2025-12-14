// internal/infrastructure/repository/master_repository_impl.go
package repository

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cockroachdb/errors"
)

type masterRepositoryImpl struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) repository.MasterRepository {
	return &masterRepositoryImpl{db: db}
}

func (r *masterRepositoryImpl) Save(ctx context.Context, entity interface{}) error {
	result := r.db.WithContext(ctx).Create(entity)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save master data")
	}
	return nil
}

func (r *masterRepositoryImpl) SaveAll(ctx context.Context, entities []interface{}) error {
	// バルクインサートを検討（GORMの機能を利用）
	// ただし、データベースの種類によってはバルクインサートがサポートされていない場合があるので注意
	result := r.db.WithContext(ctx).Create(&entities)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save all master data")
	}
	return nil
}

func (r *masterRepositoryImpl) FindByIssueCode(ctx context.Context, issueCode string, entityType string) (interface{}, error) {
	// entityTypeに基づいて、適切なモデルを選択して検索
	var entity interface{}
	switch entityType {
	case "StockMaster":
		//entity = &model.StockMaster{} // ポインタにする必要があるので修正
		var stockMaster model.StockMaster
		result := r.db.WithContext(ctx).Where("issue_code = ?", issueCode).First(&stockMaster)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil // NotFoundの場合はnilを返す
			}
			return nil, errors.Wrap(result.Error, "failed to find StockMaster by issue code")
		}
		entity = &stockMaster // アドレスを渡す
	case "TickRule":
		var tickRule model.TickRule
		result := r.db.WithContext(ctx).Where("issue_code = ?", issueCode).First(&tickRule)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil // NotFoundの場合はnilを返す
			}
			return nil, errors.Wrap(result.Error, "failed to find TickRule by issue code")
		}
		entity = &tickRule // アドレスを渡す
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
	return entity, nil
}

// dbStockMaster is a DTO for database operations, excluding relational fields.
type dbStockMaster struct {
	IssueCode               string
	IssueName               string
	IssueNameShort          string
	IssueNameKana           string
	IssueNameEnglish        string
	MarketCode              string
	IndustryCode            string
	IndustryName            string
	TradingUnit             int
	ListedSharesOutstanding int64
	UpperLimit              float64
	LowerLimit              float64
}

// TableName explicitly sets the table name for the DTO to match the domain model's table.
func (dbStockMaster) TableName() string {
	return "stock_masters"
}

func (r *masterRepositoryImpl) UpsertStockMasters(ctx context.Context, stocks []*model.StockMaster) error {
	if len(stocks) == 0 {
		return nil
	}

	// Convert domain models to DB DTOs, excluding relational fields.
	dbStocks := make([]dbStockMaster, 0, len(stocks))
	for _, s := range stocks {
		dbStocks = append(dbStocks, dbStockMaster{
			// Copy fields from MasterBase if necessary for updates, but GORM handles them.
			IssueCode:               s.IssueCode,
			IssueName:               s.IssueName,
			IssueNameShort:          s.IssueNameShort,
			IssueNameKana:           s.IssueNameKana,
			IssueNameEnglish:        s.IssueNameEnglish,
			MarketCode:              s.MarketCode,
			IndustryCode:            s.IndustryCode,
			IndustryName:            s.IndustryName,
			TradingUnit:             s.TradingUnit,
			ListedSharesOutstanding: s.ListedSharesOutstanding,
			UpperLimit:              s.UpperLimit,
			LowerLimit:              s.LowerLimit,
		})
	}

	// Define columns to be updated on conflict.
	updateColumns := []string{
		"issue_name",
		"issue_name_short",
		"issue_name_kana",
		"issue_name_english",
		"market_code",
		"industry_code",
		"industry_name",
		"trading_unit",
		"listed_shares_outstanding",
		"upper_limit",
		"lower_limit",
		"updated_at", // Explicitly update timestamp
	}

	// Use the DTO slice for the Create operation. GORM will not see the relational field.
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "issue_code"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&dbStocks)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to upsert stock masters")
	}

	return nil
}

func (r *masterRepositoryImpl) UpsertTickRules(ctx context.Context, tickRules []*model.TickRule) error {
	if len(tickRules) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, tickRule := range tickRules {
			// 親オブジェクトだけを first-class citizen として扱う
			ruleToUpsert := &model.TickRule{
				TickUnitNumber: tickRule.TickUnitNumber,
				ApplicableDate: tickRule.ApplicableDate,
			}

			// 1. 親である TickRule を Upsert
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "tick_unit_number"}},
				DoUpdates: clause.AssignmentColumns([]string{"applicable_date", "updated_at"}),
			}).Create(ruleToUpsert).Error; err != nil {
				return errors.Wrapf(err, "failed to upsert tick rule: %s", tickRule.TickUnitNumber)
			}

			// 2. 関連する古い TickLevel を削除
			if err := tx.Where("tick_rule_unit_number = ?", tickRule.TickUnitNumber).Delete(&model.TickLevel{}).Error; err != nil {
				return errors.Wrapf(err, "failed to delete old tick levels for rule: %s", tickRule.TickUnitNumber)
			}

			// 3. 新しい TickLevel を一括で挿入
			if len(tickRule.TickLevels) > 0 {
				// 各 TickLevel に親の UnitNumber を設定する
				for i := range tickRule.TickLevels {
					tickRule.TickLevels[i].TickRuleUnitNumber = tickRule.TickUnitNumber
				}
				if err := tx.Create(&tickRule.TickLevels).Error; err != nil {
					return errors.Wrapf(err, "failed to create new tick levels for rule: %s", tickRule.TickUnitNumber)
				}
			}
		}
		return nil
	})
}

// 他のメソッドも同様に実装
