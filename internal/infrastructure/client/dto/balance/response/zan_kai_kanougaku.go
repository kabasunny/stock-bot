// response/zan_kai_kanougaku.go
package response

// ResZanKaiKanougaku　は買余力のレスポンスを表すDTO
type ResZanKaiKanougaku struct {
	P_no                           string `json:"p_no"`                           // p_no
	SCLMID                         string `json:"sCLMID"`                         // 機能ID, CLMZanKaiKanougaku
	SResultCode                    string `json:"sResultCode"`                    // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                    string `json:"sResultText"`                    // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                   string `json:"sWarningCode"`                   // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                   string `json:"sWarningText"`                   // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode                     string `json:"sIssueCode"`                     // 銘柄コード, 要求設定値
	SSizyouC                       string `json:"sSizyouC"`                       // 市場, 要求設定値
	SSummaryUpdate                 string `json:"sSummaryUpdate"`                 // 更新日時, YYYYMMDDHHMM
	SSummaryGenkabuKaituke         string `json:"sSummaryGenkabuKaituke"`         // 株式現物買付可能額
	SSummaryNseityouTousiKanougaku string `json:"sSummaryNseityouTousiKanougaku"` // NISA成長投資可能額
	SHusokukinHasseiFlg            string `json:"sHusokukinHasseiFlg"`            // 不足金発生フラグ, 0：未発生, 1：発生
}
