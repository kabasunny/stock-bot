// internal/infrastructure/repository/order_repository_impl.go

package repository

import (
	"context"
	"log/slog" // 追加
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause" // 追加
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Save(ctx context.Context, order *model.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save order")
	}
	return nil
}

func (r *orderRepositoryImpl) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(result.Error, "failed to find order by id")
	}
	return &order, nil
}

func (r *orderRepositoryImpl) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.db.WithContext(ctx).Where("order_status = ?", status).Find(&orders)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to find orders by status")
	}
	return orders, nil
}

func (r *orderRepositoryImpl) UpdateOrderStatusByExecution(ctx context.Context, execution *model.Execution) error {
	var order model.Order
	result := r.db.WithContext(ctx).Where("order_id = ?", execution.OrderID).Preload("Executions").First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.ErrorContext(ctx, "order with ID not found during execution update", "order_id", execution.OrderID, "execution_id", execution.ExecutionID)
			return errors.Errorf("order with ID %s not found", execution.OrderID)
		}
		return errors.Wrap(result.Error, "failed to find order by ID for status update")
	}

	// 新しい約定情報をデータベースに保存 (ExecutionID で重複を排除)
	// ExecutionID は約定ごとにユニークなので ON CONFLICT DO NOTHING で重複挿入を避ける
	execResult := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(execution)
	if execResult.Error != nil {
		return errors.Wrap(execResult.Error, "failed to save execution")
	}

	// 新しい約定が実際に挿入された場合のみ、Executions スライスに新しく保存された約定を追加 (Preloadでロードされたものに加えて)
	// DoNothing が機能した場合 (RowsAffected == 0)、その約定は既にデータベースに存在するため、再度加算しない
	if execResult.RowsAffected > 0 {
		order.Executions = append(order.Executions, *execution)
	}

	// 約定数量の合計を計算
	totalExecutedQuantity := 0
	for _, exec := range order.Executions {
		totalExecutedQuantity += exec.Quantity
	}

	// 注文ステータスを更新
	if totalExecutedQuantity == order.Quantity {
		order.OrderStatus = model.OrderStatusFilled
	} else if totalExecutedQuantity > 0 && totalExecutedQuantity < order.Quantity {
		order.OrderStatus = model.OrderStatusPartiallyFilled
	} else if totalExecutedQuantity == 0 {
		order.OrderStatus = model.OrderStatusNew // 約定が全くない場合は新規注文状態
	} else if totalExecutedQuantity > order.Quantity {
		// これはあってはならないケースだが、念のためエラーとしておく
		return errors.Errorf("total executed quantity (%d) exceeds order quantity (%d) for order %s", totalExecutedQuantity, order.Quantity, order.OrderID)
	}

	result = r.db.WithContext(ctx).Save(&order)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update order status and add execution")
	}

	return nil
}

// FindOrderHistory は注文履歴を取得する（フィルタリング・制限付き）
func (r *orderRepositoryImpl) FindOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	query := r.db.WithContext(ctx).Preload("Executions").Order("created_at DESC")

	// ステータスでフィルタ
	if status != nil {
		query = query.Where("order_status = ?", *status)
	}

	// 銘柄でフィルタ
	if symbol != nil && *symbol != "" {
		query = query.Where("symbol = ?", *symbol)
	}

	// 制限
	if limit > 0 {
		query = query.Limit(limit)
	}

	var orders []*model.Order
	result := query.Find(&orders)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to find order history")
	}

	return orders, nil
}
