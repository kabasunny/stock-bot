// request/zan_kai_sinyou_sinkidate_syousai.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanKaiSinyouSinkidateSyousai は信用新規建て可能額詳細のリクエストを表すDTO
type ReqZanKaiSinyouSinkidateSyousai struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`       // 機能ID, CLMZanKaiSinyouSinkidateSyousai
	HitukeIndex     string `json:"sHitukeIndex"` // 日付インデックス, 0:第1営業日, 1:第2営業日, 2:第3営業日, 3:第4営業日, 4:第5営業日, 5:第6営業日
}
