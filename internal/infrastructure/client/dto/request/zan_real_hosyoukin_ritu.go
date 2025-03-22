// request/zan_real_hosyoukin_ritu.go
package request

// ReqZanRealHosyoukinRitu はリアル保証金率のリクエストを表すDTO
type ReqZanRealHosyoukinRitu struct {
	P_no      string `json:"p_no"`      // p_no
	P_sd_date string `json:"p_sd_date"` // システム日付
	SJsonOfmt string `json:"sJsonOfmt"` // JSON出力フォーマット
	SCLMID    string `json:"sCLMID"`    // 機能ID, CLMZanRealHosyoukinRitu
}
