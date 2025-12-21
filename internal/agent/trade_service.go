package agent

import (
	"context"
	"stock-bot/domain/model"
)

// TradeService はエージェントがトレードサービス（Go APIラッパー）と連携するためのインターフェース
type TradeService interface {
	// GetPositions は現在の保有ポジションを取得する
	GetPositions(ctx context.Context) ([]*model.Position, error)
	// GetOrders は発注中の注文を取得する
	GetOrders(ctx context.Context) ([]*model.Order, error)
	// GetBalance は口座残高を取得する
	GetBalance(ctx context.Context) (*Balance, error) // agent.Balance型を使用
	// PlaceOrder は注文を発行する
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error)
	// CancelOrder は注文をキャンセルする
	CancelOrder(ctx context.Context, orderID string) error
	// TODO: 他に必要なAPIを随時追加
}

// PlaceOrderRequest は注文発行に必要な情報を保持する
type PlaceOrderRequest struct {
	Symbol    string
	TradeType model.TradeType
	OrderType model.OrderType
	Quantity  int
	Price     float64 // 指値の場合のみ
}
