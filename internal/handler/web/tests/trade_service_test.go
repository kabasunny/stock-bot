package tests

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/gen/trade"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewTradeService はTradeServiceの作成をテストします
func TestNewTradeService(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	assert.NotNil(t, tradeService, "TradeService should not be nil")
}

// TestTradeService_GetSession はGetSession()の動作をテストします
func TestTradeService_GetSession(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedSession := &model.Session{
		SessionID: "test-session-id",
		UserID:    "test-user",
		LoginTime: time.Now(),
	}
	mockTradeService.On("GetSession").Return(expectedSession)

	ctx := context.Background()
	result, err := tradeService.GetSession(ctx)

	assert.NoError(t, err, "GetSession should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, "test-session-id", result.SessionID, "Should return correct session ID")
	assert.Equal(t, "test-user", result.UserID, "Should return correct user ID")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_GetSession_NoSession はセッションがない場合のGetSession()をテストします
func TestTradeService_GetSession_NoSession(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定（nilを返す）
	mockTradeService.On("GetSession").Return(nil)

	ctx := context.Background()
	result, err := tradeService.GetSession(ctx)

	assert.Error(t, err, "Should return error when no session")
	assert.Contains(t, err.Error(), "no active session", "Should indicate no active session")
	assert.Nil(t, result, "Result should be nil on error")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_GetPositions はGetPositions()の動作をテストします
func TestTradeService_GetPositions(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedPositions := []*model.Position{
		{
			Symbol:              "1301",
			PositionType:        model.PositionTypeLong,
			PositionAccountType: model.PositionAccountTypeCash,
			AveragePrice:        1500.0,
			Quantity:            100,
		},
	}
	mockTradeService.On("GetPositions", mock.Anything).Return(expectedPositions, nil)

	ctx := context.Background()
	result, err := tradeService.GetPositions(ctx)

	assert.NoError(t, err, "GetPositions should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Len(t, result.Positions, 1, "Should return 1 position")
	assert.Equal(t, "1301", result.Positions[0].Symbol, "Should return correct symbol")
	assert.Equal(t, "LONG", result.Positions[0].PositionType, "Should convert position type")
	assert.Equal(t, "CASH", result.Positions[0].PositionAccountType, "Should convert account type")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_GetOrders はGetOrders()の動作をテストします
func TestTradeService_GetOrders(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedOrders := []*model.Order{
		{
			OrderID:             "order-123",
			Symbol:              "1301",
			TradeType:           model.TradeTypeBuy,
			OrderType:           model.OrderTypeLimit,
			Quantity:            100,
			Price:               1500.0,
			OrderStatus:         model.OrderStatusNew,
			PositionAccountType: model.PositionAccountTypeCash,
		},
	}
	mockTradeService.On("GetOrders", mock.Anything).Return(expectedOrders, nil)

	ctx := context.Background()
	result, err := tradeService.GetOrders(ctx)

	assert.NoError(t, err, "GetOrders should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Len(t, result.Orders, 1, "Should return 1 order")
	assert.Equal(t, "order-123", result.Orders[0].OrderID, "Should return correct order ID")
	assert.Equal(t, "BUY", result.Orders[0].TradeType, "Should convert trade type")
	assert.Equal(t, "LIMIT", result.Orders[0].OrderType, "Should convert order type")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_GetBalance はGetBalance()の動作をテストします
func TestTradeService_GetBalance(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}
	mockTradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil)

	ctx := context.Background()
	result, err := tradeService.GetBalance(ctx)

	assert.NoError(t, err, "GetBalance should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1000000.0, result.Cash, "Should return correct cash")
	assert.Equal(t, 800000.0, result.BuyingPower, "Should return correct buying power")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_PlaceOrder はPlaceOrder()の動作をテストします
func TestTradeService_PlaceOrder(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedOrder := &model.Order{
		OrderID:             "new-order-123",
		Symbol:              "1301",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeLimit,
		Quantity:            100,
		Price:               1500.0,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: model.PositionAccountTypeCash,
	}
	mockTradeService.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*service.PlaceOrderRequest")).Return(expectedOrder, nil)

	ctx := context.Background()
	payload := &trade.PlaceOrderPayload{
		Symbol:              "1301",
		TradeType:           "BUY",
		OrderType:           "LIMIT",
		Quantity:            100,
		Price:               1500.0,
		TriggerPrice:        0.0,
		PositionAccountType: "CASH",
	}

	result, err := tradeService.PlaceOrder(ctx, payload)

	assert.NoError(t, err, "PlaceOrder should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, "new-order-123", result.OrderID, "Should return correct order ID")
	assert.Equal(t, "BUY", result.TradeType, "Should return correct trade type")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_CancelOrder はCancelOrder()の動作をテストします
func TestTradeService_CancelOrder(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	mockTradeService.On("CancelOrder", mock.Anything, "order-123").Return(nil)

	ctx := context.Background()
	payload := &trade.CancelOrderPayload{
		OrderID: "order-123",
	}

	err := tradeService.CancelOrder(ctx, payload)

	assert.NoError(t, err, "CancelOrder should not return error")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_GetPriceHistory はGetPriceHistory()の動作をテストします
func TestTradeService_GetPriceHistory(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	// モックの設定
	expectedHistory := []*service.HistoricalPrice{
		{
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Open:   1500.0,
			High:   1550.0,
			Low:    1480.0,
			Close:  1520.0,
			Volume: 1000000,
		},
	}
	mockTradeService.On("GetPriceHistory", mock.Anything, "1301", 30).Return(expectedHistory, nil)

	ctx := context.Background()
	payload := &trade.GetPriceHistoryPayload{
		Symbol: "1301",
		Days:   30,
	}

	result, err := tradeService.GetPriceHistory(ctx, payload)

	assert.NoError(t, err, "GetPriceHistory should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, "1301", result.Symbol, "Should return correct symbol")
	assert.Len(t, result.History, 1, "Should return 1 history item")
	assert.Equal(t, 1520.0, result.History[0].Close, "Should return correct close price")

	mockTradeService.AssertExpectations(t)
}

// TestTradeService_HealthCheck はHealthCheck()の動作をテストします
func TestTradeService_HealthCheck(t *testing.T) {
	mockTradeService := &MockTradeService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := web.NewTradeService(mockTradeService, logger, session)

	ctx := context.Background()
	result, err := tradeService.HealthCheck(ctx)

	// MockTradeServiceは*tradeservice.GoaTradeServiceではないため、フォールバックが実行される
	assert.NoError(t, err, "HealthCheck should not return error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, "unhealthy", result.Status, "Should return unhealthy status (fallback)")
	assert.Nil(t, result.SessionValid, "SessionValid should be nil (fallback)")

	// MockTradeServiceのHealthCheckは呼ばれない（型アサーション失敗のため）
}
