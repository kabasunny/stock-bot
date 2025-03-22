// internal/infrastructure/client/dto/master/request/get_margin_info.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetMarginInfoRequest は、証金残情報問合取得のリクエストを表すDTO
type GetMarginInfoRequest struct {
	dto.RequestBase         // 共通フィールド
	CLMID            string `json:"sCLMID"`           // 機能ID (固定値: "CLMMfdsGetSyoukinZan")
	TargetIssueCodes string `json:"sTargetIssueCode"` // 対象銘柄コード (カンマ区切りで複数指定可能)
}
