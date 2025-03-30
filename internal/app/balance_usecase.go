// internal/app/account_usecase.go
package app

import "context"

type BalanceUseCase interface {
	// CanEntry は、指定された銘柄にエントリー可能かどうかと、口座情報を返す
	CanEntry(ctx context.Context, issueCode string) (bool, float64, error)
}
