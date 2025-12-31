package model

import "gorm.io/gorm"

type PositionType string
type PositionAccountType string // New type for account type

const (
	PositionTypeLong  PositionType = "LONG"
	PositionTypeShort PositionType = "SHORT"

	PositionAccountTypeCash           PositionAccountType = "CASH"   // 現物
	PositionAccountTypeMarginNew      PositionAccountType = "MARGIN_NEW" // 信用新規
	PositionAccountTypeMarginRepay    PositionAccountType = "MARGIN_REPAY" // 信用返済
)

type Position struct {
	gorm.Model
	Symbol string `gorm:"index"` // 銘柄コード
	// AccountID        uint
	PositionType        PositionType        `gorm:"index"`
	PositionAccountType PositionAccountType // New field for cash/margin distinction
	AveragePrice        float64
	Quantity            int
	// Account      Account `gorm:"foreignKey:AccountID;references:ID"`

	// Fields for trailing stop logic
	HighestPrice      float64
	TrailingStopPrice float64 `gorm:"-"`
}
