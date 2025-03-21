// business_functions/res_zan_shinki_kano_ijiritu.go
package business_functions

// ResZanShinkiKanoIjiritu は建余力＆本日維持率のレスポンスを表すDTO
type ResZanShinkiKanoIjiritu struct {
	P_no                    string  `json:"p_no"`                    // p_no
	SCLMID                  string  `json:"sCLMID"`                  // 機能ID, CLMZanShinkiKanoIjiritu
	SResultCode             string  `json:"sResultCode"`             // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText             string  `json:"sResultText"`             // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode            string  `json:"sWarningCode"`            // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText            string  `json:"sWarningText"`            // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode              string  `json:"sIssueCode"`              // 銘柄コード, 要求設定値
	SSizyouC                string  `json:"sSizyouC"`                // 市場, 要求設定値
	SSummaryUpdate          string  `json:"sSummaryUpdate"`          // 更新日時, YYYYMMDDHHMM
	SSummarySinyouSinkidate string  `json:"sSummarySinyouSinkidate"` // 信用新規建可能額
	SItakuhosyoukin         float64 `json:"sItakuhosyoukin,string"`  // 委託保証金率(%), 0.00～9999999999.99
	SOisyouKakuteiFlg       string  `json:"sOisyouKakuteiFlg"`       // 追証フラグ, 0：未確定, 1：確定
}
