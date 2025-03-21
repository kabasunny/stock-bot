// business_functions/req_zan_kai_kanougaku_suii.go
package business_functions

// ReqZanKaiKanougakuSuii は可能額推移のリクエストを表すDTO
type ReqZanKaiKanougakuSuii struct {
	P_no      string `json:"p_no"`      // p_no
	P_sd_date string `json:"p_sd_date"` // システム日付
	SJsonOfmt string `json:"sJsonOfmt"` // JSON出力フォーマット
	SCLMID    string `json:"sCLMID"`    // 機能ID, CLMZanKaiKanougakuSuii
}
