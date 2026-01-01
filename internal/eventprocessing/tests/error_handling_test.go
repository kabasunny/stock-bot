package tests

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"stock-bot/domain/service"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/client"
	"testing"
	"time"

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

// TestWebSocketEventService_ErrorHandling はWebSocketイベントサービスのエラーハンドリングテスト
func TestWebSocketEventService_ErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// 接続エラーハンドリングテスト
	t.Run("ConnectionError", func(t *testing.T) {
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
			fmt.Errorf("connection failed: network unreachable"),
		).Once()

		ctx := context.Background()

		err := service.StartEventWatcher(ctx, session, targetSymbols)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect")
		assert.Contains(t, err.Error(), "network unreachable")
		mockEventClient.AssertExpectations(t)
	})

	// メッセージ処理エラーテスト
	t.Run("MessageProcessingError", func(t *testing.T) {
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

		// モックチャネルを作成
		messages := make(chan []byte, 1)
		errs := make(chan error, 1)

		// 不正なメッセージフォーマット
		invalidMessage := []byte("invalid message format")

		// モックの設定
		mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
			(<-chan []byte)(messages),
			(<-chan error)(errs),
			nil,
		).Once()
		mockEventClient.On("Close").Return().Maybe()

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// イベント監視を開始
		go func() {
			// 少し待ってから不正なメッセージを送信
			time.Sleep(50 * time.Millisecond)
			messages <- invalidMessage
			// 処理時間を待つ
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := service.StartEventWatcher(ctx, session, targetSymbols)

		// エラーが発生してもサービスは継続することを確認
		assert.NoError(t, err)
		mockEventClient.AssertExpectations(t)
	})

	// ディスパッチャーエラーテスト
	t.Run("DispatcherError", func(t *testing.T) {
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

		// モックチャネルを作成
		messages := make(chan []byte, 1)
		errs := make(chan error, 1)

		// 有効なメッセージフォーマット
		validMessage := []byte("p_cmd\x02EC\x01p_symbol\x021301\x01p_price\x021500")

		// モックの設定
		mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
			(<-chan []byte)(messages),
			(<-chan error)(errs),
			nil,
		).Once()
		mockEventDispatcher.On("DispatchEvent", mock.Anything, "EC", mock.AnythingOfType("map[string]string")).Return(
			fmt.Errorf("dispatcher error: handler not found"),
		).Once()
		mockEventClient.On("Close").Return().Maybe()

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// イベント監視を開始
		go func() {
			// 少し待ってからメッセージを送信
			time.Sleep(50 * time.Millisecond)
			messages <- validMessage
			// 処理時間を待つ
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := service.StartEventWatcher(ctx, session, targetSymbols)

		// ディスパッチャーエラーが発生してもサービスは継続することを確認
		assert.NoError(t, err)
		mockEventClient.AssertExpectations(t)
		mockEventDispatcher.AssertExpectations(t)
	})

	// チャネルクローズエラーテスト
	t.Run("ChannelCloseError", func(t *testing.T) {
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

		// モックチャネルを作成
		messages := make(chan []byte, 1)
		errs := make(chan error, 1)

		// モックの設定
		mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
			(<-chan []byte)(messages),
			(<-chan error)(errs),
			nil,
		).Once()
		mockEventClient.On("Close").Return().Maybe()

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// イベント監視を開始
		go func() {
			// 少し待ってからチャネルをクローズ
			time.Sleep(50 * time.Millisecond)
			close(messages)
			close(errs)
		}()

		err := service.StartEventWatcher(ctx, session, targetSymbols)

		// チャネルクローズは正常な終了として扱われることを確認
		assert.NoError(t, err)
		mockEventClient.AssertExpectations(t)
	})

	// 複数エラーの同時発生テスト
	t.Run("MultipleSimultaneousErrors", func(t *testing.T) {
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

		// モックチャネルを作成
		messages := make(chan []byte, 2)
		errs := make(chan error, 2)

		// モックの設定
		mockEventClient.On("Connect", mock.Anything, session, targetSymbols).Return(
			(<-chan []byte)(messages),
			(<-chan error)(errs),
			nil,
		).Once()
		mockEventClient.On("Close").Return().Maybe()

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		// イベント監視を開始
		go func() {
			// 複数のエラーを同時に送信
			time.Sleep(50 * time.Millisecond)
			errs <- fmt.Errorf("network error 1")
			messages <- []byte("invalid message 1")

			time.Sleep(50 * time.Millisecond)
			errs <- fmt.Errorf("network error 2")
			messages <- []byte("invalid message 2")

			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := service.StartEventWatcher(ctx, session, targetSymbols)

		// 複数エラーが発生してもサービスが適切に処理することを確認
		assert.NoError(t, err)
		mockEventClient.AssertExpectations(t)
	})
}

// TestEventDispatcher_ErrorHandling はEventDispatcherのエラーハンドリングテスト
func TestEventDispatcher_ErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// 未登録イベントタイプのテスト
	t.Run("UnregisteredEventType", func(t *testing.T) {
		dispatcher := eventprocessing.NewEventDispatcher(logger)

		ctx := context.Background()
		eventData := map[string]string{
			"p_symbol": "1301",
			"p_price":  "1500",
		}

		// 未登録のイベントタイプでディスパッチ
		err := dispatcher.DispatchEvent(ctx, "UNKNOWN_EVENT", eventData)

		// 未登録イベントはエラーにならない（ログ出力のみ）
		assert.NoError(t, err)
	})

	// ハンドラーエラーのテスト
	t.Run("HandlerError", func(t *testing.T) {
		dispatcher := eventprocessing.NewEventDispatcher(logger)

		// エラーを返すモックハンドラーを作成
		mockHandler := &MockEventHandler{}
		mockHandler.On("HandleEvent", mock.Anything, "EC", mock.AnythingOfType("map[string]string")).Return(
			fmt.Errorf("handler processing error"),
		).Once()

		// ハンドラーを登録
		dispatcher.RegisterHandler("EC", mockHandler)

		ctx := context.Background()
		eventData := map[string]string{
			"p_symbol": "1301",
			"p_price":  "1500",
		}

		// エラーを返すハンドラーでディスパッチ
		err := dispatcher.DispatchEvent(ctx, "EC", eventData)

		// ハンドラーエラーが伝播されることを確認
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "handler processing error")
		mockHandler.AssertExpectations(t)
	})

	// nilハンドラーのテスト
	t.Run("NilHandler", func(t *testing.T) {
		dispatcher := eventprocessing.NewEventDispatcher(logger)

		// nilハンドラーを登録
		dispatcher.RegisterHandler("EC", nil)

		ctx := context.Background()
		eventData := map[string]string{
			"p_symbol": "1301",
			"p_price":  "1500",
		}

		// nilハンドラーでディスパッチ
		err := dispatcher.DispatchEvent(ctx, "EC", eventData)

		// nilハンドラーはエラーにならない（スキップされる）
		assert.NoError(t, err)
	})
}

// MockEventHandler はEventHandlerのモック実装
type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) HandleEvent(ctx context.Context, eventType string, data map[string]string) error {
	args := m.Called(ctx, eventType, data)
	return args.Error(0)
}
