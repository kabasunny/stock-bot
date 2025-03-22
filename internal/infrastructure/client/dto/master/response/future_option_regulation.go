// internal/infrastructure/client/dto/master/response/future_option_regulation.go
package response

// FutureOptionRegulation は、派生銘柄別・市場別規制の情報を表すDTOです。
type FutureOptionRegulation struct {
	CLMID                 string `json:"sCLMID"`            // 機能ID (CLMIssueSizyouKiseiHasei)
	SystemAccountType     string `json:"sSystemKouzaKubun"` // システム口座区分
	IssueCode             string `json:"sIssueCode"`        // 銘柄コード
	ListingMarket         string `json:"sZyouzyouSizyou"`   // 上場市場
	TradingHaltCategory   string `json:"sTeisiKubun"`       // 停止区分
	BuyNew                string `json:"sKaitate"`          // 買建
	NextDayBuyNew         string `json:"sKaitateYoku"`      // 買建（翌営業日）
	SellNew               string `json:"sUritate"`          // 売建
	NextDaySellNew        string `json:"sUritateYoku"`      // 売建（翌営業日）
	BuyRedemption         string `json:"sKaiHensai"`        // 買返済
	NextDayBuyRedemption  string `json:"sKaiHensaiYoku"`    // 買返済（翌営業日）
	SellRedemption        string `json:"sUriHensai"`        // 売返済
	NextDaySellRedemption string `json:"sUriHensaiYoku"`    // 売返済（翌営業日）
	CreateDate            string `json:"sCreateDate"`       // 作成日時
	UpdateDate            string `json:"sUpdateDate"`       // 更新日時
	UpdateNumber          string `json:"sUpdateNumber"`     // 更新通番
}
