// internal/infrastructure/client/dto/master/response/margin_master.go
package response

// ResMarginMaster は、保証金マスタの情報を表すDTOです。
type ResMarginMaster struct {
	CLMID             string `json:"sCLMID"`              // 機能ID (CLMHosyoukinMst)
	SystemAccountType string `json:"sSystemKouzaKubun"`   // システム口座区分
	IssueCode         string `json:"sIssueCode"`          // 銘柄コード
	ListingMarket     string `json:"sZyouzyouSizyou"`     // 上場市場
	ChangeDate        string `json:"sHenkouDay"`          // 変更日 (YYYYMMDD)
	CollateralRate    string `json:"sDaiyoHosyokinRitu"`  // 代用保証金率(%)
	CashMarginRate    string `json:"sGenkinHosyokinRitu"` // 現金保証金率(%)
	CreateDate        string `json:"sCreateDate"`         // 作成日時
	UpdateNumber      string `json:"sUpdateNumber"`       // 更新通番
	UpdateDate        string `json:"sUpdateDate"`         // 更新日時
}
