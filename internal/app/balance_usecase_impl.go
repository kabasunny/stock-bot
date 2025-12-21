package app

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client"
	"strconv"
)

// balanceUseCaseImpl implements the BalanceUseCase interface.
type balanceUseCaseImpl struct {
	balanceClient client.BalanceClient
}

// NewBalanceUseCaseImpl creates a new BalanceUseCase.
func NewBalanceUseCaseImpl(balanceClient client.BalanceClient) BalanceUseCase {
	return &balanceUseCaseImpl{balanceClient: balanceClient}
}

// GetBalance retrieves the account balance summary, parses it, and returns it.
func (uc *balanceUseCaseImpl) GetBalance(ctx context.Context, session *client.Session) (*BalanceResult, error) {
	// Call the client to get the raw summary data
	summary, err := uc.balanceClient.GetZanKaiSummary(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance summary from client: %w", err)
	}

	// The API returns "0" for success. Any other code is an error.
	if summary.ResultCode != "0" {
		return nil, fmt.Errorf("client returned error: result_code=%s, text=%s", summary.ResultCode, summary.ResultText)
	}

	// Parse string fields into proper types.
	availableCashForStock, err := strconv.ParseFloat(summary.GenbutuKabuKaituke, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AvailableCashForStock '%s': %w", summary.GenbutuKabuKaituke, err)
	}

	availableMarginForNewPosition, err := strconv.ParseFloat(summary.SinyouSinkidate, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AvailableMarginForNewPosition '%s': %w", summary.SinyouSinkidate, err)
	}

	marginMaintenanceRate, err := strconv.ParseFloat(summary.HosyouKinritu, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MarginMaintenanceRate '%s': %w", summary.HosyouKinritu, err)
	}

	withdrawableCash, err := strconv.ParseFloat(summary.Syukkin, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WithdrawableCash '%s': %w", summary.Syukkin, err)
	}

	hasMarginCall := summary.OisyouHasseiFlg == "1"

	// Create and return the result struct
	result := &BalanceResult{
		AvailableCashForStock:         availableCashForStock,
		AvailableMarginForNewPosition: availableMarginForNewPosition,
		MarginMaintenanceRate:         marginMaintenanceRate,
		WithdrawableCash:              withdrawableCash,
		HasMarginCall:                 hasMarginCall,
	}

	return result, nil
}
