// business_functions/req_correct_order.go
package business_functions

// ReqCorrectOrder は株式訂正注文のリクエストを表すDTO
type ReqCorrectOrder struct {
	P_no              string `json:"p_no"`              // p_no
	P_sd_date         string `json:"p_sd_date"`         // システム日付
	SJsonOfmt         string `json:"sJsonOfmt"`         // JSON出力フォーマット
	SCLMID            string `json:"sCLMID"`            // 機能ID, CLMKabuCorrectOrder
	SOrderNumber      string `json:"sOrderNumber"`      // 注文番号, CLMKabuNewOrder.sOrderNumber
	SEigyouDay        string `json:"sEigyouDay"`        // 営業日, CLMKabuNewOrder.sEigyouDay
	SCondition        string `json:"sCondition"`        // 執行条件, *：変更なし, 0：指定なし, 2：寄付, 4：引け, 6：不成
	SOrderPrice       string `json:"sOrderPrice"`       // 注文値段, *：変更なし, 0：成行に変更
	SOrderSuryou      string `json:"sOrderSuryou"`      // 注文数量, *：変更なし
	SOrderExpireDay   string `json:"sOrderExpireDay"`   // 注文期日, *：変更なし, 0：当日
	SGyakusasiZyouken string `json:"sGyakusasiZyouken"` // 逆指値条件, *：変更なし, 0：成行に変更
	SGyakusasiPrice   string `json:"sGyakusasiPrice"`   // 逆指値値段, *：変更なし, 0：成行に変更
	SSecondPassword   string `json:"sSecondPassword"`   // 第二パスワード
}
