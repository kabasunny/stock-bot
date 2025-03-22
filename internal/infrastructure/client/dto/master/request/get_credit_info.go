// internal/infrastructure/client/dto/master/request/get_credit_info.go

package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetCreditInfoRequest は、信用残情報問合取得のリクエストを表すDTO
type GetCreditInfoRequest struct {
	dto.RequestBase         // 共通フィールド
	CLMID            string `json:"sCLMID"`           // 機能ID (固定値: "CLMMfdsGetShinyouZan")
	TargetIssueCodes string `json:"sTargetIssueCode"` // 対象銘柄コード (カンマ区切りで複数指定可能)
}
