// internal/infrastructure/client/dto/master/request/get_news_body.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetNewsBodyRequest は、ニュースボディー問合取得のリクエストを表すDTO
type GetNewsBodyRequest struct {
	dto.RequestBase        // 共通フィールド
	CLMID           string `json:"sCLMID"` // 機能ID (固定値: "CLMMfdsGetNewsBody")
	NewsID          string `json:"p_ID"`   // ニュースID
}
