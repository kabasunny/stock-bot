// internal/infrastructure/client/dto/master/request/get_margin_premium_info.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetMarginPremiumInfoRequest は、逆日歩情報問合取得のリクエストを表すDTO
type ReqGetMarginPremiumInfo struct {
	dto.RequestBase         // 共通フィールド
	CLMID            string `json:"sCLMID"`           // 機能ID (固定値: "CLMMfdsGetHibuInfo")
	TargetIssueCodes string `json:"sTargetIssueCode"` // 対象銘柄コード (カンマ区切りで複数指定可能)
}
