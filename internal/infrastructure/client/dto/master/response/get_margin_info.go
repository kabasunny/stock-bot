package response

// GetMarginInfoResponse は、証金残情報問合取得のレスポンスを表すDTO
// internal/infrastructure/client/dto/master/response/get_margin_info.go

type ResGetMarginInfo struct {
	CLMID             string                  `json:"sCLMID"`                       // 機能ID
	CLMMfdsSyoukinZan []ResMarginInfoListItem `json:"aCLMMfdsSyoukinZan,omitempty"` // 取得リスト
}
type ResMarginInfoListItem struct {
	IssueCode string `json:"sIssueCode"`      // 対象銘柄コード
	PSFC6     string `json:"pSFC6,omitempty"` // 証金差引残前日比
	PSFD      string `json:"pSFD,omitempty"`  // 証金更新日 YYYY/MM/DD
	PSFD6     string `json:"pSFD6,omitempty"` // 証金回転日数
	PSFF6     string `json:"pSFF6,omitempty"` // 証金融資残
	PSFG6     string `json:"pSFG6,omitempty"` // 証金融資前日比
	PSFKS     string `json:"pSFKS,omitempty"` // 速報確報ステータス 1:速報, 2:確報
	PSFL6     string `json:"pSFL6,omitempty"` // 証金融資・新規
	PSFN6     string `json:"pSFN6,omitempty"` // 証金差引残
	PSFP6     string `json:"pSFP6,omitempty"` // 証金融資・返済
	PSFR6     string `json:"pSFR6,omitempty"` // 貸借倍率
	PSFS6     string `json:"pSFS6,omitempty"` // 証金貸株残
	PSSG6     string `json:"pSSG6,omitempty"` // 証金貸株前日比
	PSSL6     string `json:"pSSL6,omitempty"` // 証金貸株・新規
	PSSP6     string `json:"pSSP6,omitempty"` // 証金貸株・返済
}
