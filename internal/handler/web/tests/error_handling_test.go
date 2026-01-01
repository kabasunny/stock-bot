package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

// ErrorHandlingTestSuite はエラーハンドリングテスト用のスイート
type ErrorHandlingTestSuite struct {
	server       *httptest.Server
	client       *http.Client
	tradeService *MockTradeService
}

// SetupErrorHandlingTest はエラーハンドリングテスト用のセットアップ
func SetupErrorHandlingTest(t *testing.T) *ErrorHandlingTestSuite {
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

	return &ErrorHandlingTestSuite{
		server:       server,
		client:       server.Client(),
		tradeService: mockTradeService,
	}
}

// TeardownErrorHandlingTest はエラーハンドリングテストのクリーンアップ
func (suite *ErrorHandlingTestSuite) TeardownErrorHandlingTest() {
	suite.server.Close()
}

// TestErrorHandling_NetworkTimeout はネットワークタイムアウトエラーのテスト
func TestErrorHandling_NetworkTimeout(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// タイムアウトを短く設定したクライアント
	timeoutClient := &http.Client{
		Timeout: 1 * time.Millisecond, // 非常に短いタイムアウト
	}

	// サービスエラーを設定（遅延をシミュレート）
	suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("timeout")).Maybe()

	// タイムアウトが発生することを期待
	resp, err := timeoutClient.Get(suite.server.URL + "/trade/balance")

	// タイムアウトエラーまたは接続エラーが発生することを確認
	assert.Error(t, err, "Should return timeout error")
	if resp != nil {
		resp.Body.Close()
	}
}

// TestErrorHandling_ServiceUnavailable はサービス利用不可エラーのテスト
func TestErrorHandling_ServiceUnavailable(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// サービス利用不可エラーを設定
	suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("service unavailable")).Once()

	resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Goaフレームワークのエラー形式を確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "fault", result["name"])

	suite.tradeService.AssertExpectations(t)
}

// TestErrorHandling_AuthenticationError は認証エラーのテスト
func TestErrorHandling_AuthenticationError(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// 認証エラーを設定（nilセッション）
	suite.tradeService.On("GetSession").Return(nil).Once()

	resp, err := suite.client.Get(suite.server.URL + "/trade/session")
	require.NoError(t, err)
	defer resp.Body.Close()

	// nilセッションの場合は500エラーが返される
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// エラー情報が含まれていることを確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "fault", result["name"])

	suite.tradeService.AssertExpectations(t)
}

// TestErrorHandling_InvalidOrderData は不正な注文データのテスト
func TestErrorHandling_InvalidOrderData(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// 不正な注文データ
	invalidOrderRequest := map[string]interface{}{
		"symbol":     "", // 空のシンボル
		"trade_type": "INVALID_TYPE",
		"quantity":   -100,            // 負の数量
		"price":      "invalid_price", // 不正な価格
	}

	requestBody, err := json.Marshal(invalidOrderRequest)
	require.NoError(t, err)

	resp, err := suite.client.Post(
		suite.server.URL+"/trade/orders",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// バリデーションエラーが返されることを確認
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Goaフレームワークのデコードエラー形式を確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "decode_payload", result["name"])
}

// TestErrorHandling_MalformedJSON は不正なJSONのテスト
func TestErrorHandling_MalformedJSON(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// 不正なJSON
	malformedJSON := `{"symbol": "1301", "trade_type": "BUY", "quantity": 100, "price": 1500,}` // 末尾のカンマが不正

	resp, err := suite.client.Post(
		suite.server.URL+"/trade/orders",
		"application/json",
		bytes.NewBufferString(malformedJSON),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// JSONパースエラーが返されることを確認
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Goaフレームワークのデコードエラー形式を確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "decode_payload", result["name"])
}

// TestErrorHandling_OrderNotFound は存在しない注文のテスト
func TestErrorHandling_OrderNotFound(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	nonExistentOrderID := "non-existent-order-123"

	// 注文が見つからないエラーを設定
	suite.tradeService.On("CancelOrder", mock.Anything, nonExistentOrderID).Return(fmt.Errorf("order not found")).Once()

	req, err := http.NewRequest("DELETE", suite.server.URL+"/trade/orders/"+nonExistentOrderID, nil)
	require.NoError(t, err)

	resp, err := suite.client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// エラー情報が含まれていることを確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "fault", result["name"])

	suite.tradeService.AssertExpectations(t)
}

// TestErrorHandling_RateLimitError はレート制限エラーのテスト
func TestErrorHandling_RateLimitError(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// レート制限エラーを設定
	suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("rate limit exceeded")).Once()

	resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// エラーメッセージが含まれていることを確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "fault", result["name"])

	suite.tradeService.AssertExpectations(t)
}

// TestErrorHandling_ConcurrentErrors は並行処理でのエラーハンドリングのテスト
func TestErrorHandling_ConcurrentErrors(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// 複数の並行エラーを設定
	for i := 0; i < 5; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(nil, fmt.Errorf("concurrent error %d", i)).Once()
	}

	// 並行リクエストを実行
	results := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func(index int) {
			resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusInternalServerError {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// 全ての結果を確認
	for i := 0; i < 5; i++ {
		err := <-results
		assert.NoError(t, err, "並行エラーハンドリング %d でエラーが発生", i+1)
	}

	suite.tradeService.AssertExpectations(t)
}

// TestErrorHandling_ContextCancellation はコンテキストキャンセレーションのテスト
func TestErrorHandling_ContextCancellation(t *testing.T) {
	suite := SetupErrorHandlingTest(t)
	defer suite.TeardownErrorHandlingTest()

	// コンテキストキャンセレーションエラーを設定
	suite.tradeService.On("GetBalance", mock.Anything).Return(nil, context.Canceled).Once()

	resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// エラー情報が含まれていることを確認
	assert.Contains(t, result, "name")
	assert.Equal(t, "fault", result["name"])

	suite.tradeService.AssertExpectations(t)
}
