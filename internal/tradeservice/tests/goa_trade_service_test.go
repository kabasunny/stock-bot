package tests

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/tradeservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoaTradeService_NewGoaTradeService はGoaTradeServiceの作成をテストします
func TestGoaTradeService_NewGoaTradeService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, // balanceClient
		nil, // orderClient
		nil, // priceClient
		nil, // orderRepo
		nil, // masterRepo
		session,
		logger,
	)

	assert.NotNil(t, tradeService, "GoaTradeService should not be nil")
}

// TestGoaTradeService_GetSession はGetSession()の動作をテストします
func TestGoaTradeService_GetSession(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	// テスト用のセッション情報を設定
	session.ResultCode = "0"
	session.ResultText = "Success"

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	result := tradeService.GetSession()

	assert.NotNil(t, result, "Session should not be nil")
	// SessionAdapterの変換結果を確認
	// 実際の変換ロジックに依存するため、基本的な存在確認のみ
}

// TestGoaTradeService_GetPositions はGetPositions()の動作をテストします
func TestGoaTradeService_GetPositions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	positions, err := tradeService.GetPositions(ctx)

	assert.NoError(t, err, "GetPositions should not return error")
	assert.NotNil(t, positions, "Positions should not be nil")
	assert.Equal(t, 0, len(positions), "Should return empty positions (stub implementation)")
}

// TestGoaTradeService_GetOrders はGetOrders()の動作をテストします
func TestGoaTradeService_GetOrders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	orders, err := tradeService.GetOrders(ctx)

	assert.NoError(t, err, "GetOrders should not return error")
	assert.NotNil(t, orders, "Orders should not be nil")
	assert.Equal(t, 0, len(orders), "Should return empty orders (stub implementation)")
}

// TestGoaTradeService_GetBalance はGetBalance()の動作をテストします
func TestGoaTradeService_GetBalance(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	// モックのBalanceClientを作成（実際の実装では依存性注入）
	tradeService := tradeservice.NewGoaTradeService(
		nil, // balanceClient - nilでもスタブ実装で動作
		nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	balance, err := tradeService.GetBalance(ctx)

	// スタブ実装では固定値を返すため、エラーが発生する可能性がある
	// nilクライアントの場合はエラーになることを確認
	if err != nil {
		assert.Contains(t, err.Error(), "failed to get balance", "Should return balance error")
	} else {
		assert.NotNil(t, balance, "Balance should not be nil")
		assert.Equal(t, 1000000.0, balance.Cash, "Should return stub cash value")
		assert.Equal(t, 800000.0, balance.BuyingPower, "Should return stub buying power")
	}
}

// TestGoaTradeService_GetPriceHistory はGetPriceHistory()の動作をテストします
func TestGoaTradeService_GetPriceHistory(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	history, err := tradeService.GetPriceHistory(ctx, "1301", 30)

	assert.NoError(t, err, "GetPriceHistory should not return error")
	assert.NotNil(t, history, "History should not be nil")
	assert.Equal(t, 0, len(history), "Should return empty history (stub implementation)")
}

// TestGoaTradeService_PlaceOrder_NilClients はnilクライアントでのPlaceOrder()をテストします
func TestGoaTradeService_PlaceOrder_NilClients(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, // balanceClient
		nil, // orderClient - nil
		nil, // priceClient
		nil, // orderRepo
		nil, // masterRepo
		session,
		logger,
	)

	ctx := context.Background()
	req := &service.PlaceOrderRequest{
		Symbol:              "1301",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeMarket,
		Quantity:            100,
		Price:               0,
		PositionAccountType: model.PositionAccountTypeCash,
	}

	order, err := tradeService.PlaceOrder(ctx, req)

	// orderClientがnilの場合はエラーになることを確認
	assert.Error(t, err, "Should return error when orderClient is nil")
	assert.Contains(t, err.Error(), "orderClient is nil", "Should indicate nil orderClient")
	assert.Nil(t, order, "Order should be nil on error")
}

// TestGoaTradeService_PlaceOrder_NilSession はnilセッションでのPlaceOrder()をテストします
func TestGoaTradeService_PlaceOrder_NilSession(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		nil, // session - nil
		logger,
	)

	ctx := context.Background()
	req := &service.PlaceOrderRequest{
		Symbol:              "1301",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeMarket,
		Quantity:            100,
		Price:               0,
		PositionAccountType: model.PositionAccountTypeCash,
	}

	order, err := tradeService.PlaceOrder(ctx, req)

	// orderClientがnilの場合はエラーになることを確認
	assert.Error(t, err, "Should return error when orderClient is nil")
	assert.Contains(t, err.Error(), "orderClient is nil", "Should indicate nil orderClient")
	assert.Nil(t, order, "Order should be nil on error")
}

// TestGoaTradeService_CancelOrder はCancelOrder()の動作をテストします
func TestGoaTradeService_CancelOrder(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	err := tradeService.CancelOrder(ctx, "test-order-id")

	// スタブ実装では注文が見つからないエラーになる
	assert.Error(t, err, "Should return error for non-existent order")
	assert.Contains(t, err.Error(), "order not found", "Should indicate order not found")
}

// TestGoaTradeService_CorrectOrder はCorrectOrder()の動作をテストします
func TestGoaTradeService_CorrectOrder(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	newPrice := 1500.0
	newQuantity := 200

	order, err := tradeService.CorrectOrder(ctx, "test-order-id", &newPrice, &newQuantity)

	// スタブ実装では成功する
	assert.NoError(t, err, "CorrectOrder should not return error (stub)")
	assert.NotNil(t, order, "Order should not be nil")
	assert.Equal(t, "test-order-id", order.OrderID, "Should return correct order ID")
}

// TestGoaTradeService_CancelAllOrders はCancelAllOrders()の動作をテストします
func TestGoaTradeService_CancelAllOrders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	count, err := tradeService.CancelAllOrders(ctx)

	assert.NoError(t, err, "CancelAllOrders should not return error")
	assert.Equal(t, 0, count, "Should return 0 cancelled orders (stub)")
}

// TestGoaTradeService_GetOrderHistory はGetOrderHistory()の動作をテストします
func TestGoaTradeService_GetOrderHistory(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	status := model.OrderStatusNew
	symbol := "1301"

	orders, err := tradeService.GetOrderHistory(ctx, &status, &symbol, 100)

	assert.NoError(t, err, "GetOrderHistory should not return error")
	assert.NotNil(t, orders, "Orders should not be nil")
	assert.Equal(t, 0, len(orders), "Should return empty orders (stub)")
}

// TestGoaTradeService_HealthCheck はHealthCheck()の動作をテストします
func TestGoaTradeService_HealthCheck(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()
	session.ResultCode = "0" // 有効なセッション

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	health, err := tradeService.HealthCheck(ctx)

	assert.NoError(t, err, "HealthCheck should not return error")
	assert.NotNil(t, health, "Health should not be nil")
	assert.Equal(t, "healthy", health.Status, "Should return healthy status")
	assert.True(t, health.SessionValid, "Session should be valid")
	assert.True(t, health.DatabaseConnected, "Database should be connected (stub)")
	assert.True(t, health.WebSocketConnected, "WebSocket should be connected (stub)")
	assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second, "Timestamp should be recent")
}

// TestGoaTradeService_GetStockInfo はGetStockInfo()の動作をテストします
func TestGoaTradeService_GetStockInfo(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	ctx := context.Background()
	stockInfo, err := tradeService.GetStockInfo(ctx, "1301")

	assert.NoError(t, err, "GetStockInfo should not return error")
	assert.NotNil(t, stockInfo, "StockInfo should not be nil")
	assert.Equal(t, "1301", stockInfo.Symbol, "Should return correct symbol")
	assert.Equal(t, "テスト銘柄", stockInfo.Name, "Should return stub name")
	assert.Equal(t, 100, stockInfo.TradingUnit, "Should return stub trading unit")
	assert.Equal(t, "東証プライム", stockInfo.Market, "Should return stub market")
}

// TestGoaTradeService_SetUnifiedClient はSetUnifiedClient()の動作をテストします
func TestGoaTradeService_SetUnifiedClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	// UnifiedClientを設定
	unifiedClient := &client.TachibanaUnifiedClient{}

	// パニックしないことを確認
	require.NotPanics(t, func() {
		tradeService.SetUnifiedClient(unifiedClient)
	}, "SetUnifiedClient should not panic")
}
