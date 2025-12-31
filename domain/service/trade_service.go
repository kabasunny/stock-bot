package service

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// TradeService はトレード関連のドメインサービスインターフェース
type TradeService interface {
	// GetSession は現在のAPIセッション情報を取得する
	GetSession() *client.Session
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
}

// Balance は残高情報を表現する
type Balance struct {
	Cash        float64 // 現金残高
	BuyingPower float64 // 買付余力
}

// HistoricalPrice は時系列データの一点を表現する
type HistoricalPrice struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

// PlaceOrderRequest は注文発行に必要な情報を保持する
type PlaceOrderRequest struct {
	Symbol              string
	TradeType           model.TradeType
	OrderType           model.OrderType
	Quantity            int
	Price               float64                   // 指値の場合のみ
	TriggerPrice        float64                   // 逆指値の場合のみ
	PositionAccountType model.PositionAccountType // ポジションの口座タイプ（現物/信用）
}
