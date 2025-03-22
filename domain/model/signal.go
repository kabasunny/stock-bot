package model

import "time"

type SignalType string

const (
	SignalTypeBuy  SignalType = "BUY"
	SignalTypeSell SignalType = "SELL"
	// SignalTypeExit SignalType = "EXIT" // 手仕舞いシグナル (必要に応じて)
)

type Signal struct {
	ID          uint   `gorm:"primaryKey"`
	Symbol      string `gorm:"index"` // 銘柄コード
	SignalType  SignalType
	GeneratedAt time.Time
	Rationale   string // シグナルの根拠
	Price       float64
	// PriceRange  string // シグナルが有効な価格帯 (必要に応じて)
	// Expiration time.Time // 有効期限 (必要に応じて)
}
