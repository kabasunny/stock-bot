// request/cancel_order.go
package request

// ReqCancelOrder は株式取消注文のリクエストを表すDTO
type ReqCancelOrder struct {
	P_no            string `json:"p_no"`            // p_no
	P_sd_date       string `json:"p_sd_date"`       // システム日付
	SJsonOfmt       string `json:"sJsonOfmt"`       // JSON出力フォーマット
	SCLMID          string `json:"sCLMID"`          // 機能ID, CLMKabuCancelOrder
	SOrderNumber    string `json:"sOrderNumber"`    // 注文番号, CLMKabuNewOrder.sOrderNumber
	SEigyouDay      string `json:"sEigyouDay"`      // 営業日, CLMKabuNewOrder.sEigyouDay
	SSecondPassword string `json:"sSecondPassword"` // 第二パスワード
}
