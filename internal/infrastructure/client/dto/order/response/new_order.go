// response/new_order.go
package response

// ResNewOrder は株式新規注文のレスポンスを表すDTO
type ResNewOrder struct {
	P_no                  string `json:"p_no"`                   // p_no
	CLMID                 string `json:"sCLMID"`                 // 機能ID, CLMKabuNewOrder
	ResultCode            string `json:"sResultCode"`            // 結果コード, 0：正常
	ResultText            string `json:"sResultText"`            // 結果テキスト, 正常:""
	WarningCode           string `json:"sWarningCode"`           // 警告コード, 0：正常
	WarningText           string `json:"sWarningText"`           // 警告テキスト, 正常:""
	OrderNumber           string `json:"sOrderNumber"`           // 注文番号
	EigyouDay             string `json:"sEigyouDay"`             // 営業日, YYYYMMDD
	OrderUkewatasiKingaku string `json:"sOrderUkewatasiKingaku"` // 注文受渡金額
	OrderTesuryou         string `json:"sOrderTesuryou"`         // 注文手数料
	OrderSyouhizei        string `json:"sOrderSyouhizei"`        // 注文消費税
	Kinri                 string `json:"sKinri"`                 // 金利, -：現物取引場合
	OrderDate             string `json:"sOrderDate"`             // 注文日時, YYYYMMDDHHMMSS
}
