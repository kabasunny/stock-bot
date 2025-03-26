// request/shinyou_tategyoku_list.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqShinyouTategyokuList は信用建玉一覧のリクエストを表すDTO
type ReqShinyouTategyokuList struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`     // 機能ID, CLMShinyouTategyokuList
	IssueCode       string `json:"sIssueCode"` // 銘柄コード, 指定あり：指定１銘柄, 指定なし：全保有銘柄
}
