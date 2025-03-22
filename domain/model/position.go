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
}
