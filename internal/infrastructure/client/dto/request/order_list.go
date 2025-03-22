// request/order_list.go
package request

// ReqOrderList は注文一覧のリクエストを表すDTO
type ReqOrderList struct {
	P_no                string `json:"p_no"`                // p_no
	P_sd_date           string `json:"p_sd_date"`           // システム日付
	SJsonOfmt           string `json:"sJsonOfmt"`           // JSON出力フォーマット
	SCLMID              string `json:"sCLMID"`              // 機能ID, CLMOrderList
	SIssueCode          string `json:"sIssueCode"`          // 銘柄コード  (任意), 指定あり：指定１銘柄, 指定なし：全保有銘柄
	SSikkouDay          string `json:"sSikkouDay"`          // 注文執行予定日（営業日） (任意), 指定あり：指定１営業日, 指定なし：全保有営業日
	SOrderSyoukaiStatus string `json:"sOrderSyoukaiStatus"` // 注文照会状態 (任意), ""：指定なし, 1：未約定, 2：全部約定, 3：一部約定, 4：訂正取消(可能な注文）, 5：未約定+一部約定
}
