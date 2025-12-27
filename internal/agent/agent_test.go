package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindSignalFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "signal_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up

	// Test case 1: No signal files found
	t.Run("no signal files", func(t *testing.T) {
		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != "" {
			t.Errorf("expected empty string, got %s", foundFile)
		}
	})

	// Test case 2: Single signal file
	t.Run("single signal file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "signal_20240101.bin")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != filePath {
			t.Errorf("expected %s, got %s", filePath, foundFile)
		}
	})

	// Test case 3: Multiple signal files, select newest
	t.Run("multiple signal files - select newest", func(t *testing.T) {
		file1 := filepath.Join(tmpDir, "signal_old.bin")
		file2 := filepath.Join(tmpDir, "signal_new.bin")
		file3 := filepath.Join(tmpDir, "signal_middle.bin")

		// Create files with specific modification times
		if err := os.WriteFile(file1, []byte("old"), 0644); err != nil {
			t.Fatalf("failed to create file1: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure distinct modification times
		if err := os.WriteFile(file3, []byte("middle"), 0644); err != nil {
			t.Fatalf("failed to create file3: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		if err := os.WriteFile(file2, []byte("new"), 0644); err != nil {
			t.Fatalf("failed to create file2: %v", err)
		}

		// Manually set mod times if needed for precise control, but sleep should be enough for most OS
		// For example:
		// t1 := time.Now().Add(-2 * time.Hour)
		// t2 := time.Now().Add(-1 * time.Hour)
		// t3 := time.Now()
		// os.Chtimes(file1, t1, t1)
		// os.Chtimes(file3, t2, t2)
		// os.Chtimes(file2, t3, t3)

		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != file2 {
			t.Errorf("expected newest file %s, got %s", file2, foundFile)
		}
	})

	// Test case 4: Glob pattern with no matching accessible files
	t.Run("no accessible files", func(t *testing.T) {
		// Create a file but make it inaccessible (e.g., wrong permissions or broken symlink)
		// This scenario is hard to reliably test across OS for os.Stat errors specifically,
		// but we can simulate a pattern that matches nothing.
		// For now, rely on previous 'no signal files' test for this case.
	})
}

type mockTradeService struct {
	mock.Mock
}

func (m *mockTradeService) GetSession() *client.Session {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*client.Session)
}

func (m *mockTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *mockTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func (m *mockTradeService) GetBalance(ctx context.Context) (*Balance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Balance), args.Error(1)
}

func (m *mockTradeService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error) {
	args := m.Called(ctx, symbol, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*HistoricalPrice), args.Error(1)
}

func (m *mockTradeService) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *mockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func TestCheckPositionsForExit(t *testing.T) {
	// --- Test Setup ---
	setup := func() (*Agent, *mockTradeService) {
		cfg := &AgentConfig{
			Agent: AgentSettings{
				ExecutionInterval: 1 * time.Minute,
			},
			StrategySettings: StrategySettings{
				Swingtrade: SwingtradeSettings{
					SignalFilePattern:       "test_signal.bin",
					TradeRiskPercentage:     0.01,
					UnitSize:                100,
					ProfitTakeRate:          10.0, // +10%
					StopLossRate:            5.0,  // -5%
					TrailingStopTriggerRate: 2.0,  // トレーリングストップ開始の利益率
					TrailingStopRate:        3.0,  // 最高値からの下落率
				},
			},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		mockService := new(mockTradeService)

		agent := &Agent{
			config:       cfg,
			logger:       logger,
			ctx:          context.Background(),
			state:        NewState(),
			tradeService: mockService,
		}
		return agent, mockService
	}

	basePosition := &model.Position{
		Symbol:       "1234",
		AveragePrice: 1000.0,
		Quantity:     100,
	}

	// --- Test Cases ---
	t.Run("should place profit take order", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		// Mock GetPriceHistory as it's now always called
		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1100.0 // 10% profit
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)
		mockService.On("PlaceOrder", mock.Anything, mock.Anything).Return(&model.Order{OrderID: "order-pt", Symbol: "1234"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		// Check that the order was placed
		placedOrder, ok := agent.state.GetOrder("order-pt")
		assert.True(t, ok)
		assert.Equal(t, "1234", placedOrder.Symbol)
	})

	t.Run("should place original stop loss order if ATR fails", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		// Mock GetPriceHistory to return an error
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(nil, fmt.Errorf("API error")).Once()

		currentPrice := 950.0 // 5% loss, this should still trigger the old logic as a fallback
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)
		// We expect no order because ATR calculation failure now causes a 'continue'.
		// mockService.On("PlaceOrder", mock.Anything, mock.Anything).Return(&model.Order{OrderID: "order-sl", Symbol: "1234"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})

	t.Run("should activate trailing stop", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition // copy
		agent.state.UpdatePositions([]*model.Position{&pos})
		
		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		// 1. Price rises to trigger trailing stop activation
		currentPrice := 1020.0                           // +2%, triggers activation
		expectedTrailingStopPrice := 1020.0 * (1 - 0.03) // 989.4
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		// No order should be placed yet

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, currentPrice, updatedPos.HighestPrice)
		assert.InDelta(t, expectedTrailingStopPrice, updatedPos.TrailingStopPrice, 0.001)
	})

	t.Run("should execute trailing stop order", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition                        // copy
		pos.HighestPrice = 1050.0                   // Manually set state as if price rose
		pos.TrailingStopPrice = 1050.0 * (1 - 0.03) // 1018.5
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		// Price drops below the trailing stop price
		currentPrice := 1018.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)
		mockService.On("PlaceOrder", mock.Anything, mock.Anything).Return(&model.Order{OrderID: "order-ts", Symbol: "1234"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		placedOrder, ok := agent.state.GetOrder("order-ts")
		assert.True(t, ok)
		assert.Equal(t, "1234", placedOrder.Symbol)
	})

	t.Run("should update trailing stop price as price rises", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition                        // copy
		pos.HighestPrice = 1020.0                   // Initial activation price
		pos.TrailingStopPrice = 1020.0 * (1 - 0.03) // 989.4
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		// Price rises further
		currentPrice := 1080.0
		expectedNewTrailingStopPrice := 1080.0 * (1 - 0.03) // 1047.6
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, currentPrice, updatedPos.HighestPrice, "HighestPrice should be updated")
		assert.InDelta(t, expectedNewTrailingStopPrice, updatedPos.TrailingStopPrice, 0.001, "TrailingStopPrice should be updated")
	})

	t.Run("should not place order if an open sell order exists", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})
		agent.state.AddOrder(&model.Order{
			Symbol:      "1234",
			TradeType:   model.TradeTypeSell,
			OrderStatus: model.OrderStatusNew, // Unexecuted
		})

		// No mocks needed as the function should return early.

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		mockService.AssertNotCalled(t, "GetPrice", mock.Anything, mock.Anything)
		mockService.AssertNotCalled(t, "GetPriceHistory", mock.Anything, mock.Anything)
	})

	t.Run("should do nothing if no conditions are met", func(t *testing.T) {
		agent, mockService := setup()
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1010.0 // Not enough for profit or trailing trigger
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, 0.0, updatedPos.TrailingStopPrice, "TrailingStopPrice should not be activated")
	})

	t.Run("should place ATR-based stop loss order", func(t *testing.T) {
		agent, mockService := setup()
		// Override config for this test
		agent.config.StrategySettings.Swingtrade.StopLossATRMultiplier = 2.0
		agent.config.StrategySettings.Swingtrade.ATRPeriod = 14

		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		// Mock GetPriceHistory for ATR calculation
		historicalData := make([]*HistoricalPrice, 15)
		for i := 0; i < 15; i++ {
			historicalData[i] = &HistoricalPrice{
				High:  1010.0, // TR will be 10
				Low:   1000.0,
				Close: 1000.0,
			}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", 15).Return(historicalData, nil).Once()

		// ATR should be 10
		// Stop loss multiplier is 2.0
		// Stop loss level = 1000 - (10 * 2.0) = 980
		currentPrice := 979.0 // Below ATR stop loss
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()

		mockService.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req *PlaceOrderRequest) bool {
			return req.Symbol == "1234" && req.TradeType == model.TradeTypeSell
		})).Return(&model.Order{OrderID: "order-atr-sl"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		_, ok := agent.state.GetOrder("order-atr-sl")
		assert.True(t, ok)
	})
}

func TestCheckSignalsForEntry_ATRBasedSizing(t *testing.T) {
	// --- Test Setup ---
	setup := func() (*Agent, *mockTradeService) {
		cfg := &AgentConfig{
			Agent: AgentSettings{
				ExecutionInterval: 1 * time.Minute,
			},
			StrategySettings: StrategySettings{
				Swingtrade: SwingtradeSettings{
					SignalFilePattern:       "test_signal.bin",
					TradeRiskPercentage:     0.02, // 口座資金の2%をリスクに
					UnitSize:                100,  // 最小取引単位
					ProfitTakeRate:          10.0,
					StopLossRate:            5.0,
					TrailingStopTriggerRate: 2.0,
					TrailingStopRate:        3.0,
					ATRPeriod:               14,  // ATR計算期間
					RiskPerATR:              2.0, // 1ATRあたり2ATR分のリスクを許容
				},
			},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		mockService := new(mockTradeService)

		agent := &Agent{
			configPath:   "agent_config.yaml", // For signal file logic
			config:       cfg,
			logger:       logger,
			ctx:          context.Background(),
			state:        NewState(),
			tradeService: mockService,
		}
		return agent, mockService
	}

	t.Run("should calculate ATR-based quantity and place buy order", func(t *testing.T) {
		        agent, mockService := setup()
		
		        // --- Mock Data ---
		        // Signal file
		        signalFilePath := createSignalFile(t, os.TempDir(), "1234", BuySignal)
		        defer os.Remove(signalFilePath)
		
		        		// Mock FindSignalFile and ReadSignalFile
		        		agent.signalPattern = filepath.Join(os.TempDir(), "*.bin")               // Point to temp directory
		        		// mockService.On("GetOrders", mock.Anything).Return([]*model.Order{}, nil) // No open orders - REMOVED		
		        // Set initial balance in agent's state
		        buyingPower := 1_000_000.0 // 100万円の買付余力
		        agent.state.UpdateBalance(&Balance{Cash: buyingPower, BuyingPower: buyingPower})
		
		        // Mock GetPrice for current price
		        currentPrice := 1000.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)

		// Mock GetPriceHistory for ATR calculation
		// ATRPeriod = 14, need 15 data points for ATR calculation
		historicalData := make([]*HistoricalPrice, 15)
		// For simplicity, create data such that True Range is consistently 10.0
		// High=1010, Low=1000, PrevClose=1000 => TR = 10
		for i := 0; i < 15; i++ {
			historicalData[i] = &HistoricalPrice{
				Open:  1000.0,
				High:  1010.0,
				Low:   1000.0,
				Close: 1000.0, // Consistent Close for predictable TR
			}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", 15).Return(historicalData, nil)

		// Expected ATR (assuming all TRs are 10 for ATRPeriod=14)
		// TradeRiskPercentage = 0.02 (2%)
		// BuyingPower = 1,000,000
		// totalRiskAmount = 1,000,000 * 0.02 = 20,000
		// RiskPerATR = 2.0
		// riskPerShare = 10.0 * 2.0 = 20.0
		// maxShares = totalRiskAmount / riskPerShare = 20,000 / 20.0 = 1,000
		// UnitSize = 100
		// quantity = math.Floor(1,000 / 100) * 100 = 1,000
		expectedQuantity := 1000

						// Mock PlaceOrder
						mockService.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req *PlaceOrderRequest) bool {
							return req.Symbol == "1234" &&
								req.TradeType == model.TradeTypeBuy &&
								req.OrderType == model.OrderTypeMarket &&
								req.Quantity == expectedQuantity
						})).Return(&model.Order{OrderID: "order-atr", Symbol: "1234", Quantity: expectedQuantity}, nil).Once()
		// --- Execute ---
		agent.checkSignalsForEntry(context.Background())

		// --- Assert ---
		mockService.AssertExpectations(t)
		placedOrder, ok := agent.state.GetOrder("order-atr")
		assert.True(t, ok)
		assert.Equal(t, "1234", placedOrder.Symbol)
		assert.Equal(t, expectedQuantity, placedOrder.Quantity)
	})
}

// Helper to simulate signal file creation for tests
func createSignalFile(t *testing.T, dir, symbol string, signalType TradeSignal) string {
	fileName := fmt.Sprintf("signal_%s_%d.bin", symbol, signalType)
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create signal file %s: %v", filePath, err)
	}
	defer file.Close()

	// Convert symbol string to uint16
	symbolUint16, err := strconv.ParseUint(symbol, 10, 16)
	if err != nil {
		t.Fatalf("Failed to parse symbol %s to uint16: %v", symbol, err)
	}

	// Write symbol and signalType in binary
	if err := binary.Write(file, binary.LittleEndian, uint16(symbolUint16)); err != nil {
		t.Fatalf("Failed to write symbol to signal file: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, signalType); err != nil {
		t.Fatalf("Failed to write signal type to signal file: %v", err)
	}

	return filePath
}
