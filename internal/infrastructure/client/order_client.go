package client

import (
	"context"

	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
)

type NewOrderParams struct {
	ZyoutoekiKazeiC          string
	IssueCode                string
	SizyouC                  string
	BaibaiKubun              string
	Condition                string
	OrderPrice               string
	OrderSuryou              string
	GenkinShinyouKubun       string
	OrderExpireDay           string
	GyakusasiOrderType       string
	GyakusasiZyouken         string
	GyakusasiPrice           string
	TatebiType               string
	TategyokuZyoutoekiKazeiC string
	CLMKabuHensaiData        []request.ReqHensaiData
}

// OrderClient は、注文関連の API を扱うインターフェース
type OrderClient interface {
	// NewOrder は、新規の株式注文を行う
	NewOrder(ctx context.Context, params NewOrderParams) (*response.ResNewOrder, error)
	// CorrectOrder は、既存の株式注文を訂正する
	CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error)
	// CancelOrder は、既存の株式注文を取り消す
	CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error)
	// CancelOrderAll は、顧客の全ての未約定注文を一括で取り消す
	CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error)
	// GetOrderList は、注文の一覧を取得す
	GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error)
	// GetOrderListDetail は、指定した注文の約定情報（詳細）を取得す
	GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error)
}
