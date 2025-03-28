// internal/infrastructure/repository/order_repository_impl.go

package repository

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
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
