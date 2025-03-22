// internal/infrastructure/client/dto/master/response/stock_market_master.go
package response

// StockMarketMaster は、株式銘柄市場マスタの情報を表すDTOです。
type StockMarketMaster struct {
	CLMID                 string `json:"sCLMID"`                    // 機能ID (CLMIssueSizyouMstKabu)
	IssueCode             string `json:"sIssueCode"`                // 銘柄コード
	ListingMarket         string `json:"sZyouzyouSizyou"`           // 上場市場
	SystemCategory        string `json:"sSystemC"`                  // システム区分
	LowerLimit            string `json:"sNehabaMin"`                // 値幅下限
	UpperLimit            string `json:"sNehabaMax"`                // 値幅上限
	IssueClassification   string `json:"sIssueKubunC"`              // 銘柄区分
	PriceLimitMarket      string `json:"sNehabaSizyouC"`            // 値幅市場区分
	MarginEligibility     string `json:"sSinyouC"`                  // 信用区分 (1:貸借銘柄, 2:信用制度銘柄, 3:一般信用銘柄)
	InitialListingDate    string `json:"sSinkiZyouzyouDay"`         // 新規上場日 (YYYYMMDD)
	PriceLimitExpiryDate  string `json:"sNehabaKigenDay"`           // 値幅期限日 (YYYYMMDD)
	PriceLimitRestriction string `json:"sNehabaKiseiC"`             // 値幅規制区分
	PriceLimitValue       string `json:"sNehabaKiseiTi"`            // 値幅規制値
	PriceLimitCheckFlag   string `json:"sNehabaCheckKahiC"`         // 値幅チェック可否区分
	IssueSection          string `json:"sIssueBubetuC"`             // 銘柄部別区分
	PreviousClose         string `json:"sZenzituOwarine"`           // 前日終値
	PriceCalcMarket       string `json:"sNehabaSansyutuSizyouC"`    // 値幅算出市場区分
	IssueRegulation1      string `json:"sIssueKisei1C"`             // 銘柄規制１区分
	IssueRegulation2      string `json:"sIssueKisei2C"`             // 銘柄規制２区分
	ListingCategory       string `json:"sZyouzyouKubun"`            // 上場区分
	DelistingDate         string `json:"sZyouzyouHaisiDay"`         // 上場廃止日 (YYYYMMDD)
	MarketTradingUnit     string `json:"sSizyoubetuBaibaiTani"`     // 売買単位
	NextMarketTradingUnit string `json:"sSizyoubetuBaibaiTaniYoku"` // 売買単位(翌営業日)
	TickUnitNumber        string `json:"sYobineTaniNumber"`         // 呼値の単位番号
	NextTickUnitNumber    string `json:"sYobineTaniNumberYoku"`     // 呼値の単位番号(翌営業日)
	InformationSource     string `json:"sZyouhouSource"`            // 情報系ソース
	InformationCode       string `json:"sZyouhouCode"`              // 情報系コード
	PublicOfferingPrice   string `json:"sKouboPrice"`               // 公募価格
	CreateDate            string `json:"sCreateDate"`               // 作成日時
	UpdateDate            string `json:"sUpdateDate"`               // 更新日時
	UpdateNumber          string `json:"sUpdateNumber"`             // 更新通番
}
