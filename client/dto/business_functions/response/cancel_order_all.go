// business_functions/res_cancel_order_all.go
package business_functions

// ResCancelOrderAll は株式一括取消のレスポンスを表すDTO
type ResCancelOrderAll struct {
	P_no        string `json:"p_no"`        // p_no
	SCLMID      string `json:"sCLMID"`      // 機能ID, CLMKabuCancelOrderAll
	SResultCode string `json:"sResultCode"` // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText string `json:"sResultText"` // 結果テキスト, CLMKabuNewOrder.sResultText 参照
}
