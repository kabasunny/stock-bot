// request/zan_real_hosyoukin_ritu.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanRealHosyoukinRitu はリアル保証金率のリクエストを表すDTO
type ReqZanRealHosyoukinRitu struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"` // 機能ID, CLMZanRealHosyoukinRitu
}
