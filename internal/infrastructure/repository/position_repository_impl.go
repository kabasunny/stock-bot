// internal/infrastructure/repository/position_repository_impl.go

package repository

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type positionRepositoryImpl struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) repository.PositionRepository {
	return &positionRepositoryImpl{db: db}
}

func (r *positionRepositoryImpl) Save(ctx context.Context, position *model.Position) error {
	result := r.db.WithContext(ctx).Create(position)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save position")
	}
	return nil
}

func (r *positionRepositoryImpl) FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) {
	var position model.Position
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&position)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(result.Error, "failed to find position by symbol")
	}
	return &position, nil
}

func (r *positionRepositoryImpl) FindAll(ctx context.Context) ([]*model.Position, error) {
	var positions []*model.Position
	result := r.db.WithContext(ctx).Find(&positions)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to find all positions")
	}
	return positions, nil
}

func (r *positionRepositoryImpl) UpdateHighestPrice(ctx context.Context, symbol string, price float64) error {
	result := r.db.WithContext(ctx).Model(&model.Position{}).Where("symbol = ?", symbol).Update("highest_price", price)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update highest price")
	}
	if result.RowsAffected == 0 {
		// 見つからないことをエラーとするか、単に何もしないかは設計によるが、
		// 呼び出し元はポジションが存在することを期待しているはずなのでエラーとする
		return errors.New("no position found to update highest price for symbol: " + symbol)
	}
	return nil
}
