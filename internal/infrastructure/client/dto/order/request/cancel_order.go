// request/cancel_order.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqCancelOrder は株式取消注文のリクエストを表すDTO
type ReqCancelOrder struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`          // 機能ID, CLMKabuCancelOrder
	OrderNumber     string `json:"sOrderNumber"`    // 注文番号, CLMKabuNewOrder.sOrderNumber
	EigyouDay       string `json:"sEigyouDay"`      // 営業日, CLMKabuNewOrder.sEigyouDay
	SecondPassword  string `json:"sSecondPassword"` // 第二パスワード
}
