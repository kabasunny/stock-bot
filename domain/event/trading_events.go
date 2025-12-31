package event

import (
	"stock-bot/domain/model"
)

// 取引関連のドメインイベント

// OrderPlacedEvent は注文発行イベント
type OrderPlacedEvent struct {
	BaseDomainEvent
	Order *model.Order `json:"order"`
}

// NewOrderPlacedEvent は新しい注文発行イベントを作成
func NewOrderPlacedEvent(order *model.Order) *OrderPlacedEvent {
	return &OrderPlacedEvent{
		BaseDomainEvent: NewBaseDomainEvent("OrderPlaced", order.OrderID, 1),
		Order:           order,
	}
}

// OrderExecutedEvent は注文約定イベント
type OrderExecutedEvent struct {
	BaseDomainEvent
	Execution *model.Execution `json:"execution"`
}

// NewOrderExecutedEvent は新しい注文約定イベントを作成
func NewOrderExecutedEvent(execution *model.Execution) *OrderExecutedEvent {
	return &OrderExecutedEvent{
		BaseDomainEvent: NewBaseDomainEvent("OrderExecuted", execution.OrderID, 1),
		Execution:       execution,
	}
}

// OrderCancelledEvent は注文キャンセルイベント
type OrderCancelledEvent struct {
	BaseDomainEvent
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}

// NewOrderCancelledEvent は新しい注文キャンセルイベントを作成
func NewOrderCancelledEvent(orderID, reason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		BaseDomainEvent: NewBaseDomainEvent("OrderCancelled", orderID, 1),
		OrderID:         orderID,
		Reason:          reason,
	}
}

// PositionOpenedEvent はポジション開始イベント
type PositionOpenedEvent struct {
	BaseDomainEvent
	Position *model.Position `json:"position"`
}

// NewPositionOpenedEvent は新しいポジション開始イベントを作成
func NewPositionOpenedEvent(position *model.Position) *PositionOpenedEvent {
	return &PositionOpenedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PositionOpened", position.Symbol, 1),
		Position:        position,
	}
}

// PositionClosedEvent はポジション決済イベント
type PositionClosedEvent struct {
	BaseDomainEvent
	Position   *model.Position `json:"position"`
	RealizedPL float64         `json:"realized_pl"`
}

// NewPositionClosedEvent は新しいポジション決済イベントを作成
func NewPositionClosedEvent(position *model.Position, realizedPL float64) *PositionClosedEvent {
	return &PositionClosedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PositionClosed", position.Symbol, 1),
		Position:        position,
		RealizedPL:      realizedPL,
	}
}

// BalanceUpdatedEvent は残高更新イベント
type BalanceUpdatedEvent struct {
	BaseDomainEvent
	PreviousCash   float64 `json:"previous_cash"`
	NewCash        float64 `json:"new_cash"`
	PreviousPower  float64 `json:"previous_buying_power"`
	NewBuyingPower float64 `json:"new_buying_power"`
}

// NewBalanceUpdatedEvent は新しい残高更新イベントを作成
func NewBalanceUpdatedEvent(prevCash, newCash, prevPower, newPower float64) *BalanceUpdatedEvent {
	return &BalanceUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("BalanceUpdated", "balance", 1),
		PreviousCash:    prevCash,
		NewCash:         newCash,
		PreviousPower:   prevPower,
		NewBuyingPower:  newPower,
	}
}

// PriceUpdatedEvent は価格更新イベント
type PriceUpdatedEvent struct {
	BaseDomainEvent
	Symbol        string  `json:"symbol"`
	PreviousPrice float64 `json:"previous_price"`
	NewPrice      float64 `json:"new_price"`
	Volume        int64   `json:"volume"`
}

// NewPriceUpdatedEvent は新しい価格更新イベントを作成
func NewPriceUpdatedEvent(symbol string, prevPrice, newPrice float64, volume int64) *PriceUpdatedEvent {
	return &PriceUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PriceUpdated", symbol, 1),
		Symbol:          symbol,
		PreviousPrice:   prevPrice,
		NewPrice:        newPrice,
		Volume:          volume,
	}
}

// RiskLimitExceededEvent はリスク制限超過イベント
type RiskLimitExceededEvent struct {
	BaseDomainEvent
	RiskType     string  `json:"risk_type"`
	CurrentValue float64 `json:"current_value"`
	Limit        float64 `json:"limit"`
	Symbol       string  `json:"symbol,omitempty"`
}

// NewRiskLimitExceededEvent は新しいリスク制限超過イベントを作成
func NewRiskLimitExceededEvent(riskType string, currentValue, limit float64, symbol string) *RiskLimitExceededEvent {
	aggregateID := "risk"
	if symbol != "" {
		aggregateID = symbol
	}

	return &RiskLimitExceededEvent{
		BaseDomainEvent: NewBaseDomainEvent("RiskLimitExceeded", aggregateID, 1),
		RiskType:        riskType,
		CurrentValue:    currentValue,
		Limit:           limit,
		Symbol:          symbol,
	}
}
