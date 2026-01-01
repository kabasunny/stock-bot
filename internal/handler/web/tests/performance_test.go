package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	tradesvr "stock-bot/gen/http/trade/server"
	"stock-bot/gen/trade"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	goahttp "goa.design/goa/v3/http"
)

// PerformanceTestSuite はパフォーマンステスト用のスイート
type PerformanceTestSuite struct {
	server       *httptest.Server
	client       *http.Client
	tradeService *MockTradeService
}

// SetupPerformanceTest はパフォーマンステスト用のセットアップ
func SetupPerformanceTest(t *testing.T) *PerformanceTestSuite {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})) // ログレベルを下げてパフォーマンス向上
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

	return &PerformanceTestSuite{
		server:       server,
		client:       server.Client(),
		tradeService: mockTradeService,
	}
}

// TeardownPerformanceTest はパフォーマンステストのクリーンアップ
func (suite *PerformanceTestSuite) TeardownPerformanceTest() {
	suite.server.Close()
}

// TestPerformance_ConcurrentConnections は同時接続数テスト
func TestPerformance_ConcurrentConnections(t *testing.T) {
	suite := SetupPerformanceTest(t)
	defer suite.TeardownPerformanceTest()

	concurrentUsers := 50
	requestsPerUser := 10

	// 期待される残高レスポンス
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}

	// 大量のモック期待値を設定
	totalRequests := concurrentUsers * requestsPerUser
	for i := 0; i < totalRequests; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()
	}

	var wg sync.WaitGroup
	results := make(chan error, totalRequests)
	startTime := time.Now()

	// 並行ユーザーを起動
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			for j := 0; j < requestsPerUser; j++ {
				resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
				if err != nil {
					results <- fmt.Errorf("user %d request %d: %w", userID, j, err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					resp.Body.Close()
					results <- fmt.Errorf("user %d request %d: unexpected status %d", userID, j, resp.StatusCode)
					continue
				}

				resp.Body.Close()
				results <- nil
			}
		}(i)
	}

	wg.Wait()
	close(results)

	duration := time.Since(startTime)

	// 結果を集計
	successCount := 0
	errorCount := 0
	for result := range results {
		if result == nil {
			successCount++
		} else {
			errorCount++
			t.Logf("Error: %v", result)
		}
	}

	// パフォーマンス指標を計算
	requestsPerSecond := float64(successCount) / duration.Seconds()

	t.Logf("同時接続テスト結果:")
	t.Logf("  並行ユーザー数: %d", concurrentUsers)
	t.Logf("  ユーザーあたりリクエスト数: %d", requestsPerUser)
	t.Logf("  総リクエスト数: %d", totalRequests)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  エラー数: %d", errorCount)
	t.Logf("  実行時間: %v", duration)
	t.Logf("  スループット: %.2f req/sec", requestsPerSecond)

	// 成功率が95%以上であることを確認
	successRate := float64(successCount) / float64(totalRequests) * 100
	assert.GreaterOrEqual(t, successRate, 95.0, "成功率が95%以上である必要があります")

	// スループットが最低限の値以上であることを確認
	assert.GreaterOrEqual(t, requestsPerSecond, 100.0, "スループットが100 req/sec以上である必要があります")

	suite.tradeService.AssertExpectations(t)
}

// TestPerformance_BulkOrderProcessing は大量注文処理テスト
func TestPerformance_BulkOrderProcessing(t *testing.T) {
	suite := SetupPerformanceTest(t)
	defer suite.TeardownPerformanceTest()

	orderCount := 100

	// 期待される注文レスポンス
	expectedOrder := &model.Order{
		OrderID:             "bulk-order",
		Symbol:              "1301",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeLimit,
		Quantity:            100,
		Price:               1500.0,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: model.PositionAccountTypeCash,
	}

	// 大量のモック期待値を設定
	for i := 0; i < orderCount; i++ {
		suite.tradeService.On("PlaceOrder", mock.Anything, mock.AnythingOfType("*service.PlaceOrderRequest")).Return(expectedOrder, nil).Once()
	}

	var wg sync.WaitGroup
	results := make(chan error, orderCount)
	startTime := time.Now()

	// 注文リクエストのテンプレート
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

	// 並行で注文を発行
	for i := 0; i < orderCount; i++ {
		wg.Add(1)
		go func(orderID int) {
			defer wg.Done()

			resp, err := suite.client.Post(
				suite.server.URL+"/trade/orders",
				"application/json",
				bytes.NewBuffer(requestBody),
			)
			if err != nil {
				results <- fmt.Errorf("order %d: %w", orderID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				results <- fmt.Errorf("order %d: unexpected status %d", orderID, resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	wg.Wait()
	close(results)

	duration := time.Since(startTime)

	// 結果を集計
	successCount := 0
	errorCount := 0
	for result := range results {
		if result == nil {
			successCount++
		} else {
			errorCount++
			t.Logf("Error: %v", result)
		}
	}

	// パフォーマンス指標を計算
	ordersPerSecond := float64(successCount) / duration.Seconds()

	t.Logf("大量注文処理テスト結果:")
	t.Logf("  総注文数: %d", orderCount)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  エラー数: %d", errorCount)
	t.Logf("  実行時間: %v", duration)
	t.Logf("  注文処理速度: %.2f orders/sec", ordersPerSecond)

	// 成功率が95%以上であることを確認
	successRate := float64(successCount) / float64(orderCount) * 100
	assert.GreaterOrEqual(t, successRate, 95.0, "成功率が95%以上である必要があります")

	// 注文処理速度が最低限の値以上であることを確認
	assert.GreaterOrEqual(t, ordersPerSecond, 50.0, "注文処理速度が50 orders/sec以上である必要があります")

	suite.tradeService.AssertExpectations(t)
}

// TestPerformance_MemoryUsage はメモリ使用量テスト
func TestPerformance_MemoryUsage(t *testing.T) {
	suite := SetupPerformanceTest(t)
	defer suite.TeardownPerformanceTest()

	// 初期メモリ使用量を測定
	runtime.GC()
	var initialMemStats runtime.MemStats
	runtime.ReadMemStats(&initialMemStats)

	requestCount := 1000

	// 期待される残高レスポンス
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}

	// 大量のモック期待値を設定
	for i := 0; i < requestCount; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()
	}

	startTime := time.Now()

	// 大量のリクエストを順次実行
	for i := 0; i < requestCount; i++ {
		resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// 定期的にGCを実行してメモリリークを検出
		if i%100 == 0 {
			runtime.GC()
		}
	}

	duration := time.Since(startTime)

	// 最終メモリ使用量を測定
	runtime.GC()
	var finalMemStats runtime.MemStats
	runtime.ReadMemStats(&finalMemStats)

	// メモリ使用量の変化を計算
	memoryIncrease := finalMemStats.Alloc - initialMemStats.Alloc
	memoryIncreaseKB := float64(memoryIncrease) / 1024

	t.Logf("メモリ使用量テスト結果:")
	t.Logf("  リクエスト数: %d", requestCount)
	t.Logf("  実行時間: %v", duration)
	t.Logf("  初期メモリ使用量: %d KB", initialMemStats.Alloc/1024)
	t.Logf("  最終メモリ使用量: %d KB", finalMemStats.Alloc/1024)
	t.Logf("  メモリ増加量: %.2f KB", memoryIncreaseKB)
	t.Logf("  リクエストあたりメモリ増加: %.2f bytes", float64(memoryIncrease)/float64(requestCount))

	// メモリ増加量が合理的な範囲内であることを確認（2MB以下）
	assert.LessOrEqual(t, memoryIncreaseKB, 2048.0, "メモリ増加量が2MB以下である必要があります")

	suite.tradeService.AssertExpectations(t)
}

// TestPerformance_ResponseTime はレスポンス時間テスト
func TestPerformance_ResponseTime(t *testing.T) {
	suite := SetupPerformanceTest(t)
	defer suite.TeardownPerformanceTest()

	requestCount := 100

	// 期待される残高レスポンス
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}

	// モック期待値を設定
	for i := 0; i < requestCount; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Once()
	}

	responseTimes := make([]time.Duration, requestCount)

	// レスポンス時間を測定
	for i := 0; i < requestCount; i++ {
		startTime := time.Now()

		resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
		require.NoError(t, err)

		responseTimes[i] = time.Since(startTime)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// 統計を計算
	var totalTime time.Duration
	minTime := responseTimes[0]
	maxTime := responseTimes[0]

	for _, responseTime := range responseTimes {
		totalTime += responseTime
		if responseTime < minTime {
			minTime = responseTime
		}
		if responseTime > maxTime {
			maxTime = responseTime
		}
	}

	avgTime := totalTime / time.Duration(requestCount)

	// 95パーセンタイルを計算（簡易版）
	sortedTimes := make([]time.Duration, len(responseTimes))
	copy(sortedTimes, responseTimes)

	// 簡易ソート
	for i := 0; i < len(sortedTimes)-1; i++ {
		for j := i + 1; j < len(sortedTimes); j++ {
			if sortedTimes[i] > sortedTimes[j] {
				sortedTimes[i], sortedTimes[j] = sortedTimes[j], sortedTimes[i]
			}
		}
	}

	p95Index := int(float64(len(sortedTimes)) * 0.95)
	if p95Index >= len(sortedTimes) {
		p95Index = len(sortedTimes) - 1
	}
	p95Time := sortedTimes[p95Index]

	t.Logf("レスポンス時間テスト結果:")
	t.Logf("  リクエスト数: %d", requestCount)
	t.Logf("  平均レスポンス時間: %v", avgTime)
	t.Logf("  最小レスポンス時間: %v", minTime)
	t.Logf("  最大レスポンス時間: %v", maxTime)
	t.Logf("  95パーセンタイル: %v", p95Time)

	// パフォーマンス要件を確認
	assert.LessOrEqual(t, avgTime, 100*time.Millisecond, "平均レスポンス時間が100ms以下である必要があります")
	assert.LessOrEqual(t, p95Time, 200*time.Millisecond, "95パーセンタイルが200ms以下である必要があります")

	suite.tradeService.AssertExpectations(t)
}

// TestPerformance_StressTest はストレステスト
func TestPerformance_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("ストレステストはshortモードでスキップされます")
	}

	suite := SetupPerformanceTest(t)
	defer suite.TeardownPerformanceTest()

	duration := 30 * time.Second // 30秒間のストレステスト
	concurrentUsers := 20

	// 期待される残高レスポンス
	expectedBalance := &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}

	// 大量のモック期待値を設定（概算）
	estimatedRequests := int(duration.Seconds()) * concurrentUsers * 10 // 1秒あたり10リクエスト想定
	for i := 0; i < estimatedRequests; i++ {
		suite.tradeService.On("GetBalance", mock.Anything).Return(expectedBalance, nil).Maybe()
	}

	var wg sync.WaitGroup
	results := make(chan error, estimatedRequests)
	stopChan := make(chan struct{})

	startTime := time.Now()

	// 並行ユーザーを起動
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			requestCount := 0

			for {
				select {
				case <-stopChan:
					t.Logf("User %d completed %d requests", userID, requestCount)
					return
				default:
					resp, err := suite.client.Get(suite.server.URL + "/trade/balance")
					if err != nil {
						results <- fmt.Errorf("user %d: %w", userID, err)
					} else {
						if resp.StatusCode == http.StatusOK {
							results <- nil
						} else {
							results <- fmt.Errorf("user %d: status %d", userID, resp.StatusCode)
						}
						resp.Body.Close()
					}
					requestCount++

					// 少し待機してサーバーに負荷をかけすぎないようにする
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(i)
	}

	// 指定時間後に停止
	time.Sleep(duration)
	close(stopChan)
	wg.Wait()
	close(results)

	actualDuration := time.Since(startTime)

	// 結果を集計
	successCount := 0
	errorCount := 0
	for result := range results {
		if result == nil {
			successCount++
		} else {
			errorCount++
		}
	}

	totalRequests := successCount + errorCount
	requestsPerSecond := float64(totalRequests) / actualDuration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100

	t.Logf("ストレステスト結果:")
	t.Logf("  テスト時間: %v", actualDuration)
	t.Logf("  並行ユーザー数: %d", concurrentUsers)
	t.Logf("  総リクエスト数: %d", totalRequests)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  エラー数: %d", errorCount)
	t.Logf("  成功率: %.2f%%", successRate)
	t.Logf("  スループット: %.2f req/sec", requestsPerSecond)

	// ストレステストの要件を確認
	assert.GreaterOrEqual(t, successRate, 90.0, "ストレステストでの成功率が90%以上である必要があります")
	assert.GreaterOrEqual(t, requestsPerSecond, 50.0, "ストレステストでのスループットが50 req/sec以上である必要があります")
}
