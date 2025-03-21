// business_functions/res_new_order.go
package business_functions

// ResNewOrder は株式新規注文のレスポンスを表すDTO
type ResNewOrder struct {
	P_no                   string `json:"p_no"`                   // p_no
	SCLMID                 string `json:"sCLMID"`                 // 機能ID, CLMKabuNewOrder
	SResultCode            string `json:"sResultCode"`            // 結果コード, 0：正常
	SResultText            string `json:"sResultText"`            // 結果テキスト, 正常:""
	SWarningCode           string `json:"sWarningCode"`           // 警告コード, 0：正常
	SWarningText           string `json:"sWarningText"`           // 警告テキスト, 正常:""
	SOrderNumber           string `json:"sOrderNumber"`           // 注文番号
	SEigyouDay             string `json:"sEigyouDay"`             // 営業日, YYYYMMDD
	SOrderUkewatasiKingaku string `json:"sOrderUkewatasiKingaku"` // 注文受渡金額
	SOrderTesuryou         string `json:"sOrderTesuryou"`         // 注文手数料
	SOrderSyouhizei        string `json:"sOrderSyouhizei"`        // 注文消費税
	SKinri                 string `json:"sKinri"`                 // 金利, -：現物取引場合
	SOrderDate             string `json:"sOrderDate"`             // 注文日時, YYYYMMDDHHMMSS
}
