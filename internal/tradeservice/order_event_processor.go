package tradeservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"stock-bot/domain/repository"
	"stock-bot/internal/app"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/client"
)

// OrderEventProcessor はGoaサービス用の注文イベント処理器
type OrderEventProcessor struct {
	executionHandler *eventprocessing.ExecutionEventHandlerImpl
	eventClient      client.EventClient
	session          *client.Session
	logger           *slog.Logger
	stopCh           chan struct{}
}

// NewOrderEventProcessor は新しい注文イベント処理器を作成する
func NewOrderEventProcessor(
	orderRepo repository.OrderRepository,
	positionRepo repository.PositionRepository,
	eventClient client.EventClient,
	session *client.Session,
	logger *slog.Logger,
) *OrderEventProcessor {
	// ExecutionUseCaseを作成
	executionUseCase := app.NewExecutionUseCaseImpl(orderRepo, positionRepo)

	// ExecutionEventHandlerを作成
	executionHandler := eventprocessing.NewExecutionEventHandler(executionUseCase, logger)

	return &OrderEventProcessor{
		executionHandler: executionHandler,
		eventClient:      eventClient,
		session:          session,
		logger:           logger,
		stopCh:           make(chan struct{}),
	}
}

// Start はWebSocketイベント処理を開始する
func (p *OrderEventProcessor) Start(ctx context.Context, symbols []string) error {
	p.logger.Info("Starting order event processor for WebSocket events")

	// WebSocket接続を開始
	messages, errs, err := p.eventClient.Connect(ctx, p.session, symbols)
	if err != nil {
		return err
	}

	go func() {
		defer p.logger.Info("Order event processor stopped")

		for {
			select {
			case message, ok := <-messages:
				if !ok {
					p.logger.Warn("Message channel closed")
					return
				}

				// メッセージを解析してイベント処理
				if err := p.processMessage(ctx, message); err != nil {
					p.logger.Error("Failed to process message", "error", err)
				}

			case err, ok := <-errs:
				if !ok {
					p.logger.Warn("Error channel closed")
					return
				}
				p.logger.Error("WebSocket error", "error", err)

			case <-p.stopCh:
				p.logger.Info("Received stop signal for order event processor")
				return
			case <-ctx.Done():
				p.logger.Info("Context cancelled, stopping order event processor")
				return
			}
		}
	}()

	return nil
}

// processMessage はWebSocketメッセージを解析してイベント処理する
func (p *OrderEventProcessor) processMessage(ctx context.Context, message []byte) error {
	// メッセージをJSONとして解析
	var eventData map[string]interface{}
	if err := json.Unmarshal(message, &eventData); err != nil {
		p.logger.Debug("Failed to parse message as JSON, skipping", "message", string(message))
		return nil // JSONでない場合はスキップ
	}

	// イベントタイプを取得
	eventTypeRaw, ok := eventData["p_event"]
	if !ok {
		return nil // イベントタイプがない場合はスキップ
	}

	eventType, ok := eventTypeRaw.(string)
	if !ok {
		return nil
	}

	// 約定通知（EC）イベントのみ処理
	if eventType == "EC" {
		// map[string]interface{} を map[string]string に変換
		stringData := make(map[string]string)
		for k, v := range eventData {
			if str, ok := v.(string); ok {
				stringData[k] = str
			}
		}

		if err := p.executionHandler.HandleEvent(ctx, eventType, stringData); err != nil {
			p.logger.Error("Failed to handle execution event", "error", err, "event_type", eventType)
			return err
		}

		p.logger.Debug("Successfully processed execution event", "event_type", eventType)
	}

	return nil
}

// Stop はWebSocketイベント処理を停止する
func (p *OrderEventProcessor) Stop() {
	p.logger.Info("Stopping order event processor")
	close(p.stopCh)
	if p.eventClient != nil {
		p.eventClient.Close()
	}
}
