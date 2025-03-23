// request/order_list.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqOrderList は注文一覧のリクエストを表すDTO
type ReqOrderList struct {
	dto.RequestBase           // 共通フィールドを埋め込む
	CLMID              string `json:"sCLMID"`              // 機能ID, CLMOrderList
	IssueCode          string `json:"sIssueCode"`          // 銘柄コード  (任意), 指定あり：指定１銘柄, 指定なし：全保有銘柄
	SikkouDay          string `json:"sSikkouDay"`          // 注文執行予定日（営業日） (任意), 指定あり：指定１営業日, 指定なし：全保有営業日
	OrderSyoukaiStatus string `json:"sOrderSyoukaiStatus"` // 注文照会状態 (任意), ""：指定なし, 1：未約定, 2：全部約定, 3：一部約定, 4：訂正取消(可能な注文）, 5：未約定+一部約定
}
