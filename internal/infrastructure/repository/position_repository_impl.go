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
		return errors.New("no position found to update highest price for symbol: " + symbol)
	}
	return nil
}

func (r *positionRepositoryImpl) UpsertPositionByExecution(ctx context.Context, execution *model.Execution) error {
	var position model.Position
	result := r.db.WithContext(ctx).Where("symbol = ?", execution.Symbol).First(&position)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// ポジションが存在しない場合
		if execution.TradeType == model.TradeTypeBuy {
			// 買い約定の場合は新規ポジションを作成
			newPosition := model.Position{
				Symbol:       execution.Symbol,
				AveragePrice: execution.Price,
				Quantity:     execution.Quantity,
				HighestPrice: execution.Price,
			}
			result := r.db.WithContext(ctx).Create(&newPosition)
			if result.Error != nil {
				return errors.Wrap(result.Error, "failed to create new position on buy execution")
			}
			return nil
		} else {
			// 売り約定でポジションが存在しない場合はエラー（または無視）
			return errors.Errorf("sell execution for non-existent position with symbol: %s", execution.Symbol)
		}
	} else if result.Error != nil {
		return errors.Wrap(result.Error, "failed to find position for upsert by execution")
	}

	// ポジションが存在する場合
	if execution.TradeType == model.TradeTypeBuy {
		// 買い増しの場合、平均取得単価と数量を更新
		position.AveragePrice = (position.AveragePrice*float64(position.Quantity) + execution.Price*float64(execution.Quantity)) / float64(position.Quantity+execution.Quantity)
		position.Quantity += execution.Quantity
		if execution.Price > position.HighestPrice {
			position.HighestPrice = execution.Price
		}
	} else if execution.TradeType == model.TradeTypeSell {
		// 売り約定の場合、数量を減らす
		position.Quantity -= execution.Quantity
	}

	// 最高値の更新 (買い・売り両方で可能性あり)
	if execution.Price > position.HighestPrice {
		position.HighestPrice = execution.Price
	}

	if position.Quantity <= 0 {
		// 数量が0以下になったらポジションを削除
		return r.DeletePosition(ctx, position.Symbol)
	} else {
		// 数量が正の場合はポジションを更新
		result = r.db.WithContext(ctx).Save(&position)
		if result.Error != nil {
			return errors.Wrap(result.Error, "failed to update position on execution")
		}
	}
	return nil
}

func (r *positionRepositoryImpl) DeletePosition(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&model.Position{})
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete position")
	}
	if result.RowsAffected == 0 {
		return errors.New("no position found to delete for symbol: " + symbol)
	}
	return nil
}
