// response/cancel_order.go
package response

// ResCancelOrder は株式取消注文のレスポンスを表すDTO
type ResCancelOrder struct {
	P_no                   string `json:"p_no"`                   // p_no
	SCLMID                 string `json:"sCLMID"`                 // 機能ID, CLMKabuCancelOrder
	SResultCode            string `json:"sResultCode"`            // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText            string `json:"sResultText"`            // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SOrderNumber           string `json:"sOrderNumber"`           // 注文番号, 要求設定値
	SEigyouDay             string `json:"sEigyouDay"`             // 営業日, 要求設定値
	SOrderUkewatasiKingaku string `json:"sOrderUkewatasiKingaku"` // 注文受渡金額
	SOrderDate             string `json:"sOrderDate"`             // 注文日時, YYYYMMDDHHMMSS
}
