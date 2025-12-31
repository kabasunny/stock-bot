package agent

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTradeService はTradeServiceのモック実装
type MockTradeService struct {
	session *model.Session
}

func (m *MockTradeService) GetSession() *model.Session {
	return m.session
}

func (m *MockTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	return []*model.Position{}, nil
}

func (m *MockTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

func (m *MockTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	return &service.Balance{Cash: 1000000, BuyingPower: 800000}, nil
}

func (m *MockTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	return []*service.HistoricalPrice{}, nil
}

func (m *MockTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	return &model.Order{
		OrderID:     "test-order-123",
		Symbol:      req.Symbol,
		TradeType:   req.TradeType,
		OrderType:   req.OrderType,
		Quantity:    req.Quantity,
		Price:       req.Price,
		OrderStatus: model.OrderStatusNew,
	}, nil
}

func (m *MockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

func (m *MockTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	return &model.Order{OrderID: orderID}, nil
}

func (m *MockTradeService) CancelAllOrders(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *MockTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

func (m *MockTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	return &service.HealthStatus{
		Status:             "healthy",
		Timestamp:          time.Now(),
		SessionValid:       true,
		DatabaseConnected:  true,
		WebSocketConnected: true,
	}, nil
}

// MockEventClient はEventClientのモック実装
type MockEventClient struct{}

func (m *MockEventClient) Connect(ctx context.Context, session *client.Session, symbols []string) (<-chan []byte, <-chan error, error) {
	messages := make(chan []byte)
	errors := make(chan error)
	close(messages)
	close(errors)
	return messages, errors, nil
}

func (m *MockEventClient) Close() {
	// モック実装では何もしない
}

// MockPositionRepository はPositionRepositoryのモック実装
type MockPositionRepository struct{}

func (m *MockPositionRepository) Save(ctx context.Context, position *model.Position) error {
	return nil
}

func (m *MockPositionRepository) FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) {
	return nil, nil
}

func (m *MockPositionRepository) FindAll(ctx context.Context) ([]*model.Position, error) {
	return []*model.Position{}, nil
}

func (m *MockPositionRepository) UpdateHighestPrice(ctx context.Context, symbol string, price float64) error {
	return nil
}

func (m *MockPositionRepository) UpsertPositionByExecution(ctx context.Context, execution *model.Execution) error {
	return nil
}

func (m *MockPositionRepository) DeletePosition(ctx context.Context, symbol string, accountType model.PositionAccountType) error {
	return nil
}

// MockExecutionUseCase はExecutionUseCaseのモック実装
type MockExecutionUseCase struct{}

func (m *MockExecutionUseCase) Execute(ctx context.Context, execution *model.Execution) error {
	return nil
}

// TestNewAgent はNewAgent関数をテストします
func TestNewAgent(t *testing.T) {
	// 既存の設定ファイルを使用
	configPath := "../../agent_config.yaml"

	mockTradeService := &MockTradeService{
		session: &model.Session{
			SessionID: "test-session",
			UserID:    "test-user",
			IsActive:  true,
		},
	}
	mockEventClient := &MockEventClient{}
	mockPositionRepo := &MockPositionRepository{}
	mockExecutionUseCase := &MockExecutionUseCase{}

	agent, err := NewAgent(
		configPath,
		mockTradeService,
		mockEventClient,
		mockPositionRepo,
		mockExecutionUseCase,
	)

	require.NoError(t, err)
	require.NotNil(t, agent)
	assert.Equal(t, configPath, agent.configPath)
	assert.NotNil(t, agent.tradeService)
	assert.NotNil(t, agent.state)
	assert.NotNil(t, agent.webSocketEventService)
	assert.NotNil(t, agent.eventDispatcher)
}

// TestAgent_Tick はTick関数をテストします
func TestAgent_Tick(t *testing.T) {
	configPath := "../../agent_config.yaml"

	mockTradeService := &MockTradeService{
		session: &model.Session{
			SessionID: "test-session",
			UserID:    "test-user",
			IsActive:  true,
		},
	}
	mockEventClient := &MockEventClient{}
	mockPositionRepo := &MockPositionRepository{}
	mockExecutionUseCase := &MockExecutionUseCase{}

	agent, err := NewAgent(
		configPath,
		mockTradeService,
		mockEventClient,
		mockPositionRepo,
		mockExecutionUseCase,
	)

	require.NoError(t, err)

	// Tickメソッドを実行（エラーが発生しないことを確認）
	agent.Tick()
}

// TestAgent_Stop はStop関数をテストします
func TestAgent_Stop(t *testing.T) {
	configPath := "../../agent_config.yaml"

	mockTradeService := &MockTradeService{
		session: &model.Session{
			SessionID: "test-session",
			UserID:    "test-user",
			IsActive:  true,
		},
	}
	mockEventClient := &MockEventClient{}
	mockPositionRepo := &MockPositionRepository{}
	mockExecutionUseCase := &MockExecutionUseCase{}

	agent, err := NewAgent(
		configPath,
		mockTradeService,
		mockEventClient,
		mockPositionRepo,
		mockExecutionUseCase,
	)

	require.NoError(t, err)

	// Stopメソッドを実行（エラーが発生しないことを確認）
	agent.Stop()

	// コンテキストがキャンセルされていることを確認
	select {
	case <-agent.ctx.Done():
		// 期待通り
	default:
		t.Error("context should be cancelled after Stop()")
	}
}
