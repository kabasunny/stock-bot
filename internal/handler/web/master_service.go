package web

import (
	"context"
	"errors"
	"log/slog"
	"stock-bot/gen/master"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"

	goa "goa.design/goa/v3/pkg"
)

// MasterService implements the master.Service interface.
type MasterService struct {
	masterUseCase app.MasterUseCase
	logger        *slog.Logger
    session       *client.Session
}

// NewMasterService creates a new MasterService.
func NewMasterService(masterUseCase app.MasterUseCase, logger *slog.Logger, session *client.Session) *MasterService {
	return &MasterService{
		masterUseCase: masterUseCase,
		logger:        logger,
        session:       session,
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

// Update implements the update action.
func (s *MasterService) Update(ctx context.Context) (err error) {
	s.logger.Info("Update master data triggered.")

	if err := s.masterUseCase.DownloadAndStoreMasterData(ctx, s.session); err != nil {
		s.logger.Error("Failed to download and store master data", slog.Any("error", err))
		return err
	}

	s.logger.Info("Master data update completed successfully.")
	return nil
}
