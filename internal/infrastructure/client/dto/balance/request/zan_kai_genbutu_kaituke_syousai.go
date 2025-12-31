// request/zan_kai_genbutu_kaituke_syousai.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanKaiGenbutuKaitukeSyousai は現物株式買付可能額詳細のリクエストを表すDTO
type ReqZanKaiGenbutuKaitukeSyousai struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`       // 機能ID, CLMZanKaiGenbutuKaitukeSyousai
	HitukeIndex     string `json:"sHitukeIndex"` // 日付インデックス, 3:第4営業日, 4:第5営業日, 5:第6営業日
}
