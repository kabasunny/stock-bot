// request/zan_kai_summary.go
package request

// ReqZanKaiSummary は可能額サマリーのリクエストを表すDTO
type ReqZanKaiSummary struct {
	P_no      string `json:"p_no"`      // p_no
	P_sd_date string `json:"p_sd_date"` // システム日付
	SJsonOfmt string `json:"sJsonOfmt"` // JSON出力フォーマット
	SCLMID    string `json:"sCLMID"`    // 機能ID, CLMZanKaiSummary
}
