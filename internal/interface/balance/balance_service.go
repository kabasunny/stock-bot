package balance

import (
	"context"
	"time"

	"stock-bot/internal/app"
	genbalance "stock-bot/gen/balance"
)

// balanceService は balance.Service インターフェースを実装します。
type balanceService struct {
	usecase app.BalanceUseCase
}

// NewBalanceService は新しい balance サービスを作成します。
func NewBalanceService(usecase app.BalanceUseCase) genbalance.Service {
	return &balanceService{usecase: usecase}
}

// Summary はユースケースを呼び出し、ドメインモデルをGoaのレスポンス型に変換します。
func (s *balanceService) Summary(ctx context.Context) (*genbalance.StockBalanceSummary, error) {
	// ユースケースを呼び出してドメインモデルを取得
	domainSummary, err := s.usecase.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	// ドメインモデルをGoaのレスポンス型に変換
	res := &genbalance.StockBalanceSummary{
		TotalAssets:              domainSummary.TotalAssets,
		CashBuyingPower:          domainSummary.CashBuyingPower,
		MarginBuyingPower:        domainSummary.MarginBuyingPower,
		WithdrawalPossibleAmount: domainSummary.WithdrawalPossibleAmount,
		MarginRate:               domainSummary.MarginRate,
		UpdatedAt:                domainSummary.UpdatedAt.Format(time.RFC3339), // time.Time を string に変換
	}

	return res, nil
}
