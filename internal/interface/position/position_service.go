package position

import (
	"context"
	"log/slog"
	positionsvr "stock-bot/gen/position"
	"stock-bot/internal/app"
	"time"
)

// positionService は position.Service インターフェースを実装します。
type positionService struct {
	usecase app.PositionUseCase
	logger  *slog.Logger
}

// NewPositionService は新しい position サービスを作成します。
func NewPositionService(usecase app.PositionUseCase) positionsvr.Service {
	return &positionService{
		usecase: usecase,
		logger:  slog.Default(),
	}
}

// List は建玉の一覧を取得します。
func (s *positionService) List(ctx context.Context) ([]*positionsvr.StockPosition, error) {
	s.logger.InfoContext(ctx, "position.List")

	// 1. UsecaseのListを呼び出し
	positions, err := s.usecase.List(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list positions", "error", err)
		return nil, err
	}

	// 2. Usecaseの結果をGoaのレスポンス型に変換
	res := make([]*positionsvr.StockPosition, len(positions))
	for i, p := range positions {
		res[i] = &positionsvr.StockPosition{
			ID:           p.ID,
			Symbol:       p.Symbol,
			PositionType: string(p.PositionType),
			AveragePrice: p.AveragePrice,
			Quantity:     p.Quantity,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
		}
	}

	return res, nil
}
