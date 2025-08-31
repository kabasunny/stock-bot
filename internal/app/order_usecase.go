// internal/app/order_usecase.go
package app

import (
	"context"
	"stock-bot/domain/model"
)

type OrderUseCase interface {
	ExecuteOrder(ctx context.Context, orderParams OrderParams) (*model.Order, error)
}

type OrderParams struct {
	Symbol       string
	TradeType    model.TradeType
	OrderType    model.OrderType
	Quantity     int
	Price        float64
	TriggerPrice float64
	TimeInForce  model.TimeInForce
	IsMargin     bool
}
