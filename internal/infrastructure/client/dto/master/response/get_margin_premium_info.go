package response

// GetMarginPremiumInfoResponse 逆日歩情報問合取得　レスポンス
// internal/infrastructure/client/dto/master/response/get_margin_premium_info.go

type ResGetMarginPremiumInfo struct {
	CLMID           string                         `json:"sCLMID"`                     // 機能ID
	CLMMfdsHibuInfo []ResMarginPremiumInfoListItem `json:"aCLMMfdsHibuInfo,omitempty"` // 取得リスト
}
type ResMarginPremiumInfoListItem struct {
	IssueCode string `json:"sIssueCode"`      // 対象銘柄コード
	PBWRQ     string `json:"pBWRQ,omitempty"` // 逆日歩
}
