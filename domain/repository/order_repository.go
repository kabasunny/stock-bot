package repository

import (
	"context"
	"stock-bot/domain/model"
)

type OrderRepository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) // 例: 特定のステータスの注文を検索
	// UpdateOrderStatusByExecution は約定情報に基づいて注文の状態を更新します。
	UpdateOrderStatusByExecution(ctx context.Context, execution *model.Execution) error
	// FindOrderHistory は注文履歴を取得します（フィルタリング・制限付き）
	FindOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error)
	// 他の必要なメソッドを定義
}
