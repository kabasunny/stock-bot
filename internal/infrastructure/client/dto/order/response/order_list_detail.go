// internal/infrastructure/client/dto/order/response/order_list_detail.go
package response

// ResOrderListDetail は注文約定一覧（詳細）のレスポンスを表すDTO
type ResOrderListDetail struct {
	P_no                     string                    `json:"p_no"`
	SCLMID                   string                    `json:"sCLMID"`
	ResultCode               string                    `json:"sResultCode"`
	ResultText               string                    `json:"sResultText"`
	WarningCode              string                    `json:"sWarningCode"`
	WarningText              string                    `json:"sWarningText"`
	OrderNumber              string                    `json:"sOrderNumber"`
	EigyouDay                string                    `json:"sEigyouDay"`
	IssueCode                string                    `json:"sIssueCode"`
	OrderSizyouC             string                    `json:"sOrderSizyouC"`
	OrderBaibaiKubun         string                    `json:"sOrderBaibaiKubun"`
	GenkinSinyouKubun        string                    `json:"sGenkinSinyouKubun"`
	OrderBensaiKubun         string                    `json:"sOrderBensaiKubun"`
	OrderCondition           string                    `json:"sOrderCondition"`
	OrderOrderPriceKubun     string                    `json:"sOrderOrderPriceKubun"`
	OrderOrderPrice          string                    `json:"sOrderOrderPrice"`
	OrderOrderSuryou         string                    `json:"sOrderOrderSuryou"`
	OrderCurrentSuryou       string                    `json:"sOrderCurrentSuryou"`
	OrderStatusCode          string                    `json:"sOrderStatusCode"`
	OrderStatus              string                    `json:"sOrderStatus"`
	OrderOrderDateTime       string                    `json:"sOrderOrderDateTime"`
	OrderOrderExpireDay      string                    `json:"sOrderOrderExpireDay"`
	Channel                  string                    `json:"sChannel"`
	GenbutuZyoutoekiKazeiC   string                    `json:"sGenbutuZyoutoekiKazeiC"`
	SinyouZyoutoekiKazeiC    string                    `json:"sSinyouZyoutoekiKazeiC"`
	GyakusasiOrderType       string                    `json:"sGyakusasiOrderType"`
	GyakusasiZyouken         string                    `json:"sGyakusasiZyouken"`
	GyakusasiKubun           string                    `json:"sGyakusasiKubun"`
	GyakusasiPrice           string                    `json:"sGyakusasiPrice"`
	TriggerType              string                    `json:"sTriggerType"`
	TriggerTime              string                    `json:"sTriggerTime"`
	UkewatasiDay             string                    `json:"sUkewatasiDay"`
	YakuzyouPrice            string                    `json:"sYakuzyouPrice"`
	YakuzyouSuryou           string                    `json:"sYakuzyouSuryou"`
	BaiBaiDaikin             string                    `json:"sBaiBaiDaikin"`
	UtidekiKubun             string                    `json:"sUtidekiKubun"`
	GaisanDaikin             string                    `json:"sGaisanDaikin"`
	BaiBaiTesuryo            string                    `json:"sBaiBaiTesuryo"`
	Shouhizei                string                    `json:"sShouhizei"`
	TatebiType               string                    `json:"sTatebiType"`
	SizyouErrorCode          string                    `json:"sSizyouErrorCode"`
	Zougen                   string                    `json:"sZougen"`
	OrderAcceptTime          string                    `json:"sOrderAcceptTime"`
	OrderExpireDayLimit      string                    `json:"sOrderExpireDayLimit"`
	YakuzyouSikkouList       []ResYakuzyouSikkou       `json:"aYakuzyouSikkouList"`
	KessaiOrderTategyokuList []ResKessaiOrderTategyoku `json:"aKessaiOrderTategyokuList"`
}

// ResKessaiOrderTategyoku 決済注文建株指定リストの要素 (変更なし)
type ResKessaiOrderTategyoku struct {
	KessaiWarningCode    string `json:"sKessaiWarningCode"`
	KessaiWarningText    string `json:"sKessaiWarningText"`
	KessaiTatebiZyuni    string `json:"sKessaiTatebiZyuni"`
	KessaiTategyokuDay   string `json:"sKessaiTategyokuDay"`
	KessaiTategyokuPrice string `json:"sKessaiTategyokuPrice"`
	KessaiOrderSuryo     string `json:"sKessaiOrderSuryo"`
	KessaiYakuzyouSuryo  string `json:"sKessaiYakuzyouSuryo"`
	KessaiYakuzyouPrice  string `json:"sKessaiYakuzyouPrice"`
	KessaiTateTesuryou   string `json:"sKessaiTateTesuryou"`
	KessaiZyunHibu       string `json:"sKessaiZyunHibu"`
	KessaiGyakuhibu      string `json:"sKessaiGyakuhibu"`
	KessaiKakikaeryou    string `json:"sKessaiKakikaeryou"`
	KessaiKanrihi        string `json:"sKessaiKanrihi"`
	KessaiKasikaburyou   string `json:"sKessaiKasikaburyou"`
	KessaiSonota         string `json:"sKessaiSonota"`
	KessaiSoneki         string `json:"sKessaiSoneki"`
}

// ResYakuzyouSikkou 約定失効リストの要素 (変更なし)
type ResYakuzyouSikkou struct {
	YakuzyouWarningCode string `json:"sYakuzyouWarningCode"`
	YakuzyouWarningText string `json:"sYakuzyouWarningText"`
	YakuzyouSuryou      string `json:"sYakuzyouSuryou"`
	YakuzyouPrice       string `json:"sYakuzyouPrice"`
	YakuzyouDate        string `json:"sYakuzyouDate"`
}
