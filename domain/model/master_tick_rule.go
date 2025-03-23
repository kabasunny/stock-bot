// domain/model/tick_rule.go
package model

import "gorm.io/gorm"

// TickRule は、呼値の情報を表すモデル
// 特定の銘柄に対する呼値のルール CLMYobine の情報に対応
type TickRule struct {
	gorm.Model
	IssueCode      string      `gorm:"index;size:255"` // 銘柄コード (FK)
	TickUnitNumber string      // 呼値の単位番号
	ApplicableDate string      // 適用日 (YYYYMMDD)
	TickLevels     []TickLevel `gorm:"foreignKey:TickRuleID;references:ID"` // 呼値の段階 (1対多の関係)
}

// TickLevel は、呼値の各段階を表すモデル
// 呼値のルールにおける各段階 CLMYobine の情報に対応
type TickLevel struct {
	gorm.Model
	TickRuleID uint    // TickRule の ID (FK)
	LowerPrice float64 // 基準値段の下限
	UpperPrice float64 // 基準値段の上限
	TickValue  float64 // 呼値
}
