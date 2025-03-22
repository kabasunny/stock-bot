// internal/infrastructure/client/dto/master/response/stock_issue_regulation.go
package response

// StockIssueRegulation は、株式銘柄別・市場別規制の情報を表すDTOです。
type StockIssueRegulation struct {
	CLMID                                      string `json:"sCLMID"`                       // 機能ID (CLMIssueSizyouKiseiKabu)
	SystemAccountType                          string `json:"sSystemKouzaKubun"`            // システム口座区分 (102：ｅ支店)
	IssueCode                                  string `json:"sIssueCode"`                   // 銘柄コード
	ListingMarket                              string `json:"sZyouzyouSizyou"`              // 上場市場 (00：東証)
	TradingHaltCategory                        string `json:"sTeisiKubun"`                  // 停止区分 (0:通常（無）, 1:取引禁止, 2:成行禁止, 3:端株禁止)
	CashBuy                                    string `json:"sGenbutuKaituke"`              // 現物/買付 (sTeisiKubun 参照)
	NextDayCashBuy                             string `json:"sGenbutuKaitukeYoku"`          // 現物/買付（翌営業日）(sTeisiKubun 参照)
	CashSell                                   string `json:"sGenbutuUrituke"`              // 現物/売付 (sTeisiKubun 参照)
	NextDayCashSell                            string `json:"sGenbutuUritukeYoku"`          // 現物/売付（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginNewBuy                  string `json:"sSeidoSinyouSinkiKaitate"`     // 制度信用/買建 (sTeisiKubun 参照)
	NextDayInstitutionalMarginNewBuy           string `json:"sSeidoSinyouSinkiKaitateYoku"` // 制度信用/買建（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginNewSell                 string `json:"sSeidoSinyouSinkiUritate"`     // 制度信用/売建 (sTeisiKubun 参照)
	NextDayInstitutionalMarginNewSell          string `json:"sSeidoSinyouSinkiUritateYoku"` // 制度信用/売建（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginBuyRedemption           string `json:"sSeidoSinyouKaiHensai"`        // 制度信用/買返済 (sTeisiKubun 参照)
	NextDayInstitutionalMarginBuyRedemption    string `json:"sSeidoSinyouKaiHensaiYoku"`    // 制度信用/買返済（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginSellRedemption          string `json:"sSeidoSinyouUriHensai"`        // 制度信用/売返済 (sTeisiKubun 参照)
	NextDayInstitutionalMarginSellRedemption   string `json:"sSeidoSinyouUriHensaiYoku"`    // 制度信用/売返済（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginPhysicalDelivery        string `json:"sSeidoSinyouGenbiki"`          // 制度信用/現引 (sTeisiKubun 参照)
	NextDayInstitutionalMarginPhysicalDelivery string `json:"sSeidoSinyouGenbikiYoku"`      // 制度信用/現引（翌営業日）(sTeisiKubun 参照)
	InstitutionalMarginPhysicalReceipt         string `json:"sSeidoSinyouGenwatasi"`        // 制度信用/現渡 (sTeisiKubun 参照)
	NextDayInstitutionalMarginPhysicalReceipt  string `json:"sSeidoSinyouGenwatasiYoku"`    // 制度信用/現渡（翌営業日）(sTeisiKubun 参照)
	GeneralMarginNewBuy                        string `json:"sIppanSinyouSinkiKaitate"`     // 一般信用/買建 (sTeisiKubun 参照)
	NextDayGeneralMarginNewBuy                 string `json:"sIppanSinyouSinkiKaitateYoku"` // 一般信用/買建（翌営業日）(sTeisiKubun 参照)
	GeneralMarginNewSell                       string `json:"sIppanSinyouSinkiUritate"`     // 一般信用/売建 (sTeisiKubun 参照)
	NextDayGeneralMarginNewSell                string `json:"sIppanSinyouSinkiUritateYoku"` // 一般信用/売建（翌営業日）(sTeisiKubun 参照)
	GeneralMarginBuyRedemption                 string `json:"sIppanSinyouKaiHensai"`        // 一般信用/買返済 (sTeisiKubun 参照)
	NextDayGeneralMarginBuyRedemption          string `json:"sIppanSinyouKaiHensaiYoku"`    // 一般信用/買返済（翌営業日）(sTeisiKubun 参照)
	GeneralMarginSellRedemption                string `json:"sIppanSinyouUriHensai"`        // 一般信用/売返済 (sTeisiKubun 参照)
	NextDayGeneralMarginSellRedemption         string `json:"sIppanSinyouUriHensaiYoku"`    // 一般信用/売返済（翌営業日）(sTeisiKubun 参照)
	GeneralMarginPhysicalDelivery              string `json:"sIppanSinyouGenbiki"`          // 一般信用/現引 (sTeisiKubun 参照)
	NextDayGeneralMarginPhysicalDelivery       string `json:"sIppanSinyouGenbikiYoku"`      // 一般信用/現引（翌営業日）(sTeisiKubun 参照)
	GeneralMarginPhysicalReceipt               string `json:"sIppanSinyouGenwatasi"`        // 一般信用/現渡 (sTeisiKubun 参照)
	NextDayGeneralMarginPhysicalReceipt        string `json:"sIppanSinyouGenwatasiYoku"`    // 一般信用/現渡（翌営業日）(sTeisiKubun 参照)
	PreAdjustmentFlag                          string `json:"sZizenCyouseiC"`               // 事前調整有無 (0：なし（無効）, 1：あり（有効）)
	NextDayPreAdjustmentFlag                   string `json:"sZizenCyouseiCYoku"`           // 事前調整有無（翌営業日）(sZizenCyouseiC 参照)
	SameDaySettlementRegulation                string `json:"sSokuzituNyukinC"`             // 即日入金規制有無 (sZizenCyouseiC 参照)
	NextDaySameDaySettlementRegulation         string `json:"sSokuzituNyukinCYoku"`         // 即日入金規制有無（翌営業日）(sZizenCyouseiC 参照)
	SameDaySettlementRegulationDate            string `json:"sSokuzituNyukinKiseiDate"`     // 即日入金規制日時 (YYYYMMDDHHMMSS)
	ConcentratedMarginCategory                 string `json:"sSinyouSyutyuKubun"`           // 信用一極集中区分 (0：なし, 1：あり, 2：日々公表銘柄)
	NextDayConcentratedMarginCategory          string `json:"sSinyouSyutyuKubunYoku"`       // 信用一極集中区分（翌営業日）(sSinyouSyutyuKubun 参照)
	CreateDate                                 string `json:"sCreateDate"`                  // 作成日時
	UpdateDate                                 string `json:"sUpdateDate"`                  // 更新日時
	UpdateNumber                               string `json:"sUpdateNumber"`                // 更新通番
}
