// request/order_list_detail.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqOrderListDetail は注文約定一覧（詳細）のリクエストを表すDTO
type ReqOrderListDetail struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`       // 機能ID, CLMOrderListDetail
	OrderNumber     string `json:"sOrderNumber"` // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	EigyouDay       string `json:"sEigyouDay"`   // 営業日, CLMKabuNewOrder.sEigyouDay 参照
}
