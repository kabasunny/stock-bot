package event

import (
	"context"
	"log/slog"
	"sync"
)

// EventHandler はイベントハンドラーのインターフェース
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventPublisher はドメインイベントの発行者
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler EventHandler)
	Unsubscribe(eventType string, handler EventHandler)
	GetAllEventTypes() []string
}

// InMemoryEventPublisher はインメモリのイベント発行者実装
type InMemoryEventPublisher struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
	logger   *slog.Logger
}

// NewInMemoryEventPublisher は新しいインメモリイベント発行者を作成
func NewInMemoryEventPublisher(logger *slog.Logger) *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

// Publish はイベントを発行し、登録されたハンドラーに配信
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	p.mutex.RLock()
	handlers, exists := p.handlers[event.EventType()]
	p.mutex.RUnlock()

	if !exists {
		p.logger.Debug("no handlers registered for event type",
			slog.String("event_type", event.EventType()),
			slog.String("event_id", event.EventID()))
		return nil
	}

	p.logger.Info("publishing domain event",
		slog.String("event_type", event.EventType()),
		slog.String("event_id", event.EventID()),
		slog.String("aggregate_id", event.AggregateID()),
		slog.Int("handler_count", len(handlers)))

	// 各ハンドラーに並行してイベントを配信
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h.Handle(ctx, event); err != nil {
				p.logger.Error("event handler failed",
					slog.String("event_type", event.EventType()),
					slog.String("event_id", event.EventID()),
					slog.Any("error", err))
				errChan <- err
			}
		}(handler)
	}

	wg.Wait()
	close(errChan)

	// エラーがあった場合は最初のエラーを返す
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Subscribe はイベントタイプにハンドラーを登録
func (p *InMemoryEventPublisher) Subscribe(eventType string, handler EventHandler) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.handlers[eventType]; !exists {
		p.handlers[eventType] = make([]EventHandler, 0)
	}

	p.handlers[eventType] = append(p.handlers[eventType], handler)

	p.logger.Info("event handler subscribed",
		slog.String("event_type", eventType),
		slog.Int("total_handlers", len(p.handlers[eventType])))
}

// Unsubscribe はイベントタイプからハンドラーを削除
func (p *InMemoryEventPublisher) Unsubscribe(eventType string, handler EventHandler) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	handlers, exists := p.handlers[eventType]
	if !exists {
		return
	}

	// ハンドラーを検索して削除
	for i, h := range handlers {
		if h == handler {
			p.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			p.logger.Info("event handler unsubscribed",
				slog.String("event_type", eventType),
				slog.Int("remaining_handlers", len(p.handlers[eventType])))
			break
		}
	}

	// ハンドラーが空になった場合はマップから削除
	if len(p.handlers[eventType]) == 0 {
		delete(p.handlers, eventType)
	}
}

// GetHandlerCount は指定されたイベントタイプのハンドラー数を返す
func (p *InMemoryEventPublisher) GetHandlerCount(eventType string) int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if handlers, exists := p.handlers[eventType]; exists {
		return len(handlers)
	}
	return 0
}

// GetAllEventTypes は登録されている全イベントタイプを返す
func (p *InMemoryEventPublisher) GetAllEventTypes() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	eventTypes := make([]string, 0, len(p.handlers))
	for eventType := range p.handlers {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes
}
