// request/zan_shinki_kano_ijiritu.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqZanShinkiKanoIjiritu は建余力＆本日維持率のリクエストを表すDTO
type ReqZanShinkiKanoIjiritu struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`     // 機能ID, CLMZanShinkiKanoIjiritu
	IssueCode       string `json:"sIssueCode"` // 銘柄コード 未使用
	SizyouC         string `json:"sSizyouC"`   // 市場 未使用
}
