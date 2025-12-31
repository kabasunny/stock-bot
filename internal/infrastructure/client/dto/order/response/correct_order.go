// response/correct_order.go
package response

// ResCorrectOrder は株式訂正注文のレスポンスを表すDTO
type ResCorrectOrder struct {
	P_no                  string `json:"p_no"`                   // p_no
	CLMID                 string `json:"sCLMID"`                 // 機能ID, CLMKabuCorrectOrder
	ResultCode            string `json:"sResultCode"`            // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText            string `json:"sResultText"`            // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	OrderNumber           string `json:"sOrderNumber"`           // 注文番号, 要求設定値
	EigyouDay             string `json:"sEigyouDay"`             // 営業日, 要求設定値
	OrderUkewatasiKingaku string `json:"sOrderUkewatasiKingaku"` // 注文受渡金額
	OrderTesuryou         string `json:"sOrderTesuryou"`         // 注文手数料
	OrderSyouhizei        string `json:"sOrderSyouhizei"`        // 注文消費税
	OrderDate             string `json:"sOrderDate"`             // 注文日時, YYYYMMDDHHMMSS
}
