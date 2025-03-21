// business_functions/res_zan_kai_genbutu_kaituke_syousai.go
package business_functions

// ResZanKaiGenbutuKaitukeSyousai は現物株式買付可能額詳細のレスポンスを表すDTO
type ResZanKaiGenbutuKaitukeSyousai struct {
	P_no                            string `json:"p_no"`                            // p_no
	SCLMID                          string `json:"sCLMID"`                          // 機能ID, CLMZanKaiGenbutuKaitukeSyousai
	SResultCode                     string `json:"sResultCode"`                     // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                     string `json:"sResultText"`                     // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                    string `json:"sWarningCode"`                    // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                    string `json:"sWarningText"`                    // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SHitukeIndex                    string `json:"sHitukeIndex"`                    // 日付インデックス, 要求設定値
	SHituke                         string `json:"sHituke"`                         // 指定日（日付）, YYYYMMDD
	SGenkinHosyoukin                string `json:"sGenkinHosyoukin"`                // 現金保証金
	SHosyoukinGenbutuKaitukeKanouga string `json:"sHosyoukinGenbutuKaitukeKanouga"` // 保証金からの現物株式買付可能額
	SGenbutuKaitukeKanougaku        string `json:"sGenbutuKaitukeKanougaku"`        // 現物株式買付可能額
	SAzukariKin                     string `json:"sAzukariKin"`                     // 預り金
	SHattyuZyutoukin                string `json:"sHattyuZyutoukin"`                // 発注済み注文充当金
	SHibakariKousokukin             string `json:"sHibakariKousokukin"`             // 日計り拘束金
	SSonotaKousokukin               string `json:"sSonotaKousokukin"`               // その他拘束金
	SHituyouGenkinHosyoukin         string `json:"sHituyouGenkinHosyoukin"`         // 必要現金保証金
	SDaiyouHyoukagaku               string `json:"sDaiyouHyoukagaku"`               // 代用証券評価額
	STatekabuHyoukaSoneki           string `json:"sTatekabuHyoukaSoneki"`           // 建株評価損益
	STatekabuSyoukeihi              string `json:"sTatekabuSyoukeihi"`              // 建株諸経費
	SMiukewatasiKessaiSon           string `json:"sMiukewatasiKessaiSon"`           // 未受渡建株決済損
	SMiukewatasiKessaiEki           string `json:"sMiukewatasiKessaiEki"`           // 未受渡建株決済益
	SUkeireHosyoukin                string `json:"sUkeireHosyoukin"`                // 受入保証金
	SHituyouHosyoukin               string `json:"sHituyouHosyoukin"`               // 必要保証金
	SHosyoukinYoryoku               string `json:"sHosyoukinYoryoku"`               // 保証金余力
}
