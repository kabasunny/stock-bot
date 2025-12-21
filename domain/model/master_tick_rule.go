// domain/model/tick_rule.go
package model

import "time"

// TickRule は、呼値の情報を表すモデル
type TickRule struct {
	TickUnitNumber string      `gorm:"primaryKey;size:255"` // 呼値の単位番号 (主キー)
	ApplicableDate string      // 適用日 (YYYYMMDD)
	TickLevels     []TickLevel `gorm:"foreignKey:TickRuleUnitNumber;references:TickUnitNumber"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TickLevel は、呼値の各段階を表すモデル
type TickLevel struct {
	ID                 uint    `gorm:"primaryKey"`
	TickRuleUnitNumber string  `gorm:"index;size:255"` // 外部キー
	LowerPrice         float64 // 基準値段の下限
	UpperPrice         float64 // 基準値段の上限
	TickValue          float64 // 呼値
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
