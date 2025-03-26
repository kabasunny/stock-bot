// request/zan_kai_kanougaku.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanKaiKanougaku は買余力のリクエストを表すDTO
type ReqZanKaiKanougaku struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`     // 機能ID, CLMZanKaiKanougaku
	IssueCode       string `json:"sIssueCode"` // 銘柄コード  未使用
	SizyouC         string `json:"sSizyouC"`   // 市場      未使用
}
