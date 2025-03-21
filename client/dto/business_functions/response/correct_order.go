// business_functions/res_correct_order.go
package business_functions

// ResCorrectOrder は株式訂正注文のレスポンスを表すDTO
type ResCorrectOrder struct {
	P_no                   string `json:"p_no"`                   // p_no
	SCLMID                 string `json:"sCLMID"`                 // 機能ID, CLMKabuCorrectOrder
	SResultCode            string `json:"sResultCode"`            // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText            string `json:"sResultText"`            // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SOrderNumber           string `json:"sOrderNumber"`           // 注文番号, 要求設定値
	SEigyouDay             string `json:"sEigyouDay"`             // 営業日, 要求設定値
	SOrderUkewatasiKingaku string `json:"sOrderUkewatasiKingaku"` // 注文受渡金額
	SOrderTesuryou         string `json:"sOrderTesuryou"`         // 注文手数料
	SOrderSyouhizei        string `json:"sOrderSyouhizei"`        // 注文消費税
	SOrderDate             string `json:"sOrderDate"`             // 注文日時, YYYYMMDDHHMMSS
}
