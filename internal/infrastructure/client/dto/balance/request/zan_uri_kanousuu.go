// request/zan_uri_kanousuu.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanUriKanousuu は売却可能数量のリクエストを表すDTO
type ReqZanUriKanousuu struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`     // 機能ID, CLMZanUriKanousuu
	IssueCode       string `json:"sIssueCode"` // 銘柄コード
}
