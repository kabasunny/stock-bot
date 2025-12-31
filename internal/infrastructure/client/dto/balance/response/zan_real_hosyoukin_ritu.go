// response/zan_real_hosyoukin_ritu.go
package response

// ResZanRealHosyoukinRitu はリアル保証金率のレスポンスを表すDTO
type ResZanRealHosyoukinRitu struct {
	P_no                      string `json:"p_no"`                      // p_no
	SCLMID                    string `json:"sCLMID"`                    // 機能ID, CLMZanRealHosyoukinRitu
	SResultCode               string `json:"sResultCode"`               // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText               string `json:"sResultText"`               // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode              string `json:"sWarningCode"`              // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText              string `json:"sWarningText"`              // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SSasiireHosyoukin         string `json:"sSasiireHosyoukin"`         // 差入保証金
	SHyoukaSonEki             string `json:"sHyoukaSonEki"`             // 評価損益
	SUkeireHosyoukin          string `json:"sUkeireHosyoukin"`          // 受入保証金
	STateKabuDaikin           string `json:"sTateKabuDaikin"`           // 建株代金
	SItakuHosyoukinRitu       string `json:"sItakuHosyoukinRitu"`       // 委託保証金率(%)
	SOisyouHituyouHosyoukin   string `json:"sOisyouHituyouHosyoukin"`   // 追証必要保証金
	SOisyouYoryoku            string `json:"sOisyouYoryoku"`            // 追証余力
	ST0SasiireHosyoukin       string `json:"sT0SasiireHosyoukin"`       // 差入保証金
	ST0HyoukaSonEki           string `json:"sT0HyoukaSonEki"`           // 評価損益
	ST0UkeireHosyoukin        string `json:"sT0UkeireHosyoukin"`        // 受入保証金
	ST0TateKabuDaikin         string `json:"sT0TateKabuDaikin"`         // 建株代金
	ST0ItakuHosyoukinRitu     string `json:"sT0ItakuHosyoukinRitu"`     // 委託保証金率(%)
	ST0OisyouHituyouHosyoukin string `json:"sT0OisyouHituyouHosyoukin"` // 追証必要保証金
	ST0OisyouYoryoku          string `json:"sT0OisyouYoryoku"`          // 追証余力
	ST5SasiireHosyoukin       string `json:"sT5SasiireHosyoukin"`       // 差入保証金
	ST5HyoukaSonEki           string `json:"sT5HyoukaSonEki"`           // 評価損益
	ST5UkeireHosyoukin        string `json:"sT5UkeireHosyoukin"`        // 受入保証金
	ST5TateKabuDaikin         string `json:"sT5TateKabuDaikin"`         // 建株代金
	ST5ItakuHosyoukinRitu     string `json:"sT5ItakuHosyoukinRitu"`     // 委託保証金率(%)
	ST5OisyouHituyouHosyoukin string `json:"sT5OisyouHituyouHosyoukin"` // 追証必要保証金
	ST5OisyouYoryoku          string `json:"sT5OisyouYoryoku"`          // 追証余力
}
