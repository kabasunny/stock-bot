// response/cancel_order_all.go
package response

// ResCancelOrderAll は株式一括取消のレスポンスを表すDTO
type ResCancelOrderAll struct {
	P_no       string `json:"p_no"`        // p_no
	CLMID      string `json:"sCLMID"`      // 機能ID, CLMKabuCancelOrderAll
	ResultCode string `json:"sResultCode"` // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText string `json:"sResultText"` // 結果テキスト, CLMKabuNewOrder.sResultText 参照
}
