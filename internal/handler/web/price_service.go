package web

import (
	"context"
	"log/slog"
	"stock-bot/gen/price"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client" // For Session
)

// PriceService implements the price.Service interface.
type PriceService struct {
	priceUseCase app.PriceUseCase
	logger       *slog.Logger
	session      *client.Session
}

// NewPriceService creates a new PriceService.
func NewPriceService(priceUseCase app.PriceUseCase, logger *slog.Logger, session *client.Session) *PriceService {
	return &PriceService{
		priceUseCase: priceUseCase,
		logger:       logger,
		session:      session,
	}
}

// Get implements price.Service.
func (s *PriceService) Get(ctx context.Context, p *price.GetPayload) (res *price.StockbotPrice, err error) {
	s.logger.Info("PriceService.Get", slog.String("symbol", p.Symbol))

	// ユースケースを呼び出す
	stockPrice, err := s.priceUseCase.Get(ctx, p.Symbol)
	if err != nil {
		s.logger.Error("failed to get price via usecase", slog.Any("error", err))
		return nil, err
	}

	return stockPrice, nil
}
