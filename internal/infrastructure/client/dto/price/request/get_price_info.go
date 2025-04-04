package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// 時価情報問合取得 リクエスト
// internal/infrastructure/client/dto/price/request/get_price_info.go
type ReqGetPriceInfo struct {
	dto.RequestBase
	CLMID           string `json:"sCLMID"`           // 機能ID (CLMMfdsGetMarketPrice)
	TargetIssueCode string `json:"sTargetIssueCode"` // 対象銘柄コード (カンマ区切り)
	TargetColumn    string `json:"sTargetColumn"`    // 対象情報コード (カンマ区切り)
}
