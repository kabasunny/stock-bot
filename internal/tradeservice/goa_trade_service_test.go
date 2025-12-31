package tradeservice

import (
	"context"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoaTradeService_GetSession(t *testing.T) {
	// Setup
	session := &client.Session{}
	logger := slog.Default()

	service := NewGoaTradeService(
		nil, // balanceClient not needed for this test
		nil, // orderClient not needed for this test
		nil, // priceClient not needed for this test
		nil, // orderRepo not needed for this test
		nil, // masterRepo not needed for this test
		session,
		logger,
	)

	// Execute
	result := service.GetSession()

	// Assert
	assert.Equal(t, session, result)
}

func TestGoaTradeService_ImplementsTradeService(t *testing.T) {
	// This test ensures that GoaTradeService implements service.TradeService interface
	var _ service.TradeService = (*GoaTradeService)(nil)
}

func TestGoaTradeService_CancelOrder_OrderNotFound(t *testing.T) {
	// 簡単なテスト：注文が見つからない場合
	service := NewGoaTradeService(
		nil, nil, nil,
		&mockOrderRepo{findResult: nil, findError: nil}, // 注文が見つからない
		nil,
		&client.Session{},
		slog.Default(),
	)

	err := service.CancelOrder(context.Background(), "NONEXISTENT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
}

func TestGoaTradeService_CancelOrder_OrderAlreadyFilled(t *testing.T) {
	// 簡単なテスト：注文が既に約定済み
	filledOrder := &model.Order{
		OrderID:     "FILLED_ORDER",
		OrderStatus: model.OrderStatusFilled,
	}

	service := NewGoaTradeService(
		nil, nil, nil,
		&mockOrderRepo{findResult: filledOrder, findError: nil},
		nil,
		&client.Session{},
		slog.Default(),
	)

	err := service.CancelOrder(context.Background(), "FILLED_ORDER")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be cancelled")
}

// 簡単なモック実装
type mockOrderRepo struct {
	findResult *model.Order
	findError  error
	saveError  error
}

func (m *mockOrderRepo) Save(ctx context.Context, order *model.Order) error {
	return m.saveError
}

func (m *mockOrderRepo) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	return m.findResult, m.findError
}

func (m *mockOrderRepo) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepo) UpdateOrderStatusByExecution(ctx context.Context, execution *model.Execution) error {
	return nil
}
