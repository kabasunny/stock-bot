// response/order_list.go
package response

// ResOrderList は注文一覧のレスポンスを表すDTO
type ResOrderList struct {
	P_no                string     `json:"p_no"`                // p_no
	SCLMID              string     `json:"sCLMID"`              // 機能ID, CLMOrderList
	SResultCode         string     `json:"sResultCode"`         // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText         string     `json:"sResultText"`         // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode        string     `json:"sWarningCode"`        // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText        string     `json:"sWarningText"`        // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode          string     `json:"sIssueCode"`          // 銘柄コード, 要求設定値
	SOrderSyoukaiStatus string     `json:"sOrderSyoukaiStatus"` // 注文照会状態, 要求設定値
	SSikkouDay          string     `json:"sSikkouDay"`          // 注文執行予定日, 要求設定値
	AOrderList          []ResOrder `json:"aOrderList"`          // 注文リスト
}

// ResOrder 注文リストの要素
type ResOrder struct {
	SOrderWarningCode          string `json:"sOrderWarningCode"`          // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SOrderWarningText          string `json:"sOrderWarningText"`          // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SOrderOrderNumber          string `json:"sOrderOrderNumber"`          // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	SOrderIssueCode            string `json:"sOrderIssueCode"`            // 銘柄コード, CLMKabuNewOrder.sIssueCode 参照
	SOrderSizyouC              string `json:"sOrderSizyouC"`              // 市場, CLMKabuNewOrder.SizyouC 参照
	SOrderZyoutoekiKazeiC      string `json:"sOrderZyoutoekiKazeiC"`      // 譲渡益課税区分, CLMKabuNewOrder.sZyoutoekiKazeiC 参照
	SGenkinSinyouKubun         string `json:"sGenkinSinyouKubun"`         // 現金信用区分, CLMKabuNewOrder.sGenkinShinyouKubun 参照
	SOrderBensaiKubun          string `json:"sOrderBensaiKubun"`          // 弁済区分, 00：なし, 26：制度信用6ヶ月, 29：制度信用無期限, 36：一般信用6ヶ月, 39：一般信用無期限
	SOrderBaibaiKubun          string `json:"sOrderBaibaiKubun"`          // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	SOrderOrderSuryou          string `json:"sOrderOrderSuryou"`          // 注文株数
	SOrderCurrentSuryou        string `json:"sOrderCurrentSuryou"`        // 有効株数, Ｎ≦CLMKabuNewOrder.sOrderSuryou
	SOrderOrderPrice           string `json:"sOrderOrderPrice"`           // 注文単価
	SOrderCondition            string `json:"sOrderCondition"`            // 執行条件, CLMKabuNewOrder.sCondition 参照
	SOrderOrderPriceKubun      string `json:"sOrderOrderPriceKubun"`      // 注文値段区分, " "：未使用, 1：成行, 2：指値, 3：親注文より高い, 4：親注文より低い
	SOrderGyakusasiOrderType   string `json:"sOrderGyakusasiOrderType"`   // 逆指値注文種別, CLMKabuNewOrder.sGyakusasiOrderType 参照
	SOrderGyakusasiZyouken     string `json:"sOrderGyakusasiZyouken"`     // 逆指値条件
	SOrderGyakusasiKubun       string `json:"sOrderGyakusasiKubun"`       // 逆指値値段区分, " "：未使用, 0：成行, 1：指値
	SOrderGyakusasiPrice       string `json:"sOrderGyakusasiPrice"`       // 逆指値値段
	SOrderTriggerType          string `json:"sOrderTriggerType"`          // トリガータイプ, 0：未トリガー（初期値）
	SOrderTatebiType           string `json:"sOrderTatebiType"`           // 建日種類, " "：指定なし, 1：個別指定, 2：建日順, 3：単価益順, 4：単価損順
	SOrderZougen               string `json:"sOrderZougen"`               // リバース増減値, 未使用
	SOrderYakuzyouSuryo        string `json:"sOrderYakuzyouSuryo"`        // 成立株数
	SOrderYakuzyouPrice        string `json:"sOrderYakuzyouPrice"`        // 成立単価
	SOrderUtidekiKbn           string `json:"sOrderUtidekiKbn"`           // 内出来区分, " "：約定分割以外, 2：約定分割
	SOrderSikkouDay            string `json:"sOrderSikkouDay"`            // 執行日, YYYYMMDD
	SOrderStatusCode           string `json:"sOrderStatusCode"`           // 状態コード
	SOrderStatus               string `json:"sOrderStatus"`               // 状態名称, 状態コードの名称
	SOrderYakuzyouStatus       string `json:"sOrderYakuzyouStatus"`       // 約定ステータス, 0：未約定, 1：一部約定, 2：全部約定, 3：約定中
	SOrderOrderDateTime        string `json:"sOrderOrderDateTime"`        // 注文日付, YYYYMMDDHHMMSS, 00000000000000
	SOrderOrderExpireDay       string `json:"sOrderOrderExpireDay"`       // 有効期限, YYYYMMDD, 00000000
	SOrderKurikosiOrderFlg     string `json:"sOrderKurikosiOrderFlg"`     // 繰越注文フラグ, 0：当日注文, 1：繰越注文, 2：無効
	SOrderCorrectCancelKahiFlg string `json:"sOrderCorrectCancelKahiFlg"` // 訂正取消可否フラグ, 0：可(取消、訂正), 1：否, 2：一部可(取消のみ)
	SGaisanDaikin              string `json:"sGaisanDaikin"`              // 概算代金
}
