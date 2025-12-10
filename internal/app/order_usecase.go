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
	Symbol    string
	TradeType model.TradeType
	OrderType model.OrderType
	Quantity  uint64
	Price     float64
	IsMargin  bool
}
