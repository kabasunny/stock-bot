package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// 蓄積情報問合取得 リクエスト
// internal/infrastructure/client/dto/price/request/get_price_info_history.go

type ReqGetPriceInfoHistory struct {
	dto.RequestBase
	SCLMID     string `json:"sCLMID"`             // 機能ID (CLMMfdsGetMarketPriceHistory)
	SIssueCode string `json:"sIssueCode"`         // 銘柄コード
	SSizyouC   string `json:"sSizyouC,omitempty"` // 市場コード (省略可能, デフォルト="00":東証)
}
