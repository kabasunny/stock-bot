// request/zan_kai_kanougaku_suii.go
package request

// ReqZanKaiKanougakuSuii は可能額推移のリクエストを表すDTO
type ReqZanKaiKanougakuSuii struct {
	P_no      string `json:"p_no"`      // p_no
	P_sd_date string `json:"p_sd_date"` // システム日付
	SJsonOfmt string `json:"sJsonOfmt"` // JSON出力フォーマット
	SCLMID    string `json:"sCLMID"`    // 機能ID, CLMZanKaiKanougakuSuii
}
