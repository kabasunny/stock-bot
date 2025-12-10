package web

import (
	"context"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/internal/app"
	ordersvr "stock-bot/gen/order"
)

// order.Serviceインターフェースを実装する構造体
type OrderService struct {
	usecase app.OrderUseCase
	logger  *slog.Logger
}

// コンストラクタ
func NewOrderService(usecase app.OrderUseCase, logger *slog.Logger) ordersvr.Service {
	return &OrderService{
		usecase: usecase,
		logger:  logger,
	}
}

// Createメソッドの実装 (Goaが生成したインターフェースを満たす)
func (s *OrderService) Create(ctx context.Context, p *ordersvr.CreatePayload) (res *ordersvr.CreateResult, err error) {
	s.logger.Info("order.create method called", slog.Any("payload", p))

	// PayloadをUsecaseが要求するパラメータに変換
	// TODO: 文字列からの変換部分のバリデーションを強化する
	orderParams := app.OrderParams{
		Symbol:    p.Symbol,
		TradeType: model.TradeType(p.TradeType),
		OrderType: model.OrderType(p.OrderType),
		Quantity:  p.Quantity,
		Price:     p.Price,
		IsMargin:  p.IsMargin,
	}

	// Usecaseの呼び出し
	createdOrder, err := s.usecase.ExecuteOrder(ctx, orderParams)
	if err != nil {
		s.logger.Error("failed to execute order", slog.Any("error", err))
		return nil, err // Goaが適切なエラーレスポンスに変換してくれる
	}

	// 結果をレスポンスに設定
	res = &ordersvr.CreateResult{OrderID: createdOrder.OrderID}
	s.logger.Info("order.create method successfully processed", slog.String("order_id", createdOrder.OrderID))

	return res, nil
}