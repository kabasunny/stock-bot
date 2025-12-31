// internal/infrastructure/client/dto/master/response/margin_rate.go
package response

// ResMarginRate は、代用掛目の情報を表すDTO
type ResMarginRate struct {
	CLMID                string `json:"sCLMID"`               // 機能ID (CLMDaiyouKakeme)
	SystemAccountType    string `json:"sSystemKouzaKubun"`    // システム口座区分
	IssueCode            string `json:"sIssueCode"`           // 銘柄コード
	ApplicableDate       string `json:"sTekiyouDay"`          // 適用日 (YYYYMMDD)
	MarginCollateralRate string `json:"sHosyokinDaiyoKakeme"` // 保証金代用掛目
	DeletionDate         string `json:"sDeleteDay"`           // 削除日
	CreateDate           string `json:"sCreateDate"`          // 作成日時
	UpdateNumber         string `json:"sUpdateNumber"`        // 更新通番
	UpdateDate           string `json:"sUpdateDate"`          // 更新日時
}
