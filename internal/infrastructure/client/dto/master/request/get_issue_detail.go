// internal/infrastructure/client/dto/master/request/get_issue_detail.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetIssueDetailRequest は、銘柄詳細情報問合取得のリクエストを表すDTO
type GetIssueDetailRequest struct {
	dto.RequestBase         // 共通フィールド
	CLMID            string `json:"sCLMID"`           // 機能ID (固定値: "CLMMfdsGetIssueDetail")
	TargetIssueCodes string `json:"sTargetIssueCode"` // 対象銘柄コード (カンマ区切りで複数指定可能)
}
