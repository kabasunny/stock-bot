// request/correct_order.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqCorrectOrder は株式訂正注文のリクエストを表すDTO
type ReqCorrectOrder struct {
	dto.RequestBase         // 共通フィールドを埋め込む
	CLMID            string `json:"sCLMID"`            // 機能ID, CLMKabuCorrectOrder
	OrderNumber      string `json:"sOrderNumber"`      // 注文番号, CLMKabuNewOrder.sOrderNumber
	EigyouDay        string `json:"sEigyouDay"`        // 営業日, CLMKabuNewOrder.sEigyouDay
	Condition        string `json:"sCondition"`        // 執行条件, *：変更なし, 0：指定なし, 2：寄付, 4：引け, 6：不成
	OrderPrice       string `json:"sOrderPrice"`       // 注文値段, *：変更なし, 0：成行に変更
	OrderSuryou      string `json:"sOrderSuryou"`      // 注文数量, *：変更なし
	OrderExpireDay   string `json:"sOrderExpireDay"`   // 注文期日, *：変更なし, 0：当日
	GyakusasiZyouken string `json:"sGyakusasiZyouken"` // 逆指値条件, *：変更なし, 0：成行に変更
	GyakusasiPrice   string `json:"sGyakusasiPrice"`   // 逆指値値段, *：変更なし, 0：成行に変更
	SecondPassword   string `json:"sSecondPassword"`   // 第二パスワード
}
