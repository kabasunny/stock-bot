package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTradeService はTradeServiceのモック実装
type MockTradeService struct {
	session         *model.Session
	positions       []*model.Position
	orders          []*model.Order
	balance         *service.Balance
	priceHistory    []*service.HistoricalPrice
	healthStatus    *service.HealthStatus
	placedOrders    []*model.Order
	cancelledOrders []string
	correctedOrders map[string]*model.Order
	cancelAllCount  int
	orderHistory    []*model.Order
}

// NewMockTradeService はMockTradeServiceの新しいインスタンスを作成
func NewMockTradeService() *MockTradeService {
	return &MockTradeService{
		session: &model.Session{
			SessionID:    "test-session-123",
			UserID:       "test-user",
			LoginTime:    time.Now(),
			ExpiresAt:    time.Now().Add(8 * time.Hour),
			IsActive:     true,
			ResultCode:   "0",
			ResultText:   "Success",
			LastActivity: time.Now(),
		},
		positions: []*model.Position{
			{
				Symbol:              "7203",
				PositionType:        model.PositionTypeLong,
				PositionAccountType: model.PositionAccountTypeCash,
				AveragePrice:        2500.0,
				Quantity:            100,
			},
		},
		orders: []*model.Order{
			{
				OrderID:     "31000001",
				Symbol:      "6658",
				TradeType:   model.TradeTypeBuy,
				OrderType:   model.OrderTypeLimit,
				Quantity:    100,
				Price:       8000.0,
				OrderStatus: model.OrderStatusNew,
			},
		},
		balance: &service.Balance{
			Cash:        1000000.0,
			BuyingPower: 800000.0,
		},
		priceHistory: []*service.HistoricalPrice{
			{
				Date:   time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
				Open:   2480.0,
				High:   2520.0,
				Low:    2470.0,
				Close:  2500.0,
				Volume: 1000000,
			},
		},
		healthStatus: &service.HealthStatus{
			Status:             "healthy",
			Timestamp:          time.Now(),
			SessionValid:       true,
			DatabaseConnected:  true,
			WebSocketConnected: true,
		},
		placedOrders:    []*model.Order{},
		cancelledOrders: []string{},
		correctedOrders: make(map[string]*model.Order),
		orderHistory:    []*model.Order{},
	}
}

// GetSession は現在のAPIセッション情報を取得する
func (m *MockTradeService) GetSession() *model.Session {
	return m.session
}

// GetPositions は現在の保有ポジションを取得する
func (m *MockTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	return m.positions, nil
}

// GetOrders は発注中の注文を取得する
func (m *MockTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return m.orders, nil
}

// GetBalance は口座残高を取得する
func (m *MockTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	return m.balance, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (m *MockTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	return m.priceHistory, nil
}

// PlaceOrder は注文を発行する
func (m *MockTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	order := &model.Order{
		OrderID:             "31000999",
		Symbol:              req.Symbol,
		TradeType:           req.TradeType,
		OrderType:           req.OrderType,
		Quantity:            req.Quantity,
		Price:               req.Price,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: req.PositionAccountType,
	}
	m.placedOrders = append(m.placedOrders, order)
	return order, nil
}

// CancelOrder は注文をキャンセルする
func (m *MockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	m.cancelledOrders = append(m.cancelledOrders, orderID)
	return nil
}

// CorrectOrder は注文を訂正する
func (m *MockTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	order := &model.Order{
		OrderID:     orderID,
		Symbol:      "7203",
		TradeType:   model.TradeTypeBuy,
		OrderType:   model.OrderTypeLimit,
		Quantity:    100,
		Price:       2500.0,
		OrderStatus: model.OrderStatusNew,
	}

	if newPrice != nil {
		order.Price = *newPrice
	}
	if newQuantity != nil {
		order.Quantity = *newQuantity
	}

	m.correctedOrders[orderID] = order
	return order, nil
}

// CancelAllOrders は全ての未約定注文をキャンセルする
func (m *MockTradeService) CancelAllOrders(ctx context.Context) (int, error) {
	m.cancelAllCount = len(m.orders)
	return m.cancelAllCount, nil
}

// GetOrderHistory は注文履歴を取得する
func (m *MockTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	return m.orderHistory, nil
}

// HealthCheck はサービスの健康状態をチェックする
func (m *MockTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	return m.healthStatus, nil
}

// TestTradeService_GetSession はGetSessionメソッドをテストします
func TestTradeService_GetSession(t *testing.T) {
	mockService := NewMockTradeService()

	session := mockService.GetSession()

	require.NotNil(t, session)
	assert.Equal(t, "test-session-123", session.SessionID)
	assert.Equal(t, "test-user", session.UserID)
	assert.Equal(t, "0", session.ResultCode)
	assert.Equal(t, "Success", session.ResultText)
	assert.True(t, session.IsActive)
}

// TestTradeService_GetPositions はGetPositionsメソッドをテストします
func TestTradeService_GetPositions(t *testing.T) {
	mockService := NewMockTradeService()

	positions, err := mockService.GetPositions(context.Background())

	require.NoError(t, err)
	require.Len(t, positions, 1)

	position := positions[0]
	assert.Equal(t, "7203", position.Symbol)
	assert.Equal(t, model.PositionTypeLong, position.PositionType)
	assert.Equal(t, model.PositionAccountTypeCash, position.PositionAccountType)
	assert.Equal(t, 2500.0, position.AveragePrice)
	assert.Equal(t, 100, position.Quantity)
}

// TestTradeService_GetOrders はGetOrdersメソッドをテストします
func TestTradeService_GetOrders(t *testing.T) {
	mockService := NewMockTradeService()

	orders, err := mockService.GetOrders(context.Background())

	require.NoError(t, err)
	require.Len(t, orders, 1)

	order := orders[0]
	assert.Equal(t, "31000001", order.OrderID)
	assert.Equal(t, "6658", order.Symbol)
	assert.Equal(t, model.TradeTypeBuy, order.TradeType)
	assert.Equal(t, model.OrderTypeLimit, order.OrderType)
	assert.Equal(t, 100, order.Quantity)
	assert.Equal(t, 8000.0, order.Price)
	assert.Equal(t, model.OrderStatusNew, order.OrderStatus)
}

// TestTradeService_GetBalance はGetBalanceメソッドをテストします
func TestTradeService_GetBalance(t *testing.T) {
	mockService := NewMockTradeService()

	balance, err := mockService.GetBalance(context.Background())

	require.NoError(t, err)
	require.NotNil(t, balance)
	assert.Equal(t, 1000000.0, balance.Cash)
	assert.Equal(t, 800000.0, balance.BuyingPower)
}

// TestTradeService_PlaceOrder はPlaceOrderメソッドをテストします
func TestTradeService_PlaceOrder(t *testing.T) {
	mockService := NewMockTradeService()

	req := &service.PlaceOrderRequest{
		Symbol:              "7203",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeLimit,
		Quantity:            100,
		Price:               2500.0,
		PositionAccountType: model.PositionAccountTypeCash,
	}

	order, err := mockService.PlaceOrder(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, "31000999", order.OrderID)
	assert.Equal(t, "7203", order.Symbol)
	assert.Equal(t, model.TradeTypeBuy, order.TradeType)
	assert.Equal(t, model.OrderTypeLimit, order.OrderType)
	assert.Equal(t, 100, order.Quantity)
	assert.Equal(t, 2500.0, order.Price)
	assert.Equal(t, model.OrderStatusNew, order.OrderStatus)

	// モックに注文が記録されていることを確認
	require.Len(t, mockService.placedOrders, 1)
	assert.Equal(t, order, mockService.placedOrders[0])
}

// TestTradeService_CancelOrder はCancelOrderメソッドをテストします
func TestTradeService_CancelOrder(t *testing.T) {
	mockService := NewMockTradeService()

	err := mockService.CancelOrder(context.Background(), "31000001")

	require.NoError(t, err)

	// モックにキャンセルが記録されていることを確認
	require.Len(t, mockService.cancelledOrders, 1)
	assert.Equal(t, "31000001", mockService.cancelledOrders[0])
}

// TestTradeService_CorrectOrder はCorrectOrderメソッドをテストします
func TestTradeService_CorrectOrder(t *testing.T) {
	mockService := NewMockTradeService()

	newPrice := 2600.0
	newQuantity := 200

	order, err := mockService.CorrectOrder(context.Background(), "31000001", &newPrice, &newQuantity)

	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, "31000001", order.OrderID)
	assert.Equal(t, 2600.0, order.Price)
	assert.Equal(t, 200, order.Quantity)

	// モックに訂正が記録されていることを確認
	correctedOrder, exists := mockService.correctedOrders["31000001"]
	require.True(t, exists)
	assert.Equal(t, order, correctedOrder)
}

// TestTradeService_CancelAllOrders はCancelAllOrdersメソッドをテストします
func TestTradeService_CancelAllOrders(t *testing.T) {
	mockService := NewMockTradeService()

	count, err := mockService.CancelAllOrders(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, count) // モックには1つの注文がある
	assert.Equal(t, 1, mockService.cancelAllCount)
}

// TestTradeService_GetPriceHistory はGetPriceHistoryメソッドをテストします
func TestTradeService_GetPriceHistory(t *testing.T) {
	mockService := NewMockTradeService()

	history, err := mockService.GetPriceHistory(context.Background(), "7203", 5)

	require.NoError(t, err)
	require.Len(t, history, 1)

	price := history[0]
	assert.Equal(t, time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), price.Date)
	assert.Equal(t, 2480.0, price.Open)
	assert.Equal(t, 2520.0, price.High)
	assert.Equal(t, 2470.0, price.Low)
	assert.Equal(t, 2500.0, price.Close)
	assert.Equal(t, int64(1000000), price.Volume)
}

// TestTradeService_GetOrderHistory はGetOrderHistoryメソッドをテストします
func TestTradeService_GetOrderHistory(t *testing.T) {
	mockService := NewMockTradeService()

	status := model.OrderStatusFilled
	symbol := "7203"

	history, err := mockService.GetOrderHistory(context.Background(), &status, &symbol, 10)

	require.NoError(t, err)
	require.NotNil(t, history)
	// モックでは空の履歴を返す
	assert.Len(t, history, 0)
}

// TestTradeService_HealthCheck はHealthCheckメソッドをテストします
func TestTradeService_HealthCheck(t *testing.T) {
	mockService := NewMockTradeService()

	health, err := mockService.HealthCheck(context.Background())

	require.NoError(t, err)
	require.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.True(t, health.SessionValid)
	assert.True(t, health.DatabaseConnected)
	assert.True(t, health.WebSocketConnected)
	assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second)
}
