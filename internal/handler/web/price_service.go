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

// GetHistory implements price.Service.
func (s *PriceService) GetHistory(ctx context.Context, p *price.GetHistoryPayload) (res *price.StockbotHistoricalPrice, err error) {
	s.logger.Info("PriceService.GetHistory", slog.String("symbol", p.Symbol), slog.Uint64("days", uint64(p.Days)))

	// ユースケースを呼び出す
	historyResult, err := s.priceUseCase.GetHistory(ctx, p.Symbol, p.Days)
	if err != nil {
		s.logger.Error("failed to get historical price via usecase", slog.Any("error", err))
		return nil, err
	}

	// app.HistoricalPriceResult を price.StockbotHistoricalPrice に変換
	res = &price.StockbotHistoricalPrice{
		Symbol: historyResult.Symbol,
		History: make([]*price.HistoricalPriceItem, len(historyResult.History)),
	}

	for i, item := range historyResult.History {
		res.History[i] = &price.HistoricalPriceItem{
			Date:   item.Date,
			Open:   item.Open,
			High:   item.High,
			Low:    item.Low,
			Close:  item.Close,
			Volume: &item.Volume, // Pass as pointer
		}
	}

	return res, nil
}
