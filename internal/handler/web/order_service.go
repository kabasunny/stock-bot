package web

import (
	"context"
	"log/slog"
	"stock-bot/domain/model"
	ordersvr "stock-bot/gen/order"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
)

// OrderService implements the order.Service interface.
type OrderService struct {
	usecase app.OrderUseCase
	logger  *slog.Logger
	session *client.Session
}

// NewOrderService creates a new order service.
func NewOrderService(usecase app.OrderUseCase, logger *slog.Logger, session *client.Session) ordersvr.Service {
	return &OrderService{
		usecase: usecase,
		logger:  logger,
		session: session,
	}
}

// Create implements create.
func (s *OrderService) Create(ctx context.Context, p *ordersvr.CreatePayload) (res *ordersvr.CreateResult, err error) {
	s.logger.Info("order.create method called", "payload", p)

	// PayloadからOrderParamsへの変換
	orderParams := app.OrderParams{
		Symbol:    p.Symbol,
		TradeType: model.TradeType(p.TradeType),
		OrderType: model.OrderType(p.OrderType),
		Quantity:  p.Quantity,
		Price:     p.Price,
		IsMargin:  p.IsMargin,
	}

	// UseCaseを呼び出す際にsessionを渡す
	createdOrder, err := s.usecase.ExecuteOrder(ctx, s.session, orderParams)
	if err != nil {
		s.logger.Error("Failed to execute order", "error", err)
		return nil, err // Goaが適切なエラーレスポンスに変換してくれる
	}

	res = &ordersvr.CreateResult{OrderID: createdOrder.OrderID}
	s.logger.Info("order.create method successfully processed.", "orderID", res.OrderID)

	return res, nil
}