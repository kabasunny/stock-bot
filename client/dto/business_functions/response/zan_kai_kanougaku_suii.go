// business_functions/res_zan_kai_kanougaku_suii.go
package business_functions

// ResZanKaiKanougakuSuii は可能額推移のレスポンスを表すDTO
type ResZanKaiKanougakuSuii struct {
	P_no               string             `json:"p_no"`               // p_no
	SCLMID             string             `json:"sCLMID"`             // 機能ID, CLMZanKaiKanougakuSuii
	SResultCode        string             `json:"sResultCode"`        // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText        string             `json:"sResultText"`        // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode       string             `json:"sWarningCode"`       // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText       string             `json:"sWarningText"`       // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SUpdateDate        string             `json:"sUpdateDate"`        // 更新日時, YYYYMMDDHHMM
	SNearaiKubun       string             `json:"sNearaiKubun"`       // 値洗い区分, 0:値洗い停止, 1:値洗い中, 2:値洗い終了
	AKanougakuSuiiList []ResKanougakuSuii `json:"aKanougakuSuiiList"` // 可能額推移リスト
}

// ResKanougakuSuii 可能額推移リストの要素
type ResKanougakuSuii struct {
	SHituke                      string `json:"sHituke"`                      // 日付, YYYYMMDD
	SAzukariKin                  string `json:"sAzukariKin"`                  // 預り金
	SHattyuZyutoukin             string `json:"sHattyuZyutoukin"`             // 発注済み注文充当金
	SHibakariKousokukin          string `json:"sHibakariKousokukin"`          // 日計り拘束金
	SSonotaKousokukin            string `json:"sSonotaKousokukin"`            // その他拘束金
	SGenkinHosyoukin             string `json:"sGenkinHosyoukin"`             // 現金保証金
	SDaiyouHyoukagaku            string `json:"sDaiyouHyoukagaku"`            // 代用証券評価額
	SSasiireHosyoukin            string `json:"sSasiireHosyoukin"`            // 差入保証金
	SSinyouTateHyoukaSon         string `json:"sSinyouTateHyoukaSon"`         // 信用建株 評価損
	SSinyouTateHyoukaEki         string `json:"sSinyouTateHyoukaEki"`         // 信用建株 評価益
	SSinyouTadeSyoukeihi         string `json:"sSinyouTadeSyoukeihi"`         // 信用建株 諸経費
	SMiukewatasiKessaiSon        string `json:"sMiukewatasiKessaiSon"`        // 信用建株 未受渡決済損
	SMiukewatasiKessaiEki        string `json:"sMiukewatasiKessaiEki"`        // 信用建株 未受渡決済益
	SUkeireHosyoukin             string `json:"sUkeireHosyoukin"`             // 受入保証金
	SMikessaiTateDaikin          string `json:"sMikessaiTateDaikin"`          // 未決済建株代金
	SGenbikiWatasiTateDaikin     string `json:"sGenbikiWatasiTateDaikin"`     // 現引/現渡建株代金
	SHituyouHosyoukin            string `json:"sHituyouHosyoukin"`            // 必要保証金
	SHituyouGenkinHosyoukin      string `json:"sHituyouGenkinHosyoukin"`      // 必要現金保証金
	SHosyoukinYoryoku            string `json:"sHosyoukinYoryoku"`            // 保証金余力
	SGenkinHosyoukinYoryoku      string `json:"sGenkinHosyoukinYoryoku"`      // 現金保証金余力
	SItakuHosyoukinRitu          string `json:"sItakuHosyoukinRitu"`          // 委託保証金率(%)
	SHosyoukinHikidasiKousokukin string `json:"sHosyoukinHikidasiKousokukin"` // 保証金引出拘束金
	SHosyoukinHikidasiYoryoku    string `json:"sHosyoukinHikidasiYoryoku"`    // 保証金引出余力
	SOisyouHituyouHosyoukin      string `json:"sOisyouHituyouHosyoukin"`      // 追証必要保証金
	SOisyouYoryoku               string `json:"sOisyouYoryoku"`               // 追証余力
	SFusokugaku                  string `json:"sFusokugaku"`                  // 追証/立替金/保証金不足額
	SGenbutuKaitukeKanougaku     string `json:"sGenbutuKaitukeKanougaku"`     // 現物株式買付可能額
	SSinyouSinkidateKanougaku    string `json:"sSinyouSinkidateKanougaku"`    // 信用新規建可能額
	SGenbikiKanougaku            string `json:"sGenbikiKanougaku"`            // 信用現引可能額
	STousinKaitukeKanougaku      string `json:"sTousinKaitukeKanougaku"`      // 投信買付可能額
	SSyukkinKanougaku            string `json:"sSyukkinKanougaku"`            // 出金可能額
}
