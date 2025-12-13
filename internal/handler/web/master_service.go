package web

import (
	"context"
	"errors"
	"log/slog"
	"stock-bot/gen/master"   // Goa-generated service interface and types
	"stock-bot/internal/app" // Our app-level use case interface and types

	goa "goa.design/goa/v3/pkg"
)

// MasterService implements the master.Service interface.
type MasterService struct {
	masterUseCase app.MasterUseCase
	logger        *slog.Logger
}

// NewMasterService creates a new MasterService.
func NewMasterService(masterUseCase app.MasterUseCase, logger *slog.Logger) *MasterService {
	return &MasterService{
		masterUseCase: masterUseCase,
		logger:        logger,
	}
}

// GetStock implements the get_stock action.
func (s *MasterService) GetStock(ctx context.Context, p *master.GetStockPayload) (res *master.StockbotStockMaster, err error) {
	s.logger.Info("GetStock called", slog.String("symbol", p.Symbol))

	// Call the use case
	appStockMaster, err := s.masterUseCase.GetStock(ctx, p.Symbol)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			s.logger.Warn("Stock master not found", slog.String("symbol", p.Symbol))
			// Returning a specific Goa error type for not_found
			return nil, &goa.ServiceError{Name: "not_found", ID: "not_found", Message: "Stock master not found", Fault: true}
		}
		s.logger.Error("Failed to get stock master from use case", slog.Any("error", err))
		return nil, err
	}

	// Map app.StockMasterResult to master.StockbotStockMaster
	res = &master.StockbotStockMaster{
		Symbol:       appStockMaster.Symbol,
		Name:         appStockMaster.Name,
		NameKana:     &appStockMaster.NameKana,
		Market:       appStockMaster.Market,
		IndustryCode: &appStockMaster.IndustryCode,
		IndustryName: &appStockMaster.IndustryName,
	}

	s.logger.Info("GetStock successful", slog.String("symbol", res.Symbol))
	return res, nil
}
