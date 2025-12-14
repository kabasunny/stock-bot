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

func (r *masterRepositoryImpl) UpsertStockMasters(ctx context.Context, stocks []*model.StockMaster) error {
	if len(stocks) == 0 {
		return nil
	}

	// GORMのリレーション処理を避けるため、マップのスライスに変換して処理する
	var stockMaps []map[string]interface{}
	for _, s := range stocks {
		stockMaps = append(stockMaps, map[string]interface{}{
			"issue_code":   s.IssueCode,
			"issue_name":   s.IssueName,
			"trading_unit": s.TradingUnit,
			"market_code":  s.MarketCode,
			"upper_limit":  s.UpperLimit,
			"lower_limit":  s.LowerLimit,
		})
	}

	// .Model()でテーブル名を指定し、.Create()にはマップのスライスを渡す
	result := r.db.WithContext(ctx).Model(&model.StockMaster{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "issue_code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"issue_name",
			"trading_unit",
			"market_code",
			"upper_limit",
			"lower_limit",
			"updated_at",
		}),
	}).Create(&stockMaps)

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
