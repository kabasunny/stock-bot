// internal/app/account_usecase.go
package app

import (
	"context"
	"stock-bot/domain/model"
)

type BalanceUseCase interface {
	// GetSummary は、口座のサマリー情報を取得する
	GetSummary(ctx context.Context) (*model.BalanceSummary, error)
	// CanEntry は、指定された銘柄にエントリー可能かどうかと、口座情報を返す
	CanEntry(ctx context.Context, issueCode string) (bool, float64, error)
}
