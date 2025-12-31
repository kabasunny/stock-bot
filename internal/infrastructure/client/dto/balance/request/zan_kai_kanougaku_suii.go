// request/zan_kai_kanougaku_suii.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanKaiKanougakuSuii は可能額推移のリクエストを表すDTO
type ReqZanKaiKanougakuSuii struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"` // 機能ID, CLMZanKaiKanougakuSuii
}
