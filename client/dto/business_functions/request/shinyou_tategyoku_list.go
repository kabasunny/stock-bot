// business_functions/req_shinyou_tategyoku_list.go
package business_functions

// ReqShinyouTategyokuList は信用建玉一覧のリクエストを表すDTO
type ReqShinyouTategyokuList struct {
	P_no       string `json:"p_no"`       // p_no
	P_sd_date  string `json:"p_sd_date"`  // システム日付
	SJsonOfmt  string `json:"sJsonOfmt"`  // JSON出力フォーマット
	SCLMID     string `json:"sCLMID"`     // 機能ID, CLMShinyouTategyokuList
	SIssueCode string `json:"sIssueCode"` // 銘柄コード, 指定あり：指定１銘柄, 指定なし：全保有銘柄
}
