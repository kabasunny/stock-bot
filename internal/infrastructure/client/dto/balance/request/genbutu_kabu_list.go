// request/genbutu_kabu_list.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqGenbutuKabuList は現物保有銘柄一覧のリクエストを表すDTO
type ReqGenbutuKabuList struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`     // 機能ID, CLMGenbutuKabuList
	IssueCode       string `json:"sIssueCode"` // 銘柄コード, 指定あり：指定１銘柄, 指定なし：全保有銘柄
}
