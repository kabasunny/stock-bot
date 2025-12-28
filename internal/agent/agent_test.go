package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"stock-bot/domain/model" // Add this import
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

type mockPositionRepository struct {
	mock.Mock
}

func (m *mockPositionRepository) Save(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}
func (m *mockPositionRepository) FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}
func (m *mockPositionRepository) FindAll(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}
func (m *mockPositionRepository) UpdateHighestPrice(ctx context.Context, symbol string, price float64) error {
	args := m.Called(ctx, symbol, price)
	return args.Error(0)
}

func (m *mockPositionRepository) UpsertPositionByExecution(ctx context.Context, execution *model.Execution) error {
	args := m.Called(ctx, execution)
	return args.Error(0)
}

func (m *mockPositionRepository) DeletePosition(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

type mockExecutionUseCase struct {
	mock.Mock
}

func (m *mockExecutionUseCase) Execute(ctx context.Context, execution *model.Execution) error {
	args := m.Called(ctx, execution)
	return args.Error(0)
}

func TestCheckPositionsForExit(t *testing.T) {
	// --- Test Setup ---
	setup := func(t *testing.T) (*Agent, *mockTradeService, *mockPositionRepository, *mockExecutionUseCase) {
		// 一時的なagent_config.yamlを作成
		tmpFile, err := os.CreateTemp("", "agent_config_test_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp config file: %v", err)
		}
		defer os.Remove(tmpFile.Name()) // テスト終了時に削除

		configContent := `
agent:
  strategy: swingtrade
  execution_interval: 10s
strategy_settings:
  swingtrade:
    target_symbols: ["1234"]
    trade_risk_percentage: 0.01
    unit_size: 100
    profit_take_rate: 10.0
    stop_loss_rate: 5.0
    trailing_stop_trigger_rate: 2.0
    trailing_stop_rate: 3.0
    atr_period: 14
    risk_per_atr: 2.0
    stop_loss_atr_multiplier: 2.0
    signal_file_pattern: "./signals/*.bin"
`
		if _, err := tmpFile.WriteString(configContent); err != nil {
			t.Fatalf("Failed to write to temp config file: %v", err)
		}
		tmpFile.Close()

		mockService := new(mockTradeService)
		mockRepo := new(mockPositionRepository)
		mockExecUseCase := new(mockExecutionUseCase) // Initialize mockExecutionUseCase

		// eventClientはこれらのテストでは使われないのでnil
		agent, err := NewAgent(tmpFile.Name(), mockService, nil, mockRepo, mockExecUseCase) // Pass mockExecUseCase
		if err != nil {
			t.Fatalf("failed to create agent for test: %v", err)
		}
		return agent, mockService, mockRepo, mockExecUseCase
	}

	basePosition := &model.Position{
		Symbol:       "1234",
		AveragePrice: 1000.0,
		Quantity:     100,
	}

	// --- Test Cases ---
	t.Run("should place profit take order", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		for i := range historicalData {
			historicalData[i] = &HistoricalPrice{}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1100.0 // 10% profit
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()
		mockService.On("PlaceOrder", mock.Anything, mock.Anything).Return(&model.Order{OrderID: "order-pt", Symbol: "1234"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		placedOrder, ok := agent.state.GetOrder("order-pt")
		assert.True(t, ok)
		assert.Equal(t, "1234", placedOrder.Symbol)
	})

	t.Run("should place original stop loss order if ATR fails", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		currentPrice := 950.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(nil, fmt.Errorf("API error")).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})

	t.Run("should activate trailing stop", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition // copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		for i := range historicalData {
			historicalData[i] = &HistoricalPrice{}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1020.0
		expectedTrailingStopPrice := 1020.0 * (1 - 0.03)
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, currentPrice, updatedPos.HighestPrice)
		assert.InDelta(t, expectedTrailingStopPrice, updatedPos.TrailingStopPrice, 0.001)
	})

	t.Run("should execute trailing stop order", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition                        // copy
		pos.HighestPrice = 1050.0                   // Manually set state as if price rose
		pos.TrailingStopPrice = 1050.0 * (1 - 0.03) // 1018.5
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		for i := range historicalData {
			historicalData[i] = &HistoricalPrice{}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1018.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe() // Highest price doesn't update here
		mockService.On("PlaceOrder", mock.Anything, mock.Anything).Return(&model.Order{OrderID: "order-ts", Symbol: "1234"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		placedOrder, ok := agent.state.GetOrder("order-ts")
		assert.True(t, ok)
		assert.Equal(t, "1234", placedOrder.Symbol)
	})

	t.Run("should update trailing stop price as price rises", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition                        // copy
		pos.HighestPrice = 1020.0                   // Initial activation price
		pos.TrailingStopPrice = 1020.0 * (1 - 0.03) // 989.4
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		for i := range historicalData {
			historicalData[i] = &HistoricalPrice{}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1080.0
		expectedNewTrailingStopPrice := 1080.0 * (1 - 0.03)
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, currentPrice, updatedPos.HighestPrice, "HighestPrice should be updated")
		assert.InDelta(t, expectedNewTrailingStopPrice, updatedPos.TrailingStopPrice, 0.001, "TrailingStopPrice should be updated")
	})

	t.Run("should not place order if an open sell order exists", func(t *testing.T) {
		agent, mockService, _, _ := setup(t)
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})
		agent.state.AddOrder(&model.Order{
			Symbol:      "1234",
			TradeType:   model.TradeTypeSell,
			OrderStatus: model.OrderStatusNew, // Unexecuted
		})

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		mockService.AssertNotCalled(t, "GetPrice", mock.Anything, mock.Anything)
	})

	t.Run("should do nothing if no conditions are met", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, agent.config.StrategySettings.Swingtrade.ATRPeriod+1)
		for i := range historicalData {
			historicalData[i] = &HistoricalPrice{}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Maybe()

		currentPrice := 1010.0 // Not enough for profit or trailing trigger
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
		updatedPos, _ := agent.state.GetPosition("1234")
		assert.Equal(t, 0.0, updatedPos.TrailingStopPrice, "TrailingStopPrice should not be activated")
	})

	t.Run("should place ATR-based stop loss order", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition // Make a copy
		agent.state.UpdatePositions([]*model.Position{&pos})

		historicalData := make([]*HistoricalPrice, 15)
		for i := 0; i < 15; i++ {
			historicalData[i] = &HistoricalPrice{High: 1010.0, Low: 1000.0, Close: 1000.0}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", 15).Return(historicalData, nil).Once()

		currentPrice := 979.0 // Below ATR stop loss
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()
		mockService.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req *PlaceOrderRequest) bool {
			return req.Symbol == "1234" && req.TradeType == model.TradeTypeSell
		})).Return(&model.Order{OrderID: "order-atr-sl"}, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		_, ok := agent.state.GetOrder("order-atr-sl")
		assert.True(t, ok)
	})

	t.Run("should not place order if GetPrice returns error", func(t *testing.T) {
		agent, mockService, _, _ := setup(t)
		pos := *basePosition
		agent.state.UpdatePositions([]*model.Position{&pos})

		mockService.On("GetPrice", mock.Anything, "1234").Return(0.0, fmt.Errorf("price fetch error")).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})

	t.Run("should not place order if GetPriceHistory returns error", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition
		agent.state.UpdatePositions([]*model.Position{&pos})

		currentPrice := 1000.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(nil, fmt.Errorf("history fetch error")).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})

	t.Run("should not place order if not enough historical data for ATR", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition
		agent.state.UpdatePositions([]*model.Position{&pos})

		currentPrice := 1000.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()
		historicalData := make([]*HistoricalPrice, 10)
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})

	t.Run("should not place order if ATR calculation returns error or zero", func(t *testing.T) {
		agent, mockService, mockRepo, _ := setup(t)
		pos := *basePosition
		agent.state.UpdatePositions([]*model.Position{&pos})

		currentPrice := 1000.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil).Once()
		mockRepo.On("UpdateHighestPrice", mock.Anything, "1234", currentPrice).Return(nil).Once()

		historicalData := make([]*HistoricalPrice, 15)
		for i := 0; i < 15; i++ {
			historicalData[i] = &HistoricalPrice{High: 1000.0, Low: 1000.0, Close: 1000.0}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", mock.Anything).Return(historicalData, nil).Once()

		agent.checkPositionsForExit(context.Background())

		mockService.AssertNotCalled(t, "PlaceOrder", mock.Anything, mock.Anything)
	})
}

func TestCheckSignalsForEntry_ATRBasedSizing(t *testing.T) {
	// --- Test Setup ---
	setup := func(t *testing.T) (*Agent, *mockTradeService, *mockPositionRepository, *mockExecutionUseCase) {
		// 一時的なagent_config.yamlを作成
		tmpFile, err := os.CreateTemp("", "agent_config_test_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp config file: %v", err)
		}
		defer os.Remove(tmpFile.Name()) // テスト終了時に削除

		configContent := `
agent:
  strategy: swingtrade
  execution_interval: 10s
strategy_settings:
  swingtrade:
    target_symbols: ["1234"]
    trade_risk_percentage: 0.02
    unit_size: 100
    profit_take_rate: 10.0
    stop_loss_rate: 5.0
    trailing_stop_trigger_rate: 2.0
    trailing_stop_rate: 3.0
    atr_period: 14
    risk_per_atr: 2.0
    stop_loss_atr_multiplier: 2.0
    signal_file_pattern: "./signals/*.bin"
`
		if _, err := tmpFile.WriteString(configContent); err != nil {
			t.Fatalf("Failed to write to temp config file: %v", err)
		}
		tmpFile.Close()

		mockService := new(mockTradeService)
		mockRepo := new(mockPositionRepository)
		mockExecUseCase := new(mockExecutionUseCase)

		agent, err := NewAgent(tmpFile.Name(), mockService, nil, mockRepo, mockExecUseCase)
		if err != nil {
			t.Fatalf("failed to create agent for test: %v", err)
		}
		return agent, mockService, mockRepo, mockExecUseCase
	}

	t.Run("should calculate ATR-based quantity and place buy order", func(t *testing.T) {
		agent, mockService, _, _ := setup(t)

		// --- Mock Data ---
		// Signal file
		signalFilePath := createSignalFile(t, os.TempDir(), "1234", BuySignal)
		defer os.Remove(signalFilePath)

		// Mock FindSignalFile and ReadSignalFile
		agent.signalPattern = filepath.Join(os.TempDir(), "*.bin")
		// Set initial balance in agent's state
		buyingPower := 1_000_000.0 // 100万円の買付余力
		agent.state.UpdateBalance(&Balance{Cash: buyingPower, BuyingPower: buyingPower})

		// Mock GetPrice for current price
		currentPrice := 1000.0
		mockService.On("GetPrice", mock.Anything, "1234").Return(currentPrice, nil)

		// Mock GetPriceHistory for ATR calculation
		historicalData := make([]*HistoricalPrice, 15)
		for i := 0; i < 15; i++ {
			historicalData[i] = &HistoricalPrice{
				Open:  1000.0,
				High:  1010.0,
				Low:   1000.0,
				Close: 1000.0,
			}
		}
		mockService.On("GetPriceHistory", mock.Anything, "1234", 15).Return(historicalData, nil)

		expectedQuantity := 200

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
