package app

import (
	"context"
)

// BalanceResult is the result type for the GetBalance use case.
// It contains curated balance information.
type BalanceResult struct {
	AvailableCashForStock         float64
	AvailableMarginForNewPosition float64
	MarginMaintenanceRate         float64
	WithdrawableCash              float64
	HasMarginCall                 bool
}

// BalanceUseCase defines the interface for balance-related use cases.
type BalanceUseCase interface {
	GetBalance(ctx context.Context) (*BalanceResult, error)
}