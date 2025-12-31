// response/zan_kai_sinyou_sinkidate_syousai.go
package response

// ResZanKaiSinyouSinkidateSyousai は信用新規建て可能額詳細のレスポンスを表すDTO
type ResZanKaiSinyouSinkidateSyousai struct {
	P_no                        string `json:"p_no"`                        // p_no
	SCLMID                      string `json:"sCLMID"`                      // 機能ID, CLMZanKaiSinyouSinkidateSyousai
	SResultCode                 string `json:"sResultCode"`                 // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                 string `json:"sResultText"`                 // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                string `json:"sWarningCode"`                // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                string `json:"sWarningText"`                // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SHitukeIndex                string `json:"sHitukeIndex"`                // 日付インデックス, 要求設定値
	SHituke                     string `json:"sHituke"`                     // 指定日（日付）, YYYYMMDD
	SUkeireHosyoukin            string `json:"sUkeireHosyoukin"`            // 受入保証金
	SHituyouHosyoukin           string `json:"sHituyouHosyoukin"`           // 必要保証金
	SHosyoukinYoryoku           string `json:"sHosyoukinYoryoku"`           // 保証金余力
	SHosyoukinTyousyuRitu       string `json:"sHosyoukinTyousyuRitu"`       // 保証金徴収率(%)
	SSinyouSinkidateKanougaku   string `json:"sSinyouSinkidateKanougaku"`   // 信用新規建可能額
	SAzukariKin                 string `json:"sAzukariKin"`                 // 預り金
	SHattyuZyutoukin            string `json:"sHattyuZyutoukin"`            // 発注済み注文充当金
	SSonotaKousokukin           string `json:"sSonotaKousokukin"`           // その他拘束金
	SGenkinHosyoukin            string `json:"sGenkinHosyoukin"`            // 現金保証金
	SDaiyouHyoukagaku           string `json:"sDaiyouHyoukagaku"`           // 代用証券評価額
	SHattyuDaiyouHyoukagaku     string `json:"sHattyuDaiyouHyoukagaku"`     // 現物買発注分代用証券評価額
	SSasiireHosyoukin           string `json:"sSasiireHosyoukin"`           // 差入保証金
	SSinkiTesuryou              string `json:"sSinkiTesuryou"`              // 新規建手数料
	SHibuGyakuhibuKousokukin    string `json:"sHibuGyakuhibuKousokukin"`    // 日歩・逆日歩・貸株料拘束金
	SHibuGyakuhibuSyueki        string `json:"sHibuGyakuhibuSyueki"`        // 日歩・逆日歩収益
	SSonotaTateSyokeihi         string `json:"sSonotaTateSyokeihi"`         // その他未収費用
	SSinyouTadeSyoukeihi        string `json:"sSinyouTadeSyoukeihi"`        // 建株諸経費
	SSinyouTateHyoukaSon        string `json:"sSinyouTateHyoukaSon"`        // 建株評価損
	SSinyouTateHyoukaEki        string `json:"sSinyouTateHyoukaEki"`        // 建株評価益
	STatekabuHyoukaSoneki       string `json:"sTatekabuHyoukaSoneki"`       // 建株評価損益
	SMiukewatasiKessaiSon       string `json:"sMiukewatasiKessaiSon"`       // 未受渡建株決済損
	SMiukewatasiKessaiEki       string `json:"sMiukewatasiKessaiEki"`       // 未受渡建株決済益
	SSaiteiHituyouHosyoukin     string `json:"sSaiteiHituyouHosyoukin"`     // 最低必要保証金
	SHosyoukin                  string `json:"sHosyoukin"`                  // 建株必要保証金
	SHattyuHosyoukin            string `json:"sHattyuHosyoukin"`            // 発注分必要保証金
	SGenbikiWatasiHosyoukin     string `json:"sGenbikiWatasiHosyoukin"`     // 現引/現渡必要保証金
	SMikessaiGenkinHosyoukin    string `json:"sMikessaiGenkinHosyoukin"`    // 建株必要保証金（現金）
	SHattyuGenkinHosyoukin      string `json:"sHattyuGenkinHosyoukin"`      // 発注分必要保証金（現金）
	SGenbwGenkinHosyoukin       string `json:"sGenbwGenkinHosyoukin"`       // 現引/現渡必要保証金（現金）
	SHituyouGenkinHosyoukin     string `json:"sHituyouGenkinHosyoukin"`     // 必要保証金（現金）
	SHosyoukinRitu              string `json:"sHosyoukinRitu"`              // 保証金率(%)
	SHosyoukinIziRitu           string `json:"sHosyoukinIziRitu"`           // 保証金維持率(%)
	SHosyoukinRituIziYoryoku    string `json:"sHosyoukinRituIziYoryoku"`    // 保証金率・維持余力
	SHosyoukinIzirituIziYoryoku string `json:"sHosyoukinIzirituIziYoryoku"` // 保証金維持率・維持余力
	SMikessaiTateDaikin         string `json:"sMikessaiTateDaikin"`         // 建株代金
	SHattyuTateDaikin           string `json:"sHattyuTateDaikin"`           // 発注分建株代金
	SGenbikiWatasiTateDaikin    string `json:"sGenbikiWatasiTateDaikin"`    // 現引/現渡建株代金
	SItakuHosyoukinRitu         string `json:"sItakuHosyoukinRitu"`         // 委託保証金率(%)
	STouzituKessaiSon           string `json:"sTouzituKessaiSon"`           // 本日決済損
	STouzituKessaiEki           string `json:"sTouzituKessaiEki"`           // 本日決済益
	SKessaiTotalToday           string `json:"sKessaiTotalToday"`           // 本日決済損益合計
	STouzituKessaiZenHyouka     string `json:"sTouzituKessaiZenHyouka"`     // 本日決済建株の前日価格評価
	SUkewatasiTategyokuSon      string `json:"sUkewatasiTategyokuSon"`      // 指定日決済損
	SUkewatasiTategyokuEki      string `json:"sUkewatasiTategyokuEki"`      // 指定日決済益
	SKessaiTotalSiteibi         string `json:"sKessaiTotalSiteibi"`         // 指定日決済損益累計
}
