// request/new_order.go
package request

// ReqNewOrder は株式新規注文のリクエストを表すDTO
type ReqNewOrder struct {
	P_no                      string          `json:"p_no"`                      // p_no
	P_sd_date                 string          `json:"p_sd_date"`                 // システム日付
	SJsonOfmt                 string          `json:"sJsonOfmt"`                 // JSON出力フォーマット
	SCLMID                    string          `json:"sCLMID"`                    // 機能ID, CLMKabuNewOrder
	SZyoutoekiKazeiC          string          `json:"sZyoutoekiKazeiC"`          // 譲渡益課税区分, 1：特定, 3：一般, 5：NISA, 6：N成長
	SIssueCode                string          `json:"sIssueCode"`                // 銘柄コード
	SSizyouC                  string          `json:"sSizyouC"`                  // 市場, 00：東証
	SBaibaiKubun              string          `json:"sBaibaiKubun"`              // 売買区分, 1：売, 3：買, 5：現渡, 7：現引
	SCondition                string          `json:"sCondition"`                // 執行条件, 0：指定なし, 2：寄付, 4：引け, 6：不成
	SOrderPrice               string          `json:"sOrderPrice"`               // 注文値段, *：指定なし, 0：成行
	SOrderSuryou              string          `json:"sOrderSuryou"`              // 注文株数
	SGenkinShinyouKubun       string          `json:"sGenkinShinyouKubun"`       // 現金信用区分, 0：現物, 2：新規(制度信用6ヶ月), 4：返済(制度信用6ヶ月), 6：新規(一般信用6ヶ月), 8：返済(一般信用6ヶ月)
	SOrderExpireDay           string          `json:"sOrderExpireDay"`           // 注文期日, 0：当日
	SGyakusasiOrderType       string          `json:"sGyakusasiOrderType"`       // 逆指値注文種別, 0：通常, 1：逆指値, 2：通常＋逆指値
	SGyakusasiZyouken         string          `json:"sGyakusasiZyouken"`         // 逆指値条件, 0：指定なし
	SGyakusasiPrice           string          `json:"sGyakusasiPrice"`           // 逆指値値段, *：指定なし, 0：成行
	STatebiType               string          `json:"sTatebiType"`               // 建日種類, *：指定なし（現物または新規）, 1：個別指定, 2：建日順, 3：単価益順, 4：単価損順
	STategyokuZyoutoekiKazeiC string          `json:"sTategyokuZyoutoekiKazeiC"` // 建玉譲渡益課税区分, *：現引、現渡以外の取引, 1：特定, 3：一般
	SSecondPassword           string          `json:"sSecondPassword"`           // 第二パスワード
	ACLMKabuHensaiData        []ReqHensaiData `json:"aCLMKabuHensaiData"`        // 返済リスト
}

// ReqHensaiData　返済リスト
type ReqHensaiData struct {
	STategyokuNumber string `json:"sTategyokuNumber"` // 	新規建玉番号
	STatebiZyuni     string `json:"sTatebiZyuni"`     // 建日順位
	SOrderSuryou     string `json:"sOrderSuryou"`     // 注文数量
}
