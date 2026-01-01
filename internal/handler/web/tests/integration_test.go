package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	tradesvr "stock-bot/gen/http/trade/server"
	"stock-bot/gen/trade"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"
)

// IntegrationTestSuite はHTTP API統合テストのスイート
type IntegrationTestSuite struct {
	server       *httptest.Server
	client       *http.Client
	tradeService *MockTradeService
}

// SetupIntegrationTest は統合テスト用のHTTPサーバーをセットアップ
func SetupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	// モックサービスを作成
	mockTradeService := &MockTradeService{}

	// Webサービスを作成
	webTradeService := web.NewTradeService(mockTradeService, logger, session)

	// Goaエンドポイントを作成
	tradeEndpoints := trade.NewEndpoints(webTradeService)

	// HTTPマルチプレクサーを作成
	mux := goahttp.NewMuxer()

	// TradeServiceのHTTPハンドラーをマウント
	tradeServer := tradesvr.New(tradeEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	tradesvr.Mount(mux, tradeServer)

	// テストサーバーを作成
	server := httptest.NewServer(middleware.RequestID()(mux))

	return &IntegrationTestSuite{
		server:       server,
		client:       server.Client(),
		tradeService: mockTradeService,
	}
}

// TeardownIntegrationTest は統合テストのクリーンアップ
func (suite *IntegrationTestSuite) TeardownIntegrationTest() {
	suite.server.Close()
}

// TestHTTPAPI_GetSession はGET /trade/sessionの統合テスト
func TestHTTPAPI_GetSession(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// モックの設定
	expectedSession := &model.Session{
		SessionID: "test-session-123",
		UserID:    "test-user",
		LoginTime: time.Now(),
	}
	suite.tradeService.On("GetSession").Return(expectedSession)

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/session")
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "test-session-123", result["session_id"])
	assert.Equal(t, "test-user", result["user_id"])
	assert.NotEmpty(t, result["login_time"])

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_GetPositions はGET /trade/positionsの統合テスト
func TestHTTPAPI_GetPositions(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

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
	suite.tradeService.On("GetPositions", mock.Anything).Return(expectedPositions, nil)

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/positions")
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	positions := result["positions"].([]interface{})
	assert.Len(t, positions, 1)

	position := positions[0].(map[string]interface{})
	assert.Equal(t, "1301", position["symbol"])
	assert.Equal(t, "LONG", position["position_type"])
	assert.Equal(t, "CASH", position["position_account_type"])
	assert.Equal(t, 1500.0, position["average_price"])
	assert.Equal(t, float64(100), position["quantity"])

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_GetBalance はGET /trade/balanceの統合テスト
func TestHTTPAPI_GetBalance(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// モックの設定
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}
	suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil)

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 1000000.0, result["cash"])
	assert.Equal(t, 800000.0, result["buying_power"])

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_PlaceOrder はPOST /trade/ordersの統合テスト
func TestHTTPAPI_PlaceOrder(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// モックの設定
	expectedOrder := &model.Order{
		OrderID:             "order-123",
		Symbol:              "1301",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeLimit,
		Quantity:            100,
		Price:               1500.0,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: model.PositionAccountTypeCash,
	}
	suite.tradeService.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req *service.PlaceOrderRequest) bool {
		return req.Symbol == "1301" &&
			req.TradeType == model.TradeTypeBuy &&
			req.OrderType == model.OrderTypeLimit &&
			req.Quantity == 100 &&
			req.Price == 1500.0 &&
			req.PositionAccountType == model.PositionAccountTypeCash
	})).Return(expectedOrder, nil)

	// リクエストボディを作成
	orderRequest := map[string]interface{}{
		"symbol":                "1301",
		"trade_type":            "BUY",
		"order_type":            "LIMIT",
		"quantity":              100,
		"price":                 1500.0,
		"position_account_type": "CASH",
	}

	requestBody, err := json.Marshal(orderRequest)
	require.NoError(t, err)

	// HTTPリクエストを送信
	resp, err := suite.client.Post(
		suite.server.URL+"/trade/orders",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "order-123", result["order_id"])
	assert.Equal(t, "1301", result["symbol"])
	assert.Equal(t, "BUY", result["trade_type"])
	assert.Equal(t, "LIMIT", result["order_type"])
	assert.Equal(t, float64(100), result["quantity"])
	assert.Equal(t, 1500.0, result["price"])
	assert.Equal(t, "NEW", result["order_status"])

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_CancelOrder はDELETE /trade/orders/{order_id}の統合テスト
func TestHTTPAPI_CancelOrder(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// モックの設定
	suite.tradeService.On("CancelOrder", mock.Anything, "order-123").Return(nil)

	// HTTPリクエストを作成
	req, err := http.NewRequest("DELETE", suite.server.URL+"/trade/orders/order-123", nil)
	require.NoError(t, err)

	// HTTPリクエストを送信
	resp, err := suite.client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_GetPriceHistory はGET /trade/price-history/{symbol}の統合テスト
func TestHTTPAPI_GetPriceHistory(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

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
	suite.tradeService.On("GetPriceHistory", mock.Anything, "1301", 30).Return(expectedHistory, nil)

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/price-history/1301?days=30")
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "1301", result["symbol"])

	history := result["history"].([]interface{})
	assert.Len(t, history, 1)

	item := history[0].(map[string]interface{})
	assert.Equal(t, 1520.0, item["close"])
	assert.Equal(t, float64(1000000), item["volume"])

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_HealthCheck はGET /trade/healthの統合テスト
func TestHTTPAPI_HealthCheck(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディを解析
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// MockTradeServiceは型アサーションに失敗するため、フォールバック結果を期待
	assert.Equal(t, "unhealthy", result["status"])
	assert.NotEmpty(t, result["timestamp"])
}

// TestHTTPAPI_ErrorHandling はエラーハンドリングの統合テスト
func TestHTTPAPI_ErrorHandling(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// モックでエラーを返すように設定
	suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("service unavailable"))

	// HTTPリクエストを送信
	resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
	require.NoError(t, err)
	defer resp.Body.Close()

	// エラーレスポンスを検証
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// エラーレスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// エラーメッセージが含まれていることを確認
	assert.Contains(t, string(body), "service unavailable")

	suite.tradeService.AssertExpectations(t)
}

// TestHTTPAPI_InvalidJSON は不正なJSONリクエストの統合テスト
func TestHTTPAPI_InvalidJSON(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.TeardownIntegrationTest()

	// 不正なJSONを送信
	invalidJSON := `{"symbol": "1301", "trade_type": "BUY", "quantity": "invalid"}`

	resp, err := suite.client.Post(
		suite.server.URL+"/trade/orders",
		"application/json",
		bytes.NewBufferString(invalidJSON),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// バリデーションエラーを期待
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
