package tests

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/gen/trade"
	"stock-bot/internal/handler/web"
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
	shouldError     bool
	errorMessage    string
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
				OrderID:             "31000001",
				Symbol:              "6658",
				TradeType:           model.TradeTypeBuy,
				OrderType:           model.OrderTypeLimit,
				Quantity:            100,
				Price:               8000.0,
				OrderStatus:         model.OrderStatusNew,
				PositionAccountType: model.PositionAccountTypeCash,
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
		orderHistory: []*model.Order{
			{
				OrderID:             "31000002",
				Symbol:              "7203",
				TradeType:           model.TradeTypeBuy,
				OrderType:           model.OrderTypeLimit,
				Quantity:            100,
				Price:               2500.0,
				OrderStatus:         model.OrderStatusFilled,
				PositionAccountType: model.PositionAccountTypeCash,
				Executions: []model.Execution{
					{
						ExecutionID: "EX001",
						Quantity:    100,
						Price:       2500.0,
						ExecutedAt:  time.Date(2024, 12, 30, 9, 5, 0, 0, time.UTC),
					},
				},
			},
		},
	}
}

// SetError はモックでエラーを発生させる設定
func (m *MockTradeService) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

// GetSession は現在のAPIセッション情報を取得する
func (m *MockTradeService) GetSession() *model.Session {
	return m.session
}

// GetPositions は現在の保有ポジションを取得する
func (m *MockTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.positions, nil
}

// GetOrders は発注中の注文を取得する
func (m *MockTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.orders, nil
}

// GetBalance は口座残高を取得する
func (m *MockTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.balance, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (m *MockTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.priceHistory, nil
}

// PlaceOrder は注文を発行する
func (m *MockTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
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
	if m.shouldError {
		return fmt.Errorf("%s", m.errorMessage)
	}
	m.cancelledOrders = append(m.cancelledOrders, orderID)
	return nil
}

// CorrectOrder は注文を訂正する
func (m *MockTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	order := &model.Order{
		OrderID:             orderID,
		Symbol:              "7203",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeLimit,
		Quantity:            100,
		Price:               2500.0,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: model.PositionAccountTypeCash,
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
	if m.shouldError {
		return 0, fmt.Errorf("%s", m.errorMessage)
	}
	m.cancelAllCount = len(m.orders)
	return m.cancelAllCount, nil
}

// GetOrderHistory は注文履歴を取得する
func (m *MockTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.orderHistory, nil
}

// HealthCheck はサービスの健康状態をチェックする
func (m *MockTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return m.healthStatus, nil
}

// TestTradeServiceHandler_GetSession はGET /trade/sessionをテストします
func TestTradeServiceHandler_GetSession(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetSession(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-session-123", result.SessionID)
	assert.Equal(t, "test-user", result.UserID)
	assert.NotEmpty(t, result.LoginTime)
}

// TestTradeServiceHandler_GetSession_NoSession はセッションがない場合をテストします
func TestTradeServiceHandler_GetSession_NoSession(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.session = nil
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetSession(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no active session")
}

// TestTradeServiceHandler_GetPositions はGET /trade/positionsをテストします
func TestTradeServiceHandler_GetPositions(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetPositions(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Positions, 1)

	position := result.Positions[0]
	assert.Equal(t, "7203", position.Symbol)
	assert.Equal(t, "LONG", position.PositionType)
	assert.Equal(t, "CASH", position.PositionAccountType)
	assert.Equal(t, 2500.0, position.AveragePrice)
	assert.Equal(t, uint(100), position.Quantity)
}

// TestTradeServiceHandler_GetPositions_Error はエラー時の動作をテストします
func TestTradeServiceHandler_GetPositions_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "positions service error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetPositions(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get positions")
}

// TestTradeServiceHandler_GetOrders はGET /trade/ordersをテストします
func TestTradeServiceHandler_GetOrders(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetOrders(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Orders, 1)

	order := result.Orders[0]
	assert.Equal(t, "31000001", order.OrderID)
	assert.Equal(t, "6658", order.Symbol)
	assert.Equal(t, "BUY", order.TradeType)
	assert.Equal(t, "LIMIT", order.OrderType)
	assert.Equal(t, uint(100), order.Quantity)
	assert.Equal(t, 8000.0, order.Price)
	assert.Equal(t, "NEW", order.OrderStatus)
	assert.Equal(t, "CASH", *order.PositionAccountType)
}

// TestTradeServiceHandler_GetOrders_Error はエラー時の動作をテストします
func TestTradeServiceHandler_GetOrders_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "orders service error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetOrders(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get orders")
}

// TestTradeServiceHandler_GetBalance はGET /trade/balanceをテストします
func TestTradeServiceHandler_GetBalance(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetBalance(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1000000.0, result.Cash)
	assert.Equal(t, 800000.0, result.BuyingPower)
}

// TestTradeServiceHandler_GetBalance_Error はエラー時の動作をテストします
func TestTradeServiceHandler_GetBalance_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "balance service error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.GetBalance(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get balance")
}

// TestTradeServiceHandler_GetPriceHistory はGET /trade/price/{symbol}/historyをテストします
func TestTradeServiceHandler_GetPriceHistory(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.GetPriceHistoryPayload{
		Symbol: "7203",
		Days:   5,
	}

	result, err := handler.GetPriceHistory(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "7203", result.Symbol)
	require.Len(t, result.History, 1)

	history := result.History[0]
	assert.Equal(t, "2024-12-30T00:00:00Z", history.Date)
	assert.Equal(t, 2480.0, history.Open)
	assert.Equal(t, 2520.0, history.High)
	assert.Equal(t, 2470.0, history.Low)
	assert.Equal(t, 2500.0, history.Close)
	assert.Equal(t, uint64(1000000), history.Volume)
}

// TestTradeServiceHandler_GetPriceHistory_DefaultDays はデフォルト日数をテストします
func TestTradeServiceHandler_GetPriceHistory_DefaultDays(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.GetPriceHistoryPayload{
		Symbol: "7203",
		Days:   0, // デフォルト値をテスト
	}

	result, err := handler.GetPriceHistory(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "7203", result.Symbol)
}

// TestTradeServiceHandler_PlaceOrder はPOST /trade/ordersをテストします
func TestTradeServiceHandler_PlaceOrder(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.PlaceOrderPayload{
		Symbol:              "7203",
		TradeType:           "BUY",
		OrderType:           "LIMIT",
		Quantity:            100,
		Price:               2500.0,
		PositionAccountType: "CASH",
	}

	result, err := handler.PlaceOrder(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "31000999", result.OrderID)
	assert.Equal(t, "7203", result.Symbol)
	assert.Equal(t, "BUY", result.TradeType)
	assert.Equal(t, "LIMIT", result.OrderType)
	assert.Equal(t, uint(100), result.Quantity)
	assert.Equal(t, 2500.0, result.Price)
	assert.Equal(t, "NEW", result.OrderStatus)
	assert.Equal(t, "CASH", *result.PositionAccountType)

	// モックに注文が記録されていることを確認
	require.Len(t, mockService.placedOrders, 1)
}

// TestTradeServiceHandler_PlaceOrder_Error はエラー時の動作をテストします
func TestTradeServiceHandler_PlaceOrder_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "place order error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.PlaceOrderPayload{
		Symbol:    "7203",
		TradeType: "BUY",
		OrderType: "LIMIT",
		Quantity:  100,
		Price:     2500.0,
	}

	result, err := handler.PlaceOrder(context.Background(), payload)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to place order")
}

// TestTradeServiceHandler_CancelOrder はDELETE /trade/orders/{id}をテストします
func TestTradeServiceHandler_CancelOrder(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.CancelOrderPayload{
		OrderID: "31000001",
	}

	err := handler.CancelOrder(context.Background(), payload)

	require.NoError(t, err)

	// モックにキャンセルが記録されていることを確認
	require.Len(t, mockService.cancelledOrders, 1)
	assert.Equal(t, "31000001", mockService.cancelledOrders[0])
}

// TestTradeServiceHandler_CancelOrder_Error はエラー時の動作をテストします
func TestTradeServiceHandler_CancelOrder_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "cancel order error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.CancelOrderPayload{
		OrderID: "31000001",
	}

	err := handler.CancelOrder(context.Background(), payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cancel order error")
}

// TestTradeServiceHandler_ValidateSymbol はGET /trade/symbols/{symbol}/validateをテストします
func TestTradeServiceHandler_ValidateSymbol(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.ValidateSymbolPayload{
		Symbol: "7203",
	}

	result, err := handler.ValidateSymbol(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Equal(t, "7203", result.Symbol)
	assert.NotNil(t, result.Name)
	assert.Equal(t, "銘柄名（取得不可）", *result.Name)
}

// TestTradeServiceHandler_GetOrderHistory はGET /trade/orders/historyをテストします
func TestTradeServiceHandler_GetOrderHistory(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	status := "FILLED"
	symbol := "7203"
	payload := &trade.GetOrderHistoryPayload{
		Status: &status,
		Symbol: &symbol,
		Limit:  10,
	}

	result, err := handler.GetOrderHistory(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Orders, 1)

	order := result.Orders[0]
	assert.Equal(t, "31000002", order.OrderID)
	assert.Equal(t, "7203", order.Symbol)
	assert.Equal(t, "BUY", order.TradeType)
	assert.Equal(t, "LIMIT", order.OrderType)
	assert.Equal(t, uint(100), order.Quantity)
	assert.Equal(t, 2500.0, order.Price)
	assert.Equal(t, "FILLED", order.OrderStatus)
	assert.NotEmpty(t, order.CreatedAt)
	assert.NotNil(t, order.UpdatedAt)

	// 約定履歴の確認
	require.Len(t, order.Executions, 1)
	execution := order.Executions[0]
	assert.Equal(t, "EX001", execution.ExecutionID)
	assert.Equal(t, uint(100), execution.ExecutedQuantity)
	assert.Equal(t, 2500.0, execution.ExecutedPrice)
	assert.NotEmpty(t, execution.ExecutedAt)
}

// TestTradeServiceHandler_GetOrderHistory_DefaultLimit はデフォルトリミットをテストします
func TestTradeServiceHandler_GetOrderHistory_DefaultLimit(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	payload := &trade.GetOrderHistoryPayload{
		Limit: 0, // デフォルト値をテスト
	}

	result, err := handler.GetOrderHistory(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestTradeServiceHandler_CorrectOrder はPUT /trade/orders/{id}をテストします
func TestTradeServiceHandler_CorrectOrder(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	newPrice := 2600.0
	newQuantity := uint(200)
	payload := &trade.CorrectOrderPayload{
		OrderID:  "31000001",
		Price:    &newPrice,
		Quantity: &newQuantity,
	}

	result, err := handler.CorrectOrder(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "31000001", result.OrderID)
	assert.Equal(t, 2600.0, result.Price)
	assert.Equal(t, uint(200), result.Quantity)

	// モックに訂正が記録されていることを確認
	correctedOrder, exists := mockService.correctedOrders["31000001"]
	require.True(t, exists)
	assert.Equal(t, 2600.0, correctedOrder.Price)
	assert.Equal(t, 200, correctedOrder.Quantity)
}

// TestTradeServiceHandler_CorrectOrder_Error はエラー時の動作をテストします
func TestTradeServiceHandler_CorrectOrder_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "correct order error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	newPrice := 2600.0
	payload := &trade.CorrectOrderPayload{
		OrderID: "31000001",
		Price:   &newPrice,
	}

	result, err := handler.CorrectOrder(context.Background(), payload)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to correct order")
}

// TestTradeServiceHandler_CancelAllOrders はDELETE /trade/ordersをテストします
func TestTradeServiceHandler_CancelAllOrders(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.CancelAllOrders(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.CancelledCount)
	assert.Equal(t, 1, mockService.cancelAllCount)
}

// TestTradeServiceHandler_CancelAllOrders_Error はエラー時の動作をテストします
func TestTradeServiceHandler_CancelAllOrders_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "cancel all orders error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.CancelAllOrders(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to cancel all orders")
}

// TestTradeServiceHandler_HealthCheck はGET /trade/healthをテストします
func TestTradeServiceHandler_HealthCheck(t *testing.T) {
	mockService := NewMockTradeService()
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.HealthCheck(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "unhealthy", result.Status) // フォールバック実装のため
	assert.NotEmpty(t, result.Timestamp)
}

// TestTradeServiceHandler_HealthCheck_Error はエラー時の動作をテストします
func TestTradeServiceHandler_HealthCheck_Error(t *testing.T) {
	mockService := NewMockTradeService()
	mockService.SetError(true, "health check error")
	handler := web.NewTradeService(mockService, slog.Default(), nil)

	result, err := handler.HealthCheck(context.Background())

	require.NoError(t, err) // HealthCheckはエラーでもフォールバック値を返す
	require.NotNil(t, result)
	assert.Equal(t, "unhealthy", result.Status)
}
