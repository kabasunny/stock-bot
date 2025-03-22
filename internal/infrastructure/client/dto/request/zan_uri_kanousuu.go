// request/zan_uri_kanousuu.go
package request

// ReqZanUriKanousuu は売却可能数量のリクエストを表すDTO
type ReqZanUriKanousuu struct {
	P_no       string `json:"p_no"`       // p_no
	P_sd_date  string `json:"p_sd_date"`  // システム日付
	SJsonOfmt  string `json:"sJsonOfmt"`  // JSON出力フォーマット
	SCLMID     string `json:"sCLMID"`     // 機能ID, CLMZanUriKanousuu
	SIssueCode string `json:"sIssueCode"` // 銘柄コード
}
