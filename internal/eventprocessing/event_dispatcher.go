package eventprocessing

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/service"
	"sync"
)

// EventDispatcherImpl はイベントディスパッチャーの実装
type EventDispatcherImpl struct {
	handlers map[string]service.EventHandler
	mutex    sync.RWMutex
	logger   *slog.Logger
}

// NewEventDispatcher は新しいイベントディスパッチャーを作成する
func NewEventDispatcher(logger *slog.Logger) *EventDispatcherImpl {
	return &EventDispatcherImpl{
		handlers: make(map[string]service.EventHandler),
		logger:   logger,
	}
}

// RegisterHandler はイベントタイプに対するハンドラーを登録する
func (d *EventDispatcherImpl) RegisterHandler(eventType string, handler service.EventHandler) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.handlers[eventType] = handler
	d.logger.Info("event handler registered", "event_type", eventType)
}

// DispatchEvent はイベントを適切なハンドラーに振り分ける
func (d *EventDispatcherImpl) DispatchEvent(ctx context.Context, eventType string, data map[string]string) error {
	d.mutex.RLock()
	handler, exists := d.handlers[eventType]
	d.mutex.RUnlock()

	if !exists {
		d.logger.Warn("no handler registered for event type", "event_type", eventType)
		return fmt.Errorf("no handler registered for event type: %s", eventType)
	}

	d.logger.Debug("dispatching event", "event_type", eventType)
	return handler.HandleEvent(ctx, eventType, data)
}
