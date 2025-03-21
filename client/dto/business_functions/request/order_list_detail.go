// business_functions/req_order_list_detail.go
package business_functions

// ReqOrderListDetail は注文約定一覧（詳細）のリクエストを表すDTO
type ReqOrderListDetail struct {
	P_no         string `json:"p_no"`         // p_no
	P_sd_date    string `json:"p_sd_date"`    // システム日付
	SJsonOfmt    string `json:"sJsonOfmt"`    // JSON出力フォーマット
	SCLMID       string `json:"sCLMID"`       // 機能ID, CLMOrderListDetail
	SOrderNumber string `json:"sOrderNumber"` // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	SEigyouDay   string `json:"sEigyouDay"`   // 営業日, CLMKabuNewOrder.sEigyouDay 参照
}
