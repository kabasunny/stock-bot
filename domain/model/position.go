package model

import "gorm.io/gorm"

type PositionType string

const (
	PositionTypeLong  PositionType = "LONG"
	PositionTypeShort PositionType = "SHORT"
)

type Position struct {
	gorm.Model
	Symbol string `gorm:"index"` // 銘柄コード
	// AccountID        uint
	PositionType PositionType `gorm:"index"`
	AveragePrice float64
	Quantity     int
	// Account      Account `gorm:"foreignKey:AccountID;references:ID"`

	// Fields for trailing stop logic (in-memory only)
	HighestPrice      float64 `gorm:"-"`
	TrailingStopPrice float64 `gorm:"-"`
}
