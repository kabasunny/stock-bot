package order

import (
	"context"
	"log/slog"
	"stock-bot/domain/model"
	ordersvr "stock-bot/gen/order"
	"stock-bot/internal/app"
	"time"
)

// orderService は order.Service インターフェースを実装します。
type orderService struct {
	usecase app.OrderUseCase
	logger  *slog.Logger
}

// NewOrderService は新しい order サービスを作成します。
func NewOrderService(usecase app.OrderUseCase) ordersvr.Service {
	return &orderService{
		usecase: usecase,
		logger:  slog.Default(),
	}
}

// NewOrder は新しい株式注文を作成します。
func (s *orderService) NewOrder(ctx context.Context, p *ordersvr.NewOrderPayload) (res *ordersvr.StockOrder, err error) {
	s.logger.InfoContext(ctx, "order.NewOrder", "payload", p)

	// 1. ペイロードをUsecaseのパラメータに変換
	params := app.OrderParams{
		Symbol:    p.Symbol,
		TradeType: model.TradeType(p.TradeType),
		OrderType: model.OrderType(p.OrderType),
		Quantity:  p.Quantity,
		IsMargin:  p.IsMargin,
	}
	if p.Price != nil {
		params.Price = *p.Price
	}
	if p.TriggerPrice != nil {
		params.TriggerPrice = *p.TriggerPrice
	}
	if p.TimeInForce != "" {
		params.TimeInForce = model.TimeInForce(p.TimeInForce)
	}

	// 2. UsecaseのExecuteOrderを呼び出し
	order, err := s.usecase.ExecuteOrder(ctx, params)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to execute order", "error", err)
		return nil, err
	}

	// 3. Usecaseの結果をGoaのレスポンス型に変換
	res = &ordersvr.StockOrder{
		OrderID:     order.OrderID,
		Symbol:      order.Symbol,
		TradeType:   string(order.TradeType),
		OrderType:   string(order.OrderType),
		Quantity:    order.Quantity,
		OrderStatus: string(order.OrderStatus),
		IsMargin:    order.IsMargin,
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}
	if order.Price != 0 {
		res.Price = &order.Price
	}
	if order.TriggerPrice != 0 {
		res.TriggerPrice = &order.TriggerPrice
	}
	if order.TimeInForce != "" {
		timeInForceStr := string(order.TimeInForce)
		res.TimeInForce = &timeInForceStr
	}

	return res, nil
}
