package agent

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// TradeService はエージェントがトレードサービス（Go APIラッパー）と連携するためのインターフェース
type TradeService interface {
	// GetSession は現在のAPIセッション情報を取得する
	GetSession() *client.Session
	// GetPositions は現在の保有ポジションを取得する
	GetPositions(ctx context.Context) ([]*model.Position, error)
	// GetOrders は発注中の注文を取得する
	GetOrders(ctx context.Context) ([]*model.Order, error)
	// GetBalance は口座残高を取得する
	GetBalance(ctx context.Context) (*Balance, error) // agent.Balance型を使用
	// GetPrice は指定した銘柄の現在価格を取得する
	GetPrice(ctx context.Context, symbol string) (float64, error)
	// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
	GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error)
	// PlaceOrder は注文を発行する
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error)
	// CancelOrder は注文をキャンセルする
	CancelOrder(ctx context.Context, orderID string) error
	// TODO: 他に必要なAPIを随時追加
}

// HistoricalPrice は時系列データの一点を表現する
type HistoricalPrice struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

// PlaceOrderRequest は注文発行に必要な情報を保持する
type PlaceOrderRequest struct {
	Symbol       string
	TradeType    model.TradeType
	OrderType    model.OrderType
	Quantity     int
	Price        float64 // 指値の場合のみ
	TriggerPrice float64 // 逆指値の場合のみ
}
