package model

import "time"

// BalanceSummary は口座サマリーのドメインモデルです。
type BalanceSummary struct {
	// 総資産 (円)
	TotalAssets float64
	// 現物買付可能額 (円)
	CashBuyingPower float64
	// 信用新規建可能額 (円)
	MarginBuyingPower float64
	// 出金可能額 (円)
	WithdrawalPossibleAmount float64
	// 委託保証金率 (%)
	MarginRate float64
	// 最終更新日時
	UpdatedAt time.Time
}
