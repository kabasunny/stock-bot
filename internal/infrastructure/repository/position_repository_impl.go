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
	db        *gorm.DB
	orderRepo repository.OrderRepository // Add order repository
}

func NewPositionRepository(db *gorm.DB, orderRepo repository.OrderRepository) repository.PositionRepository {
	return &positionRepositoryImpl{db: db, orderRepo: orderRepo}
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
	// 1. 関連するOrderを取得 (orderRepoは既に注入済み)
	order, err := r.orderRepo.FindByID(ctx, execution.OrderID)
	if err != nil {
		return errors.Wrapf(err, "failed to find order %s for execution %s", execution.OrderID, execution.ExecutionID)
	}
	if order == nil {
		return errors.Errorf("order with ID %s not found for execution %s", execution.OrderID, execution.ExecutionID)
	}

	var position model.Position
	// ポジションを検索する際に、SymbolだけでなくPositionAccountTypeも検索条件に含める
	result := r.db.WithContext(ctx).
		Where("symbol = ?", execution.Symbol).
		Where("position_account_type = ?", order.PositionAccountType). // Add PositionAccountType to search
		First(&position)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// ポジションが存在しない場合
		if execution.TradeType == model.TradeTypeBuy {
			// 買い約定の場合は新規ポジションを作成
			newPosition := model.Position{
				Symbol:              execution.Symbol,
				AveragePrice:        execution.Price,
				Quantity:            execution.Quantity,
				HighestPrice:        execution.Price,
				PositionAccountType: order.PositionAccountType, // OrderからAccountTypeを設定
				PositionType:        model.PositionTypeLong,    // 買いなのでLongポジション
			}
			result := r.db.WithContext(ctx).Create(&newPosition)
			if result.Error != nil {
				return errors.Wrap(result.Error, "failed to create new position on buy execution")
			}
			return nil
		} else {
			// 売り約定でポジションが存在しない場合はエラー（または無視）
			return errors.Errorf("sell execution for non-existent position with symbol: %s and account type: %s", execution.Symbol, order.PositionAccountType)
		}
	} else if result.Error != nil {
		return errors.Wrap(result.Error, "failed to find position for upsert by execution")
	}

	// ポジションが存在する場合 (この場合もAccountTypeをOrderから更新しておくが、検索条件にあるので通常は一致)
	// position.PositionAccountType = order.PositionAccountType // 検索条件にあるので不要だが、念のため

	// ポジションのPositionTypeも更新（売建からの買返済など、タイプが変わりうる場合を考慮）
	if execution.TradeType == model.TradeTypeBuy {
		position.PositionType = model.PositionTypeLong
	} else if execution.TradeType == model.TradeTypeSell {
		// 売りでQuantityが0になる場合は削除されるが、
		// 残る場合はPositionTypeShortに変わる可能性も考慮（例：空売り残高）
		// ただし、このAPIでは返済を別途考慮しているため、基本的にはLongの売却かShortの返済のはず。
		// ここでは、明示的に変更しない。必要であれば別のロジックで対応。
	}

	// ポジションが存在する場合の処理
	if execution.TradeType == model.TradeTypeBuy {
		// 買い約定の場合
		if order.PositionAccountType == model.PositionAccountTypeMarginRepay {
			// 信用返済（買返済）の場合：ショートポジションを減らす
			if position.PositionType != model.PositionTypeShort {
				return errors.Errorf("buy repayment for non-short position: %s", execution.Symbol)
			}
			position.Quantity -= execution.Quantity
		} else {
			// 通常の買い増し
			position.AveragePrice = (position.AveragePrice*float64(position.Quantity) + execution.Price*float64(execution.Quantity)) / float64(position.Quantity+execution.Quantity)
			position.Quantity += execution.Quantity
			position.PositionType = model.PositionTypeLong
		}

		if execution.Price > position.HighestPrice {
			position.HighestPrice = execution.Price
		}
	} else if execution.TradeType == model.TradeTypeSell {
		// 売り約定の場合
		if order.PositionAccountType == model.PositionAccountTypeMarginRepay {
			// 信用返済（売返済）の場合：ロングポジションを減らす
			if position.PositionType != model.PositionTypeLong {
				return errors.Errorf("sell repayment for non-long position: %s", execution.Symbol)
			}
			position.Quantity -= execution.Quantity
		} else if order.PositionAccountType == model.PositionAccountTypeMarginNew {
			// 信用新規売り（空売り）の場合
			if position.Quantity == 0 {
				// 新規ショートポジション
				position.PositionType = model.PositionTypeShort
				position.AveragePrice = execution.Price
				position.Quantity = execution.Quantity
			} else {
				// 既存ショートポジションに追加
				position.AveragePrice = (position.AveragePrice*float64(position.Quantity) + execution.Price*float64(execution.Quantity)) / float64(position.Quantity+execution.Quantity)
				position.Quantity += execution.Quantity
			}
		} else {
			// 現物売り：ロングポジションを減らす
			position.Quantity -= execution.Quantity
		}
	}

	// 最高値の更新 (買い・売り両方で可能性あり)
	if execution.Price > position.HighestPrice {
		position.HighestPrice = execution.Price
	}

	if position.Quantity <= 0 {
		// 数量が0以下になったらポジションを削除
		return r.DeletePosition(ctx, position.Symbol, position.PositionAccountType) // DeletePositionはSymbolとAccountTypeで検索するように修正
	} else {
		// 数量が正の場合はポジションを更新
		result = r.db.WithContext(ctx).Save(&position)
		if result.Error != nil {
			return errors.Wrap(result.Error, "failed to update position on execution")
		}
	}
	return nil
}

func (r *positionRepositoryImpl) DeletePosition(ctx context.Context, symbol string, accountType model.PositionAccountType) error {
	result := r.db.WithContext(ctx).Where("symbol = ? AND position_account_type = ?", symbol, accountType).Delete(&model.Position{})
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete position")
	}
	if result.RowsAffected == 0 {
		return errors.Errorf("no position found to delete for symbol: %s with account type: %s", symbol, accountType)
	}
	return nil
}
