// business_functions/req_zan_shinki_kano_ijiritu.go
package business_functions

// ReqZanShinkiKanoIjiritu は建余力＆本日維持率のリクエストを表すDTO
type ReqZanShinkiKanoIjiritu struct {
	P_no       string `json:"p_no"`       // p_no
	P_sd_date  string `json:"p_sd_date"`  // システム日付
	SJsonOfmt  string `json:"sJsonOfmt"`  // JSON出力フォーマット
	SCLMID     string `json:"sCLMID"`     // 機能ID, CLMZanShinkiKanoIjiritu
	SIssueCode string `json:"sIssueCode"` // 銘柄コード 未使用
	SSizyouC   string `json:"sSizyouC"`   // 市場 未使用
}
