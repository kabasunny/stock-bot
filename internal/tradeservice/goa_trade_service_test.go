package tradeservice

import (
	"context"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoaTradeService_GetSession(t *testing.T) {
	// Setup
	clientSession := &client.Session{
		ResultCode: "0",
		ResultText: "Success",
	}
	logger := slog.Default()

	service := NewGoaTradeService(
		nil, // balanceClient not needed for this test
		nil, // orderClient not needed for this test
		nil, // priceClient not needed for this test
		nil, // orderRepo not needed for this test
		nil, // masterRepo not needed for this test
		clientSession,
		logger,
	)

	// Execute
	result := service.GetSession()

	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "0", result.ResultCode)
	assert.Equal(t, "Success", result.ResultText)
	assert.True(t, result.IsActive)
	assert.NotEmpty(t, result.SessionID)
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

// TestGoaTradeService_HealthCheck はHealthCheckメソッドをテストします
func TestGoaTradeService_HealthCheck(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{}
	session := &client.Session{ResultCode: "0"}
	service := NewGoaTradeService(nil, nil, nil, mockOrderRepo, nil, session, slog.Default())

	health, err := service.HealthCheck(context.Background())

	require.NoError(t, err)
	require.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.True(t, health.SessionValid)
	assert.True(t, health.DatabaseConnected)
	assert.True(t, health.WebSocketConnected)
	assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second)
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

func (m *mockOrderRepo) FindOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepo) UpdateOrderStatusByExecution(ctx context.Context, execution *model.Execution) error {
	return nil
}
