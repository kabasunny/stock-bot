// internal/infrastructure/repository/master_repository_impl.go
package repository

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"gorm.io/gorm"

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

// 他のメソッドも同様に実装
