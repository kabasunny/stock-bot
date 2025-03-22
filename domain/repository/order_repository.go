package repository

import (
	"context"
	"myapp/domain/model"
)

type OrderRepository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) // 例: 特定のステータスの注文を検索
	// 他の必要なメソッドを定義
}
