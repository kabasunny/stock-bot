// business_functions/req_cancel_order_all.go
package business_functions

// ReqCancelOrderAll は株式一括取消のリクエストを表すDTO
type ReqCancelOrderAll struct {
	P_no            string `json:"p_no"`            // p_no
	P_sd_date       string `json:"p_sd_date"`       // システム日付
	SJsonOfmt       string `json:"sJsonOfmt"`       // JSON出力フォーマット
	SCLMID          string `json:"sCLMID"`          // 機能ID, CLMKabuCancelOrderAll
	SSecondPassword string `json:"sSecondPassword"` // 第二パスワード
}
