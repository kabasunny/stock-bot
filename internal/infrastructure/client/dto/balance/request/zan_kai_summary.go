// request/zan_kai_summary.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanKaiSummary は可能額サマリーのリクエストを表すDTO
type ReqZanKaiSummary struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"` // 機能ID, CLMZanKaiSummary
}
