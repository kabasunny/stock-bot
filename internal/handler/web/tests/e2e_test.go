package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

// E2ETestSuite はE2Eテスト用のスイート
type E2ETestSuite struct {
	server       *httptest.Server
	client       *http.Client
	tradeService *MockTradeService
}

// SetupE2ETest はE2Eテスト用のセットアップ
func SetupE2ETest(t *testing.T) *E2ETestSuite {
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
	server := httptest.NewServer(mux)

	return &E2ETestSuite{
		server:       server,
		client:       server.Client(),
		tradeService: mockTradeService,
	}
}

// TeardownE2ETest はE2Eテストのクリーンアップ
func (suite *E2ETestSuite) TeardownE2ETest() {
	suite.server.Close()
}

// TestE2E_CompleteTradingFlow は完全な取引フローのE2Eテスト
func TestE2E_CompleteTradingFlow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.TeardownE2ETest()

	// Step 1: セッション確認
	t.Run("Step1_CheckSession", func(t *testing.T) {
		expectedSession := &model.Session{
			SessionID: "test-session-123",
			UserID:    "test-user",
			LoginTime: time.Now(),
		}
		suite.tradeService.On("GetSession").Return(expectedSession).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/session")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "test-session-123", result["session_id"])
	})

	// Step 2: 残高確認
	t.Run("Step2_CheckBalance", func(t *testing.T) {
		expectedBalance := &service.Balance{
			Cash:        1000000.0,
			BuyingPower: 800000.0,
		}
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, 1000000.0, result["cash"])
		assert.Equal(t, 800000.0, result["buying_power"])
	})

	// Step 3: 現在のポジション確認
	t.Run("Step3_CheckPositions", func(t *testing.T) {
		expectedPositions := []*model.Position{} // 初期状態では空
		suite.tradeService.On("GetPositions", mock.Anything).Return(expectedPositions, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/positions")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		positions := result["positions"].([]interface{})
		assert.Len(t, positions, 0, "初期状態ではポジションは空")
	})

	// Step 4: 買い注文発行
	var orderID string
	t.Run("Step4_PlaceBuyOrder", func(t *testing.T) {
		orderID = "buy-order-123"
		expectedOrder := &model.Order{
			OrderID:             orderID,
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
				req.Price == 1500.0
		})).Return(expectedOrder, nil).Once()

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

		resp, err := suite.client.Post(
			suite.server.URL+"/trade/orders",
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, orderID, result["order_id"])
		assert.Equal(t, "NEW", result["order_status"])
	})

	// Step 5: 注文一覧確認
	t.Run("Step5_CheckOrders", func(t *testing.T) {
		expectedOrders := []*model.Order{
			{
				OrderID:             orderID,
				Symbol:              "1301",
				TradeType:           model.TradeTypeBuy,
				OrderType:           model.OrderTypeLimit,
				Quantity:            100,
				Price:               1500.0,
				OrderStatus:         model.OrderStatusNew,
				PositionAccountType: model.PositionAccountTypeCash,
			},
		}
		suite.tradeService.On("GetOrders", mock.Anything).Return(expectedOrders, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/orders")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		orders := result["orders"].([]interface{})
		assert.Len(t, orders, 1)

		order := orders[0].(map[string]interface{})
		assert.Equal(t, orderID, order["order_id"])
		assert.Equal(t, "NEW", order["order_status"])
	})

	// Step 6: 注文キャンセル
	t.Run("Step6_CancelOrder", func(t *testing.T) {
		suite.tradeService.On("CancelOrder", mock.Anything, orderID).Return(nil).Once()

		req, err := http.NewRequest("DELETE", suite.server.URL+"/trade/orders/"+orderID, nil)
		require.NoError(t, err)

		resp, err := suite.client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	// Step 7: キャンセル後の注文一覧確認
	t.Run("Step7_CheckOrdersAfterCancel", func(t *testing.T) {
		expectedOrders := []*model.Order{} // キャンセル後は空
		suite.tradeService.On("GetOrders", mock.Anything).Return(expectedOrders, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/orders")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		orders := result["orders"].([]interface{})
		assert.Len(t, orders, 0, "キャンセル後は注文一覧が空")
	})

	// 全てのモック期待値が満たされたことを確認
	suite.tradeService.AssertExpectations(t)
}

// TestE2E_StopLimitOrderFlow は逆指値指値注文のE2Eテスト
func TestE2E_StopLimitOrderFlow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.TeardownE2ETest()

	// Step 1: 逆指値指値注文発行
	var orderID string
	t.Run("Step1_PlaceStopLimitOrder", func(t *testing.T) {
		orderID = "stop-limit-order-123"
		triggerPrice := 1450.0

		expectedOrder := &model.Order{
			OrderID:             orderID,
			Symbol:              "1301",
			TradeType:           model.TradeTypeSell,
			OrderType:           model.OrderTypeStopLimit,
			Quantity:            100,
			Price:               1400.0,
			TriggerPrice:        triggerPrice,
			OrderStatus:         model.OrderStatusNew,
			PositionAccountType: model.PositionAccountTypeCash,
		}

		suite.tradeService.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req *service.PlaceOrderRequest) bool {
			return req.Symbol == "1301" &&
				req.TradeType == model.TradeTypeSell &&
				req.OrderType == model.OrderTypeStopLimit &&
				req.Quantity == 100 &&
				req.Price == 1400.0 &&
				req.TriggerPrice != nil &&
				*req.TriggerPrice == 1450.0
		})).Return(expectedOrder, nil).Once()

		orderRequest := map[string]interface{}{
			"symbol":                "1301",
			"trade_type":            "SELL",
			"order_type":            "STOP_LIMIT",
			"quantity":              100,
			"price":                 1400.0,
			"trigger_price":         1450.0,
			"position_account_type": "CASH",
		}

		requestBody, err := json.Marshal(orderRequest)
		require.NoError(t, err)

		resp, err := suite.client.Post(
			suite.server.URL+"/trade/orders",
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, orderID, result["order_id"])
		assert.Equal(t, "SELL", result["trade_type"])
		assert.Equal(t, "STOP_LIMIT", result["order_type"])
		assert.Equal(t, 1400.0, result["price"])
		assert.Equal(t, "NEW", result["order_status"])
	})

	// Step 2: 逆指値注文の確認
	t.Run("Step2_CheckStopLimitOrder", func(t *testing.T) {
		expectedOrders := []*model.Order{
			{
				OrderID:             orderID,
				Symbol:              "1301",
				TradeType:           model.TradeTypeSell,
				OrderType:           model.OrderTypeStopLimit,
				Quantity:            100,
				Price:               1400.0,
				TriggerPrice:        1450.0,
				OrderStatus:         model.OrderStatusNew,
				PositionAccountType: model.PositionAccountTypeCash,
			},
		}
		suite.tradeService.On("GetOrders", mock.Anything).Return(expectedOrders, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/orders")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		orders := result["orders"].([]interface{})
		assert.Len(t, orders, 1)

		order := orders[0].(map[string]interface{})
		assert.Equal(t, orderID, order["order_id"])
		assert.Equal(t, "STOP_LIMIT", order["order_type"])
		assert.Equal(t, 1400.0, order["price"])
	})

	suite.tradeService.AssertExpectations(t)
}

// TestE2E_ErrorRecoveryFlow はエラー回復フローのE2Eテスト
func TestE2E_ErrorRecoveryFlow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.TeardownE2ETest()

	// Step 1: サービスエラーが発生
	t.Run("Step1_ServiceError", func(t *testing.T) {
		suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("service temporarily unavailable")).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	// Step 2: サービス回復後の正常動作
	t.Run("Step2_ServiceRecovery", func(t *testing.T) {
		expectedBalance := &service.Balance{
			Cash:        1000000.0,
			BuyingPower: 800000.0,
		}
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()

		resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, 1000000.0, result["cash"])
	})

	suite.tradeService.AssertExpectations(t)
}

// TestE2E_ConcurrentRequests は並行リクエストのE2Eテスト
func TestE2E_ConcurrentRequests(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.TeardownE2ETest()

	// 複数の並行リクエストを設定
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}

	// 5回の並行リクエストを期待
	for i := 0; i < 5; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()
	}

	// 並行リクエストを実行
	results := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}()
	}

	// 全ての結果を確認
	for i := 0; i < 5; i++ {
		err := <-results
		assert.NoError(t, err, "並行リクエスト %d でエラーが発生", i+1)
	}

	suite.tradeService.AssertExpectations(t)
}
