// business_functions/req_zan_kai_sinyou_sinkidate_syousai.go
package business_functions

// ReqZanKaiSinyouSinkidateSyousai は信用新規建て可能額詳細のリクエストを表すDTO
type ReqZanKaiSinyouSinkidateSyousai struct {
	P_no         string `json:"p_no"`         // p_no
	P_sd_date    string `json:"p_sd_date"`    // システム日付
	SJsonOfmt    string `json:"sJsonOfmt"`    // JSON出力フォーマット
	SCLMID       string `json:"sCLMID"`       // 機能ID, CLMZanKaiSinyouSinkidateSyousai
	SHitukeIndex string `json:"sHitukeIndex"` // 日付インデックス, 0:第1営業日, 1:第2営業日, 2:第3営業日, 3:第4営業日, 4:第5営業日, 5:第6営業日
}
