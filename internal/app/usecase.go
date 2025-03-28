// internal/app/usecase.go
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
	// 他の注文に必要なパラメータ
}
