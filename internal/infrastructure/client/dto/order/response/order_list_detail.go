// response/order_list_detail.go
package response

// ResOrderListDetail は注文約定一覧（詳細）のレスポンスを表すDTO
type ResOrderListDetail struct {
	P_no                     string                    `json:"p_no"`                      // p_no
	CLMID                    string                    `json:"sCLMID"`                    // 機能ID, CLMOrderListDetail
	ResultCode               string                    `json:"sResultCode"`               // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText               string                    `json:"sResultText"`               // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode              string                    `json:"sWarningCode"`              // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText              string                    `json:"sWarningText"`              // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	OrderNumber              string                    `json:"sOrderNumber"`              // 注文番号, CLMKabuNewOrder.sOrderNumber 参照
	EigyouDay                string                    `json:"sEigyouDay"`                // 営業日, CLMKabuNewOrder.sEigyouDay 参照
	IssueCode                string                    `json:"sIssueCode"`                // 銘柄コード, CLMKabuNewOrder.sIssueCode 参照
	OrderSizyouC             string                    `json:"sOrderSizyouC"`             // 市場, CLMKabuNewOrder.sSizyouC 参照
	OrderBaibaiKubun         string                    `json:"sOrderBaibaiKubun"`         // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	GenkinSinyouKubun        string                    `json:"sGenkinSinyouKubun"`        // 現金信用区分, CLMKabuNewOrder.sGenkinShinyouKubun 参照
	OrderBensaiKubun         string                    `json:"sOrderBensaiKubun"`         // 弁済区分, CLMOrderList.sOrderBensaiKubun 参照
	OrderCondition           string                    `json:"sOrderCondition"`           // 執行条件, CLMKabuNewOrder.sCondition 参照
	OrderOrderPriceKubun     string                    `json:"sOrderOrderPriceKubun"`     // 注文値段区分, CLMOrderList.sOrderOrderPriceKubun 参照
	OrderOrderPrice          string                    `json:"sOrderOrderPrice"`          // 注文単価, CLMOrderList.sOrderOrderPrice 参照
	OrderOrderSuryou         string                    `json:"sOrderOrderSuryou"`         // 注文株数, CLMOrderList.sOrderOrderSuryou 参照
	OrderCurrentSuryou       string                    `json:"sOrderCurrentSuryou"`       // 有効株数, CLMOrderList.sOrderCurrentSuryou 参照
	OrderStatusCode          string                    `json:"sOrderStatusCode"`          // 状態コード, CLMOrderList.sOrderStatusCode 参照
	OrderStatus              string                    `json:"sOrderStatus"`              // 状態名称, CLMOrderList.sOrderStatus 参照
	OrderOrderDateTime       string                    `json:"sOrderOrderDateTime"`       // 注文日付, CLMOrderList.sOrderOrderDateTime 参照
	OrderOrderExpireDay      string                    `json:"sOrderOrderExpireDay"`      // 有効期限, CLMOrderList.sOrderOrderExpireDay 参照
	Channel                  string                    `json:"sChannel"`                  // チャネル
	GenbutuZyoutoekiKazeiC   string                    `json:"sGenbutuZyoutoekiKazeiC"`   // 現物口座区分, CLMOrderList.sOrderZyoutoekiKazeiC 参照
	SinyouZyoutoekiKazeiC    string                    `json:"sSinyouZyoutoekiKazeiC"`    // 建玉口座区分, 1：特定, 3：一般
	GyakusasiOrderType       string                    `json:"sGyakusasiOrderType"`       // 逆指値注文種別, 0：通常, 1：逆指値, 2：通常＋逆指値
	GyakusasiZyouken         string                    `json:"sGyakusasiZyouken"`         // 逆指値条件
	GyakusasiKubun           string                    `json:"sGyakusasiKubun"`           // 逆指値値段区分, CLMOrderList.sOrderGyakusasiKubun 参照
	GyakusasiPrice           string                    `json:"sGyakusasiPrice"`           // 逆指値値段
	TriggerType              string                    `json:"sTriggerType"`              // トリガータイプ, CLMOrderList.sOrderTriggerType 参照
	TriggerTime              string                    `json:"sTriggerTime"`              // トリガー日時, YYYYMMDDHHMMSS, 00000000000000
	UkewatasiDay             string                    `json:"sUkewatasiDay"`             // 受渡日, YYYYMMDD, 00000000
	YakuzyouPrice            string                    `json:"sYakuzyouPrice"`            // 約定単価
	YakuzyouSuryou           string                    `json:"sYakuzyouSuryou"`           // 約定株数
	BaiBaiDaikin             string                    `json:"sBaiBaiDaikin"`             // 売買代金
	UtidekiKubun             string                    `json:"sUtidekiKubun"`             // 内出来区分, CLMOrderList.sOrderUtidekiKbn 参照
	GaisanDaikin             string                    `json:"sGaisanDaikin"`             // 概算代金
	BaiBaiTesuryo            string                    `json:"sBaiBaiTesuryo"`            // 手数料
	Shouhizei                string                    `json:"sShouhizei"`                // 消費税
	TatebiType               string                    `json:"sTatebiType"`               // 建日種類, CLMOrderList.sOrderTatebiType 参照
	SizyouErrorCode          string                    `json:"sSizyouErrorCode"`          // 取引所エラー等理由コード, ""：正常
	Zougen                   string                    `json:"sZougen"`                   // リバース増減値, 未使用
	OrderAcceptTime          string                    `json:"sOrderAcceptTime"`          // 取引所受付／エラー時刻, YYYYMMDDHHMMSS, 00000000000000
	OrderExpireDayLimit      string                    `json:"sOrderExpireDayLimit"`      // 注文失効日付, YYYYMMDD
	YakuzyouSikkouList       []ResYakuzyouSikkou       `json:"aYakuzyouSikkouList"`       // 約定失効リスト
	KessaiOrderTategyokuList []ResKessaiOrderTategyoku `json:"aKessaiOrderTategyokuList"` // 決済注文建株指定リスト
}

// ResYakuzyouSikkou 約定失効リストの要素
type ResYakuzyouSikkou struct {
	YakuzyouWarningCode string `json:"sYakuzyouWarningCode"` // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	YakuzyouWarningText string `json:"sYakuzyouWarningText"` // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	YakuzyouSuryou      string `json:"sYakuzyouSuryou"`      // 約定数量
	YakuzyouPrice       string `json:"sYakuzyouPrice"`       // 約定価格
	YakuzyouDate        string `json:"sYakuzyouDate"`        // 約定日時, YYYYMMDDHHMMSS, 00000000000000
}

// ResKessaiOrderTategyoku 決済注文建株指定リストの要素
type ResKessaiOrderTategyoku struct {
	KessaiWarningCode    string `json:"sKessaiWarningCode"`    // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	KessaiWarningText    string `json:"sKessaiWarningText"`    // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	KessaiTatebiZyuni    string `json:"sKessaiTatebiZyuni"`    // 順位
	KessaiTategyokuDay   string `json:"sKessaiTategyokuDay"`   // 建日, YYYYMMDD, 00000000
	KessaiTategyokuPrice string `json:"sKessaiTategyokuPrice"` // 建単価
	KessaiOrderSuryo     string `json:"sKessaiOrderSuryo"`     // 返済注文株数
	KessaiYakuzyouSuryo  string `json:"sKessaiYakuzyouSuryo"`  // 約定株数
	KessaiYakuzyouPrice  string `json:"sKessaiYakuzyouPrice"`  // 約定単価
	KessaiTateTesuryou   string `json:"sKessaiTateTesuryou"`   // 建手数料
	KessaiZyunHibu       string `json:"sKessaiZyunHibu"`       // 順日歩
	KessaiGyakuhibu      string `json:"sKessaiGyakuhibu"`      // 逆日歩
	KessaiKakikaeryou    string `json:"sKessaiKakikaeryou"`    // 書換料
	KessaiKanrihi        string `json:"sKessaiKanrihi"`        // 管理費
	KessaiKasikaburyou   string `json:"sKessaiKasikaburyou"`   // 貸株料
	KessaiSonota         string `json:"sKessaiSonota"`         // その他
	KessaiSoneki         string `json:"sKessaiSoneki"`         // 決済損益/受渡代金
}
