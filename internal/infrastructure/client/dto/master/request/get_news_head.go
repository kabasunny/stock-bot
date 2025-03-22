// internal/infrastructure/client/dto/master/request/get_news_head.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetNewsHeadRequest は、ニュースヘッダー問合取得のリクエストを表すDTO
type GetNewsHeadRequest struct {
	dto.RequestBase        // 共通フィールド
	CLMID           string `json:"sCLMID"`               // 機能ID (固定値: "CLMMfdsGetNewsHead")
	Category        string `json:"p_CG,omitempty"`       // カテゴリコード (任意)
	Issue           string `json:"p_IS,omitempty"`       // 銘柄コード (任意)
	FromDate        string `json:"p_DT_FROM,omitempty"`  // 日付(From) (YYYYMMDD形式、任意)
	ToDate          string `json:"p_DT_TO,omitempty"`    // 日付(To) (YYYYMMDD形式、任意)
	Offset          string `json:"p_REC_OFST,omitempty"` // レコード取得位置 (デフォルト: 0、任意)
	Limit           string `json:"p_REC_LIMT,omitempty"` // レコード取得件数最大 (デフォルト: 100、任意)
}
