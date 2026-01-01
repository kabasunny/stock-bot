package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/service"

	"github.com/stretchr/testify/mock"
)

// MockTradeService はservice.TradeServiceのモック実装
type MockTradeService struct {
	mock.Mock
}

func (m *MockTradeService) GetSession() *model.Session {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*model.Session)
}

func (m *MockTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Balance), args.Error(1)
}

func (m *MockTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	args := m.Called(ctx, symbol, days)
	return args.Get(0).([]*service.HistoricalPrice), args.Error(1)
}

func (m *MockTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	args := m.Called(ctx, orderID, newPrice, newQuantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockTradeService) CancelAllOrders(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	args := m.Called(ctx, status, symbol, limit)
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *MockTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.HealthStatus), args.Error(1)
}
