// request/zan_kai_genbutu_kaituke_syousai.go
package request

// ReqZanKaiGenbutuKaitukeSyousai は現物株式買付可能額詳細のリクエストを表すDTO
type ReqZanKaiGenbutuKaitukeSyousai struct {
	P_no         string `json:"p_no"`         // p_no
	P_sd_date    string `json:"p_sd_date"`    // システム日付
	SJsonOfmt    string `json:"sJsonOfmt"`    // JSON出力フォーマット
	SCLMID       string `json:"sCLMID"`       // 機能ID, CLMZanKaiGenbutuKaitukeSyousai
	SHitukeIndex string `json:"sHitukeIndex"` // 日付インデックス, 3:第4営業日, 4:第5営業日, 5:第6営業日
}
