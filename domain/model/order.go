package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderID      string      `gorm:"uniqueIndex"` // 証券会社固有の注文ID
	Symbol       string      `gorm:"index"`       // 銘柄コード
	TradeType    TradeType   `gorm:"index"`       // 買い/売り
	OrderType    OrderType   `gorm:"index"`       // 成行/指値など
	Quantity     int         `gorm:"not null"`
	Price        float64     // 指値の場合
	TriggerPrice float64     // 逆指値の場合
	TimeInForce  TimeInForce `gorm:"index;default:'DAY'"`                   // 有効期限
	OrderStatus  OrderStatus `gorm:"index"`                                 // 注文状態
	IsMargin     bool        `gorm:"not null;default:false"`                // 信用取引かどうか
	Executions   []Execution `gorm:"foreignKey:OrderID;references:OrderID"` // 約定情報
	// Account    Account `gorm:"foreignKey:AccountID;references:ID"`
}

type TradeType string

const (
	TradeTypeBuy  TradeType = "BUY"
	TradeTypeSell TradeType = "SELL"
)

type OrderType string

const (
	OrderTypeMarket    OrderType = "MARKET"
	OrderTypeLimit     OrderType = "LIMIT"
	OrderTypeStop      OrderType = "STOP"
	OrderTypeStopLimit OrderType = "STOP_LIMIT"
)

type TimeInForce string

const (
	TimeInForceDay TimeInForce = "DAY" // 当日限り
	// 他の有効期間は、必要になったら追加
	// TimeInForceGTC TimeInForce = "GTC" // Good 'Til Canceled
	// TimeInForceIOC TimeInForce = "IOC" // Immediate Or Cancel
	// ...
)

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"              // 新規
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED" // 一部約定
	OrderStatusFilled          OrderStatus = "FILLED"           // 完全約定
	OrderStatusCanceled        OrderStatus = "CANCELED"         // 取消済
	OrderStatusRejected        OrderStatus = "REJECTED"         // 拒否
	OrderStatusExpired         OrderStatus = "EXPIRED"          // 期限切れ
)

// IsUnexecuted は注文が市場でまだ有効（未約定または一部約定）かどうかを返す
func (os OrderStatus) IsUnexecuted() bool {
	return os == OrderStatusNew || os == OrderStatusPartiallyFilled
}

type Execution struct {
	gorm.Model
	OrderID           string `gorm:"index"`       // 注文ID (Orderモデルのgorm.Model.IDを参照)
	ExecutionID       string `gorm:"uniqueIndex"` // 証券会社固有の約定ID
	ExecutionTime     time.Time
	ExecutionPrice    float64
	ExecutionQuantity int
	Commission        float64 // 手数料
}

// IsUnexecuted は注文が市場でまだ有効かどうかを返す
func (o *Order) IsUnexecuted() bool {
	return o.OrderStatus.IsUnexecuted()
}
