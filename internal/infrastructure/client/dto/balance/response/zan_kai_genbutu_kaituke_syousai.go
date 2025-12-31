// response/zan_kai_genbutu_kaituke_syousai.go
package response

// ResZanKaiGenbutuKaitukeSyousai は現物株式買付可能額詳細のレスポンスを表すDTO
type ResZanKaiGenbutuKaitukeSyousai struct {
	P_no                           string `json:"p_no"`                            // p_no
	CLMID                          string `json:"sCLMID"`                          // 機能ID, CLMZanKaiGenbutuKaitukeSyousai
	ResultCode                     string `json:"sResultCode"`                     // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText                     string `json:"sResultText"`                     // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode                    string `json:"sWarningCode"`                    // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText                    string `json:"sWarningText"`                    // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	HitukeIndex                    string `json:"sHitukeIndex"`                    // 日付インデックス, 要求設定値
	Hituke                         string `json:"sHituke"`                         // 指定日（日付）, YYYYMMDD
	GenkinHosyoukin                string `json:"sGenkinHosyoukin"`                // 現金保証金
	HosyoukinGenbutuKaitukeKanouga string `json:"sHosyoukinGenbutuKaitukeKanouga"` // 保証金からの現物株式買付可能額
	GenbutuKaitukeKanougaku        string `json:"sGenbutuKaitukeKanougaku"`        // 現物株式買付可能額
	AzukariKin                     string `json:"sAzukariKin"`                     // 預り金
	HattyuZyutoukin                string `json:"sHattyuZyutoukin"`                // 発注済み注文充当金
	HibakariKousokukin             string `json:"sHibakariKousokukin"`             // 日計り拘束金
	SonotaKousokukin               string `json:"sSonotaKousokukin"`               // その他拘束金
	HituyouGenkinHosyoukin         string `json:"sHituyouGenkinHosyoukin"`         // 必要現金保証金
	DaiyouHyoukagaku               string `json:"sDaiyouHyoukagaku"`               // 代用証券評価額
	TatekabuHyoukaSoneki           string `json:"sTatekabuHyoukaSoneki"`           // 建株評価損益
	TatekabuSyoukeihi              string `json:"sTatekabuSyoukeihi"`              // 建株諸経費
	MiukewatasiKessaiSon           string `json:"sMiukewatasiKessaiSon"`           // 未受渡建株決済損
	MiukewatasiKessaiEki           string `json:"sMiukewatasiKessaiEki"`           // 未受渡建株決済益
	UkeireHosyoukin                string `json:"sUkeireHosyoukin"`                // 受入保証金
	HituyouHosyoukin               string `json:"sHituyouHosyoukin"`               // 必要保証金
	HosyoukinYoryoku               string `json:"sHosyoukinYoryoku"`               // 保証金余力
}
