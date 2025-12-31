package model

import "time"

// Execution は約定情報を表すドメインモデルです。
type Execution struct {
	ExecutionID string    `json:"execution_id" gorm:"primaryKey"` // 約定ID
	OrderID     string    `json:"order_id" gorm:"index"`          // 関連する注文ID
	Symbol      string    `json:"symbol"`                         // 銘柄コード
	TradeType   TradeType `json:"trade_type"`                     // 売買区分 (BUY/SELL)
	Quantity    int       `json:"quantity"`                       // 約定数量
	Price       float64   `json:"price"`                          // 約定価格
	ExecutedAt  time.Time `json:"executed_at"`                    // 約定日時
	Commission  float64   `json:"commission,omitempty"`           // 手数料 (オプション)
}
