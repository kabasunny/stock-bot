package tests

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/domain/service"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/client"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventClient はEventClientのモック実装
type MockEventClient struct {
	mock.Mock
}

func (m *MockEventClient) Connect(ctx context.Context, session *client.Session, symbols []string) (<-chan []byte, <-chan error, error) {
	args := m.Called(ctx, session, symbols)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(<-chan []byte), args.Get(1).(<-chan error), args.Error(2)
}

func (m *MockEventClient) Close() {
	m.Called()
}

// MockEventDispatcher はEventDispatcherのモック実装
type MockEventDispatcher struct {
	mock.Mock
}

func (m *MockEventDispatcher) RegisterHandler(eventType string, handler service.EventHandler) {
	m.Called(eventType, handler)
}

func (m *MockEventDispatcher) DispatchEvent(ctx context.Context, eventType string, data map[string]string) error {
	args := m.Called(ctx, eventType, data)
	return args.Error(0)
}

// TestNewWebSocketEventService はWebSocketEventServiceの作成をテストします
func TestNewWebSocketEventService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mockEventClient := &MockEventClient{}
	mockEventDispatcher := &MockEventDispatcher{}

	service := eventprocessing.NewWebSocketEventService(
		mockEventClient,
		mockEventDispatcher,
		logger,
	)

	assert.NotNil(t, service, "WebSocketEventService should not be nil")
}

// TestWebSocketEventService_StartEventWatcher_ConnectionError は接続エラーのテストです
func TestWebSocketEventService_StartEventWatcher_ConnectionError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mockEventClient := &MockEventClient{}
	mockEventDispatcher := &MockEventDispatcher{}
	session := client.NewSession()

	service := eventprocessing.NewWebSocketEventService(
		mockEventClient,
		mockEventDispatcher,
		logger,
	)

	targetSymbols := []string{"1301"}
	session.EventURL = "wss://example.com/event"

	// 接続エラーを設定
	mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
		nil,
		nil,
		assert.AnError,
	).Once()

	ctx := context.Background()

	err := service.StartEventWatcher(ctx, session, targetSymbols)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to connect", "Should indicate connection failure")
	mockEventClient.AssertExpectations(t)
}

// TestWebSocketEventService_StartEventWatcher_Success は正常な接続開始をテストします
func TestWebSocketEventService_StartEventWatcher_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mockEventClient := &MockEventClient{}
	mockEventDispatcher := &MockEventDispatcher{}
	session := client.NewSession()

	service := eventprocessing.NewWebSocketEventService(
		mockEventClient,
		mockEventDispatcher,
		logger,
	)

	targetSymbols := []string{"1301", "1302"}
	session.EventURL = "wss://example.com/event"

	// モックチャネルを作成
	messages := make(chan []byte, 1)
	errs := make(chan error, 1)

	// モックの設定 - Closeは呼ばれない可能性があるのでMaybeを使用
	mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
		(<-chan []byte)(messages),
		(<-chan error)(errs),
		nil,
	).Once()

	ctx := context.Background()

	// イベント監視を開始
	err := service.StartEventWatcher(ctx, session, targetSymbols)

	assert.NoError(t, err)
	mockEventClient.AssertExpectations(t)
}

// TestWebSocketEventService_StartEventWatcher_EmptySymbols は空のシンボルリストでのテスト
func TestWebSocketEventService_StartEventWatcher_EmptySymbols(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mockEventClient := &MockEventClient{}
	mockEventDispatcher := &MockEventDispatcher{}
	session := client.NewSession()

	service := eventprocessing.NewWebSocketEventService(
		mockEventClient,
		mockEventDispatcher,
		logger,
	)

	targetSymbols := []string{} // 空のリスト
	session.EventURL = "wss://example.com/event"

	// モックチャネルを作成
	messages := make(chan []byte, 1)
	errs := make(chan error, 1)

	// モックの設定 - 空のシンボルでも接続は成功する
	mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
		(<-chan []byte)(messages),
		(<-chan error)(errs),
		nil,
	).Once()

	ctx := context.Background()

	// 空のシンボルリストでのテスト - 実装では警告ログを出すが成功する
	err := service.StartEventWatcher(ctx, session, targetSymbols)

	assert.NoError(t, err, "Should not return error for empty symbols")
	mockEventClient.AssertExpectations(t)
}
