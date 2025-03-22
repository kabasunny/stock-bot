// internal/infrastructure/client/dto/master/response/stock_master.go
package response

// StockMaster は、株式銘柄マスタの情報を表すDTOです。
type StockMaster struct {
	CLMID                     string `json:"sCLMID"`                 // 機能ID (CLMIssueMstKabu)
	IssueCode                 string `json:"sIssueCode"`             // 銘柄コード
	IssueName                 string `json:"sIssueName"`             // 銘柄名称
	IssueNameShort            string `json:"sIssueNameRyaku"`        // 銘柄名略称
	IssueNameKana             string `json:"sIssueNameKana"`         // 銘柄名（カナ）
	IssueNameEnglish          string `json:"sIssueNameEizi"`         // 銘柄名（英語表記）
	SpecialAccountEligible    string `json:"sTokuteiF"`              // 特定口座対象区分 (1：特定参加, 0：不参加)
	TaxExemptEligible         string `json:"sHikazeiC"`              // 非課税対象区分
	ListedSharesOutstanding   string `json:"sZyouzyouHakkouKabusu"`  // 上場発行株数
	RightsAllotmentFlag       string `json:"sKenriotiFlag"`          // 権利落ちフラグ
	ExRightsDate              string `json:"sKenritukiSaisyuDay"`    // 権利付最終日 (YYYYMMDD)
	ListingStatus             string `json:"sZyouzyouNyusatuC"`      // 上場・入札区分
	TenderOfferReleaseDate    string `json:"sNyusatuKaizyoDay"`      // 入札解除日 (YYYYMMDD)
	TenderOfferDate           string `json:"sNyusatuDay"`            // 入札日 (YYYYMMDD)
	TradingUnit               string `json:"sBaibaiTani"`            // 売買単位
	NextTradingUnit           string `json:"sBaibaiTaniYoku"`        // 売買単位(翌営業日)
	TradingHaltFlag           string `json:"sBaibaiTeisiC"`          // 売買停止区分
	IssueStartDate            string `json:"sHakkouKaisiDay"`        // 発行開始日 (YYYYMMDD)
	IssueEndDate              string `json:"sHakkouSaisyuDay"`       // 発行最終日 (YYYYMMDD)
	FiscalMonth               string `json:"sKessanC"`               // 決算月
	FiscalDate                string `json:"sKessanDay"`             // 決算日 (YYYYMMDD)
	ListedDate                string `json:"sZyouzyouOutouDay"`      // 上場応答日 (YYYYMMDD)
	Category2ExpirationFlag   string `json:"sNiruiKizituC"`          // 二類期日区分
	LargeLotQuantity          string `json:"sOogutiKabusu"`          // 大口株数
	LargeLotAmount            string `json:"sOogutiKingaku"`         // 大口金額
	FloorSlipOutputFlag       string `json:"sBadenpyouOutputYNC"`    // 場伝票出力有無区分
	MarginCollateralRate      string `json:"sHosyoukinDaiyouKakeme"` // 保証金代用掛目
	CollateralValuationPrice  string `json:"sDaiyouHyoukaTanka"`     // 代用証券評価単価
	InstitutionParticipation  string `json:"sKikoSankaC"`            // 機構参加区分
	ProvisionalSettlementFlag string `json:"sKarikessaiC"`           // 仮決済区分
	PreferredMarket           string `json:"sYusenSizyou"`           // 優先市場
	UnlimitedTermEligible     string `json:"sMukigenC"`              // 無期限対象区分
	IndustryCode              string `json:"sGyousyuCode"`           // 業種コード
	IndustryName              string `json:"sGyousyuName"`           // 業種コード名
	SOR                       string `json:"sSorC"`                  // ＳＯＲ対象銘柄区分
	CreateDate                string `json:"sCreateDate"`            // 作成日時
	UpdateDate                string `json:"sUpdateDate"`            // 更新日時
	UpdateNumber              string `json:"sUpdateNumber"`          // 更新通番
}
