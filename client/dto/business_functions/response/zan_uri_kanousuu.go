// business_functions/res_zan_uri_kanousuu.go
package business_functions

// ResZanUriKanousuu は売却可能数量のレスポンスを表すDTO
type ResZanUriKanousuu struct {
	P_no                           string `json:"p_no"`                           // p_no
	SCLMID                         string `json:"sCLMID"`                         // 機能ID, CLMZanUriKanousuu
	SResultCode                    string `json:"sResultCode"`                    // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                    string `json:"sResultText"`                    // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                   string `json:"sWarningCode"`                   // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                   string `json:"sWarningText"`                   // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode                     string `json:"sIssueCode"`                     // 銘柄コード, 要求設定値
	SSummaryUpdate                 string `json:"sSummaryUpdate"`                 // 更新日時, YYYYMMDDHHMM
	SZanKabuSuryouUriKanouIppan    string `json:"sZanKabuSuryouUriKanouIppan"`    // 売付可能株数(一般)
	SZanKabuSuryouUriKanouTokutei  string `json:"sZanKabuSuryouUriKanouTokutei"`  // 売付可能株数(特定)
	SZanKabuSuryouUriKanouNisa     string `json:"sZanKabuSuryouUriKanouNisa"`     // 売付可能株数(NISA)
	SZanKabuSuryouUriKanouNseityou string `json:"sZanKabuSuryouUriKanouNseityou"` // 売付可能株数(N成長)
}
