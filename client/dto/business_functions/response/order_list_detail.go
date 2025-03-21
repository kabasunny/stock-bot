// business_functions/res_order_list_detail.go
package business_functions

// ResOrderListDetail は注文約定一覧（詳細）のレスポンスを表すDTO
type ResOrderListDetail struct {
	P_no                      string                    `json:"p_no"`                      // p_no
	SCLMID                    string                    `json:"sCLMID"`                    // 機能ID, CLMOrderListDetail
	SResultCode               string                    `json:"sResultCode"`               // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText               string                    `json:"sResultText"`               // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode              string                    `json:"sWarningCode"`              // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText              string                    `json:"sWarningText"`              // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SOrderNumber              string                    `json:"sOrderNumber"`              // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	SEigyouDay                string                    `json:"sEigyouDay"`                // 営業日, CLMKabuNewOrder.sEigyouDay 参照
	SIssueCode                string                    `json:"sIssueCode"`                // 銘柄コード, CLMKabuNewOrder.sIssueCode 参照
	SOrderSizyouC             string                    `json:"sOrderSizyouC"`             // 市場, CLMKabuNewOrder.sSizyouC 参照
	SOrderBaibaiKubun         string                    `json:"sOrderBaibaiKubun"`         // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	SGenkinSinyouKubun        string                    `json:"sGenkinSinyouKubun"`        // 現金信用区分, CLMKabuNewOrder.sGenkinShinyouKubun 参照
	SOrderBensaiKubun         string                    `json:"sOrderBensaiKubun"`         // 弁済区分, CLMOrderList.sOrderBensaiKubun 参照
	SOrderCondition           string                    `json:"sOrderCondition"`           // 執行条件, CLMKabuNewOrder.sCondition 参照
	SOrderOrderPriceKubun     string                    `json:"sOrderOrderPriceKubun"`     // 注文値段区分, CLMOrderList.sOrderOrderPriceKubun 参照
	SOrderOrderPrice          string                    `json:"sOrderOrderPrice"`          // 注文単価, CLMOrderList.sOrderOrderPrice 参照
	SOrderOrderSuryou         string                    `json:"sOrderOrderSuryou"`         // 注文株数, CLMOrderList.sOrderOrderSuryou 参照
	SOrderCurrentSuryou       string                    `json:"sOrderCurrentSuryou"`       // 有効株数, CLMOrderList.sOrderCurrentSuryou 参照
	SOrderStatusCode          string                    `json:"sOrderStatusCode"`          // 状態コード, CLMOrderList.sOrderStatusCode 参照
	SOrderStatus              string                    `json:"sOrderStatus"`              // 状態名称, CLMOrderList.sOrderStatus 参照
	SOrderOrderDateTime       string                    `json:"sOrderOrderDateTime"`       // 注文日付, CLMOrderList.sOrderOrderDateTime 参照
	SOrderOrderExpireDay      string                    `json:"sOrderOrderExpireDay"`      // 有効期限, CLMOrderList.sOrderOrderExpireDay 参照
	SChannel                  string                    `json:"sChannel"`                  // チャネル
	SGenbutuZyoutoekiKazeiC   string                    `json:"sGenbutuZyoutoekiKazeiC"`   // 現物口座区分, CLMOrderList.sOrderZyoutoekiKazeiC 参照
	SSinyouZyoutoekiKazeiC    string                    `json:"sSinyouZyoutoekiKazeiC"`    // 建玉口座区分, 1：特定, 3：一般
	SGyakusasiOrderType       string                    `json:"sGyakusasiOrderType"`       // 逆指値注文種別, 0：通常, 1：逆指値, 2：通常＋逆指値
	SGyakusasiZyouken         string                    `json:"sGyakusasiZyouken"`         // 逆指値条件
	SGyakusasiKubun           string                    `json:"sGyakusasiKubun"`           // 逆指値値段区分, CLMOrderList.sOrderGyakusasiKubun 参照
	SGyakusasiPrice           string                    `json:"sGyakusasiPrice"`           // 逆指値値段
	STriggerType              string                    `json:"sTriggerType"`              // トリガータイプ, CLMOrderList.sOrderTriggerType 参照
	STriggerTime              string                    `json:"sTriggerTime"`              // トリガー日時, YYYYMMDDHHMMSS, 00000000000000
	SUkewatasiDay             string                    `json:"sUkewatasiDay"`             // 受渡日, YYYYMMDD, 00000000
	SYakuzyouPrice            string                    `json:"sYakuzyouPrice"`            // 約定単価
	SYakuzyouSuryou           string                    `json:"sYakuzyouSuryou"`           // 約定株数
	SBaiBaiDaikin             string                    `json:"sBaiBaiDaikin"`             // 売買代金
	SUtidekiKubun             string                    `json:"sUtidekiKubun"`             // 内出来区分, CLMOrderList.sOrderUtidekiKbn 参照
	SGaisanDaikin             string                    `json:"sGaisanDaikin"`             // 概算代金
	SBaiBaiTesuryo            string                    `json:"sBaiBaiTesuryo"`            // 手数料
	SShouhizei                string                    `json:"sShouhizei"`                // 消費税
	STatebiType               string                    `json:"sTatebiType"`               // 建日種類, CLMOrderList.sOrderTatebiType 参照
	SSizyouErrorCode          string                    `json:"sSizyouErrorCode"`          // 取引所エラー等理由コード, ""：正常
	SZougen                   string                    `json:"sZougen"`                   // リバース増減値, 未使用
	SOrderAcceptTime          string                    `json:"sOrderAcceptTime"`          // 取引所受付／エラー時刻, YYYYMMDDHHMMSS, 00000000000000
	SOrderExpireDayLimit      string                    `json:"sOrderExpireDayLimit"`      // 注文失効日付, YYYYMMDD
	AYakuzyouSikkouList       []ResYakuzyouSikkou       `json:"aYakuzyouSikkouList"`       // 約定失効リスト
	AKessaiOrderTategyokuList []ResKessaiOrderTategyoku `json:"aKessaiOrderTategyokuList"` // 決済注文建株指定リスト
}

// ResYakuzyouSikkou 約定失効リストの要素
type ResYakuzyouSikkou struct {
	SYakuzyouWarningCode string `json:"sYakuzyouWarningCode"` // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SYakuzyouWarningText string `json:"sYakuzyouWarningText"` // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SYakuzyouSuryou      string `json:"sYakuzyouSuryou"`      // 約定数量
	SYakuzyouPrice       string `json:"sYakuzyouPrice"`       // 約定価格
	SYakuzyouDate        string `json:"sYakuzyouDate"`        // 約定日時, YYYYMMDDHHMMSS, 00000000000000
}

// ResKessaiOrderTategyoku 決済注文建株指定リストの要素
type ResKessaiOrderTategyoku struct {
	SKessaiWarningCode    string `json:"sKessaiWarningCode"`    // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SKessaiWarningText    string `json:"sKessaiWarningText"`    // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SKessaiTatebiZyuni    string `json:"sKessaiTatebiZyuni"`    // 順位
	SKessaiTategyokuDay   string `json:"sKessaiTategyokuDay"`   // 建日, YYYYMMDD, 00000000
	SKessaiTategyokuPrice string `json:"sKessaiTategyokuPrice"` // 建単価
	SKessaiOrderSuryo     string `json:"sKessaiOrderSuryo"`     // 返済注文株数
	SKessaiYakuzyouSuryo  string `json:"sKessaiYakuzyouSuryo"`  // 約定株数
	SKessaiYakuzyouPrice  string `json:"sKessaiYakuzyouPrice"`  // 約定単価
	SKessaiTateTesuryou   string `json:"sKessaiTateTesuryou"`   // 建手数料
	SKessaiZyunHibu       string `json:"sKessaiZyunHibu"`       // 順日歩
	SKessaiGyakuhibu      string `json:"sKessaiGyakuhibu"`      // 逆日歩
	SKessaiKakikaeryou    string `json:"sKessaiKakikaeryou"`    // 書換料
	SKessaiKanrihi        string `json:"sKessaiKanrihi"`        // 管理費
	SKessaiKasikaburyou   string `json:"sKessaiKasikaburyou"`   // 貸株料
	SKessaiSonota         string `json:"sKessaiSonota"`         // その他
	SKessaiSoneki         string `json:"sKessaiSoneki"`         // 決済損益/受渡代金
}
