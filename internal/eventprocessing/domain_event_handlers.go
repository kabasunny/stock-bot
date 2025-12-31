package eventprocessing

import (
	"context"
	"log/slog"
	"stock-bot/domain/event"
	"stock-bot/domain/repository"
)

// OrderEventHandler は注文関連のドメインイベントハンドラー
type OrderEventHandler struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

// NewOrderEventHandler は新しい注文イベントハンドラーを作成
func NewOrderEventHandler(orderRepo repository.OrderRepository, logger *slog.Logger) *OrderEventHandler {
	return &OrderEventHandler{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

// Handle はイベントを処理
func (h *OrderEventHandler) Handle(ctx context.Context, domainEvent event.DomainEvent) error {
	switch e := domainEvent.(type) {
	case *event.OrderPlacedEvent:
		return h.handleOrderPlaced(ctx, e)
	case *event.OrderExecutedEvent:
		return h.handleOrderExecuted(ctx, e)
	case *event.OrderCancelledEvent:
		return h.handleOrderCancelled(ctx, e)
	default:
		h.logger.Warn("unhandled event type in OrderEventHandler",
			slog.String("event_type", domainEvent.EventType()))
		return nil
	}
}

// CanHandle はハンドル可能なイベントタイプかどうかを判定
func (h *OrderEventHandler) CanHandle(eventType string) bool {
	switch eventType {
	case "OrderPlaced", "OrderExecuted", "OrderCancelled":
		return true
	default:
		return false
	}
}

func (h *OrderEventHandler) handleOrderPlaced(ctx context.Context, event *event.OrderPlacedEvent) error {
	h.logger.Info("handling order placed event",
		slog.String("order_id", event.Order.OrderID),
		slog.String("symbol", event.Order.Symbol))

	// 注文履歴の記録、通知送信などの処理
	// 実装は必要に応じて追加
	return nil
}

func (h *OrderEventHandler) handleOrderExecuted(ctx context.Context, event *event.OrderExecutedEvent) error {
	h.logger.Info("handling order executed event",
		slog.String("order_id", event.Execution.OrderID),
		slog.String("symbol", event.Execution.Symbol))

	// 約定処理、ポジション更新などの処理
	// 実装は必要に応じて追加
	return nil
}

func (h *OrderEventHandler) handleOrderCancelled(ctx context.Context, event *event.OrderCancelledEvent) error {
	h.logger.Info("handling order cancelled event",
		slog.String("order_id", event.OrderID),
		slog.String("reason", event.Reason))

	// キャンセル処理、通知送信などの処理
	// 実装は必要に応じて追加
	return nil
}

// PositionEventHandler はポジション関連のドメインイベントハンドラー
type PositionEventHandler struct {
	positionRepo repository.PositionRepository
	logger       *slog.Logger
}

// NewPositionEventHandler は新しいポジションイベントハンドラーを作成
func NewPositionEventHandler(positionRepo repository.PositionRepository, logger *slog.Logger) *PositionEventHandler {
	return &PositionEventHandler{
		positionRepo: positionRepo,
		logger:       logger,
	}
}

// Handle はイベントを処理
func (h *PositionEventHandler) Handle(ctx context.Context, domainEvent event.DomainEvent) error {
	switch e := domainEvent.(type) {
	case *event.PositionOpenedEvent:
		return h.handlePositionOpened(ctx, e)
	case *event.PositionClosedEvent:
		return h.handlePositionClosed(ctx, e)
	default:
		h.logger.Warn("unhandled event type in PositionEventHandler",
			slog.String("event_type", domainEvent.EventType()))
		return nil
	}
}

// CanHandle はハンドル可能なイベントタイプかどうかを判定
func (h *PositionEventHandler) CanHandle(eventType string) bool {
	switch eventType {
	case "PositionOpened", "PositionClosed":
		return true
	default:
		return false
	}
}

func (h *PositionEventHandler) handlePositionOpened(ctx context.Context, event *event.PositionOpenedEvent) error {
	h.logger.Info("handling position opened event",
		slog.String("symbol", event.Position.Symbol),
		slog.Int("quantity", event.Position.Quantity))

	// ポジション開始時の処理（リスク計算、通知など）
	// 実装は必要に応じて追加
	return nil
}

func (h *PositionEventHandler) handlePositionClosed(ctx context.Context, event *event.PositionClosedEvent) error {
	h.logger.Info("handling position closed event",
		slog.String("symbol", event.Position.Symbol),
		slog.Float64("realized_pl", event.RealizedPL))

	// ポジション決済時の処理（損益記録、通知など）
	// 実装は必要に応じて追加
	return nil
}

// RiskEventHandler はリスク関連のドメインイベントハンドラー
type RiskEventHandler struct {
	logger *slog.Logger
}

// NewRiskEventHandler は新しいリスクイベントハンドラーを作成
func NewRiskEventHandler(logger *slog.Logger) *RiskEventHandler {
	return &RiskEventHandler{
		logger: logger,
	}
}

// Handle はイベントを処理
func (h *RiskEventHandler) Handle(ctx context.Context, domainEvent event.DomainEvent) error {
	switch e := domainEvent.(type) {
	case *event.RiskLimitExceededEvent:
		return h.handleRiskLimitExceeded(ctx, e)
	default:
		h.logger.Warn("unhandled event type in RiskEventHandler",
			slog.String("event_type", domainEvent.EventType()))
		return nil
	}
}

// CanHandle はハンドル可能なイベントタイプかどうかを判定
func (h *RiskEventHandler) CanHandle(eventType string) bool {
	switch eventType {
	case "RiskLimitExceeded":
		return true
	default:
		return false
	}
}

func (h *RiskEventHandler) handleRiskLimitExceeded(ctx context.Context, event *event.RiskLimitExceededEvent) error {
	h.logger.Error("risk limit exceeded",
		slog.String("risk_type", event.RiskType),
		slog.Float64("current_value", event.CurrentValue),
		slog.Float64("limit", event.Limit),
		slog.String("symbol", event.Symbol))

	// リスク制限超過時の処理（緊急停止、通知など）
	// 実装は必要に応じて追加
	return nil
}
