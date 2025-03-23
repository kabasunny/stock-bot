// response/order_list.go
package response

// ResOrderList は注文一覧のレスポンスを表すDTO
type ResOrderList struct {
	P_no               string     `json:"p_no"`                // p_no
	CLMID              string     `json:"sCLMID"`              // 機能ID, CLMOrderList
	ResultCode         string     `json:"sResultCode"`         // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText         string     `json:"sResultText"`         // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode        string     `json:"sWarningCode"`        // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText        string     `json:"sWarningText"`        // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	IssueCode          string     `json:"sIssueCode"`          // 銘柄コード, 要求設定値
	OrderSyoukaiStatus string     `json:"sOrderSyoukaiStatus"` // 注文照会状態, 要求設定値
	SikkouDay          string     `json:"sSikkouDay"`          // 注文執行予定日, 要求設定値
	OrderList          []ResOrder `json:"aOrderList"`          // 注文リスト
}

// ResOrder 注文リストの要素
type ResOrder struct {
	OrderWarningCode          string `json:"sOrderWarningCode"`          // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	OrderWarningText          string `json:"sOrderWarningText"`          // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	OrderOrderNumber          string `json:"sOrderOrderNumber"`          // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	OrderIssueCode            string `json:"sOrderIssueCode"`            // 銘柄コード, CLMKabuNewOrder.sIssueCode 参照
	OrderSizyouC              string `json:"sOrderSizyouC"`              // 市場, CLMKabuNewOrder.SizyouC 参照
	OrderZyoutoekiKazeiC      string `json:"sOrderZyoutoekiKazeiC"`      // 譲渡益課税区分, CLMKabuNewOrder.sZyoutoekiKazeiC 参照
	GenkinSinyouKubun         string `json:"sGenkinSinyouKubun"`         // 現金信用区分, CLMKabuNewOrder.sGenkinShinyouKubun 参照
	OrderBensaiKubun          string `json:"sOrderBensaiKubun"`          // 弁済区分, 00：なし, 26：制度信用6ヶ月, 29：制度信用無期限, 36：一般信用6ヶ月, 39：一般信用無期限
	OrderBaibaiKubun          string `json:"sOrderBaibaiKubun"`          // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	OrderOrderSuryou          string `json:"sOrderOrderSuryou"`          // 注文株数
	OrderCurrentSuryou        string `json:"sOrderCurrentSuryou"`        // 有効株数, Ｎ≦CLMKabuNewOrder.sOrderSuryou
	OrderOrderPrice           string `json:"sOrderOrderPrice"`           // 注文単価
	OrderCondition            string `json:"sOrderCondition"`            // 執行条件, CLMKabuNewOrder.sCondition 参照
	OrderOrderPriceKubun      string `json:"sOrderOrderPriceKubun"`      // 注文値段区分, " "：未使用, 1：成行, 2：指値, 3：親注文より高い, 4：親注文より低い
	OrderGyakusasiOrderType   string `json:"sOrderGyakusasiOrderType"`   // 逆指値注文種別, CLMKabuNewOrder.sGyakusasiOrderType 参照
	OrderGyakusasiZyouken     string `json:"sOrderGyakusasiZyouken"`     // 逆指値条件
	OrderGyakusasiKubun       string `json:"sOrderGyakusasiKubun"`       // 逆指値値段区分, " "：未使用, 0：成行, 1：指値
	OrderGyakusasiPrice       string `json:"sOrderGyakusasiPrice"`       // 逆指値値段
	OrderTriggerType          string `json:"sOrderTriggerType"`          // トリガータイプ, 0：未トリガー（初期値）
	OrderTatebiType           string `json:"sOrderTatebiType"`           // 建日種類, " "：指定なし, 1：個別指定, 2：建日順, 3：単価益順, 4：単価損順
	OrderZougen               string `json:"sOrderZougen"`               // リバース増減値, 未使用
	OrderYakuzyouSuryo        string `json:"sOrderYakuzyouSuryo"`        // 成立株数
	OrderYakuzyouPrice        string `json:"sOrderYakuzyouPrice"`        // 成立単価
	OrderUtidekiKbn           string `json:"sOrderUtidekiKbn"`           // 内出来区分, " "：約定分割以外, 2：約定分割
	OrderSikkouDay            string `json:"sOrderSikkouDay"`            // 執行日, YYYYMMDD
	OrderStatusCode           string `json:"sOrderStatusCode"`           // 状態コード
	OrderStatus               string `json:"sOrderStatus"`               // 状態名称, 状態コードの名称
	OrderYakuzyouStatus       string `json:"sOrderYakuzyouStatus"`       // 約定ステータス, 0：未約定, 1：一部約定, 2：全部約定, 3：約定中
	OrderOrderDateTime        string `json:"sOrderOrderDateTime"`        // 注文日付, YYYYMMDDHHMMSS, 00000000000000
	OrderOrderExpireDay       string `json:"sOrderOrderExpireDay"`       // 有効期限, YYYYMMDD, 00000000
	OrderKurikosiOrderFlg     string `json:"sOrderKurikosiOrderFlg"`     // 繰越注文フラグ, 0：当日注文, 1：繰越注文, 2：無効
	OrderCorrectCancelKahiFlg string `json:"sOrderCorrectCancelKahiFlg"` // 訂正取消可否フラグ, 0：可(取消、訂正), 1：否, 2：一部可(取消のみ)
	GaisanDaikin              string `json:"sGaisanDaikin"`              // 概算代金
}
