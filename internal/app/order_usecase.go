// internal/app/order_usecase.go
package app

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
)

type OrderUseCase interface {
	ExecuteOrder(ctx context.Context, session *client.Session, orderParams OrderParams) (*model.Order, error)
}

type OrderParams struct {
	Symbol              string
	TradeType           model.TradeType
	OrderType           model.OrderType
	Quantity            uint64
	Price               float64
	PositionAccountType model.PositionAccountType
}