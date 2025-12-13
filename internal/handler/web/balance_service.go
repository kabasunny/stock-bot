package web

import (
	"context"
	"log/slog"
	"stock-bot/gen/balance" // This is the Goa-generated service interface and types
	"stock-bot/internal/app" // This is our app-level use case interface and types
)

// BalanceService implements the balance.Service interface.
type BalanceService struct {
	balanceUseCase app.BalanceUseCase
	logger         *slog.Logger
}

// NewBalanceService creates a new BalanceService.
func NewBalanceService(balanceUseCase app.BalanceUseCase, logger *slog.Logger) *BalanceService {
	return &BalanceService{
		balanceUseCase: balanceUseCase,
		logger:         logger,
	}
}

// Get implements the get action.
func (s *BalanceService) Get(ctx context.Context) (res *balance.StockbotBalance, err error) {
	s.logger.Info("GetBalance called")

	// Call the use case
	balanceResult, err := s.balanceUseCase.GetBalance(ctx)
	if err != nil {
		s.logger.Error("Failed to get balance from use case", slog.Any("error", err))
		return nil, err
	}

	// Map app.BalanceResult to balance.StockbotBalance (Goa generated type)
	res = &balance.StockbotBalance{
		AvailableCashForStock:         balanceResult.AvailableCashForStock,
		AvailableMarginForNewPosition: balanceResult.AvailableMarginForNewPosition,
		MarginMaintenanceRate:         balanceResult.MarginMaintenanceRate,
		WithdrawableCash:              balanceResult.WithdrawableCash,
		HasMarginCall:                 balanceResult.HasMarginCall,
	}

	s.logger.Info("GetBalance successful", slog.Any("result", res))
	return res, nil
}
