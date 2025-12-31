package service

import (
	"context"
	"stock-bot/domain/model"
	"time"
)

// TradeService はトレード関連のドメインサービスインターフェース
type TradeService interface {
	// GetSession は現在のAPIセッション情報を取得する
	GetSession() *model.Session
	// GetPositions は現在の保有ポジションを取得する
	GetPositions(ctx context.Context) ([]*model.Position, error)
	// GetOrders は発注中の注文を取得する
	GetOrders(ctx context.Context) ([]*model.Order, error)
	// GetBalance は口座残高を取得する
	GetBalance(ctx context.Context) (*Balance, error)
	// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
	GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error)
	// PlaceOrder は注文を発行する
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error)
	// CancelOrder は注文をキャンセルする
	CancelOrder(ctx context.Context, orderID string) error
	// CorrectOrder は注文を訂正する
	CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error)
	// CancelAllOrders は全ての未約定注文をキャンセルする
	CancelAllOrders(ctx context.Context) (int, error)
	// GetOrderHistory は注文履歴を取得する
	GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error)
	// HealthCheck はサービスの健康状態をチェックする
	HealthCheck(ctx context.Context) (*HealthStatus, error)
}

// Balance は残高情報を表す
type Balance struct {
	Cash        float64 `json:"cash"`         // 現金残高
	BuyingPower float64 `json:"buying_power"` // 買付余力
}

// HistoricalPrice は過去の価格情報を表す
type HistoricalPrice struct {
	Date   time.Time `json:"date"`   // 日付
	Open   float64   `json:"open"`   // 始値
	High   float64   `json:"high"`   // 高値
	Low    float64   `json:"low"`    // 安値
	Close  float64   `json:"close"`  // 終値
	Volume int64     `json:"volume"` // 出来高
}

// PlaceOrderRequest は注文発行リクエストを表す
type PlaceOrderRequest struct {
	Symbol              string                    `json:"symbol"`                  // 銘柄コード
	TradeType           model.TradeType           `json:"trade_type"`              // 売買区分
	OrderType           model.OrderType           `json:"order_type"`              // 注文種別
	Quantity            int                       `json:"quantity"`                // 数量
	Price               float64                   `json:"price,omitempty"`         // 価格（指値の場合）
	TriggerPrice        *float64                  `json:"trigger_price,omitempty"` // トリガー価格（逆指値の場合）
	PositionAccountType model.PositionAccountType `json:"position_account_type"`   // ポジション口座区分
}

// HealthStatus はサービスの健康状態を表す
type HealthStatus struct {
	Status             string    `json:"status"`              // ステータス（healthy/unhealthy）
	Timestamp          time.Time `json:"timestamp"`           // チェック時刻
	SessionValid       bool      `json:"session_valid"`       // セッション有効性
	DatabaseConnected  bool      `json:"database_connected"`  // データベース接続状態
	WebSocketConnected bool      `json:"websocket_connected"` // WebSocket接続状態
}
