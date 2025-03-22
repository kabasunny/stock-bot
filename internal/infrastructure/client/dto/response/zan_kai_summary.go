// response/zan_kai_summary.go
package response

// ResZanKaiSummary は可能額サマリーのレスポンスを表すDTO
type ResZanKaiSummary struct {
	P_no                         string                       `json:"p_no"`                         // p_no
	SCLMID                       string                       `json:"sCLMID"`                       // 機能ID, CLMZanKaiSummary
	SResultCode                  string                       `json:"sResultCode"`                  // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                  string                       `json:"sResultText"`                  // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                 string                       `json:"sWarningCode"`                 // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                 string                       `json:"sWarningText"`                 // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SUpdateDate                  string                       `json:"sUpdateDate"`                  // 更新日時, YYYYMMDDHHMM
	SOisyouHasseiFlg             string                       `json:"sOisyouHasseiFlg"`             // 追証発生フラグ, 1:追証発生, 0:追証未発生
	SOhzsKeisanDay               string                       `json:"sOhzsKeisanDay"`               // 追証発生状況詳細.計算日, YYYYMMDD
	SOhzsGenkinHosyoukin         string                       `json:"sOhzsGenkinHosyoukin"`         // 追証発生状況詳細.現金保証金
	SOhzsDaiyouHyoukagaku        string                       `json:"sOhzsDaiyouHyoukagaku"`        // 追証発生状況詳細.代用証券評価額
	SOhzsSasiireHosyoukin        string                       `json:"sOhzsSasiireHosyoukin"`        // 追証発生状況詳細.差入保証金
	SOhzsHyoukaSoneki            string                       `json:"sOhzsHyoukaSoneki"`            // 追証発生状況詳細.評価損益
	SOhzsSyokeihi                string                       `json:"sOhzsSyokeihi"`                // 追証発生状況詳細.諸経費
	SOhzsMiukeKessaiSon          string                       `json:"sOhzsMiukeKessaiSon"`          // 追証発生状況詳細.未受渡決済損
	SOhzsMiukeKessaiEki          string                       `json:"sOhzsMiukeKessaiEki"`          // 追証発生状況詳細.未受渡決済益
	SOhzsUkeireHosyoukin         string                       `json:"sOhzsUkeireHosyoukin"`         // 追証発生状況詳細.受入保証金
	SOhzsTatekabuDaikin          string                       `json:"sOhzsTatekabuDaikin"`          // 追証発生状況詳細.建株代金
	SOhzsItakuHosyoukinRitu      string                       `json:"sOhzsItakuHosyoukinRitu"`      // 追証発生状況詳細.委託保証金率(%)
	STatekaekinHasseiFlg         string                       `json:"sTatekaekinHasseiFlg"`         // 立替金発生フラグ, 1:立替金発生, 0:立替金未発生
	SThzNyukinKigenDay           string                       `json:"sThzNyukinKigenDay"`           // 立替金発生状況.入金期限, YYYYMMDD
	SThzSeisangaku               string                       `json:"sThzSeisangaku"`               // 立替金発生状況.精算額
	SThzHibakariKousokukin       string                       `json:"sThzHibakariKousokukin"`       // 立替金発生状況.日計り拘束金
	SThzHurikaegaku              string                       `json:"sThzHurikaegaku"`              // 立替金発生状況.振替額
	SThzHituyouNyukingaku        string                       `json:"sThzHituyouNyukingaku"`        // 立替金発生状況.必要入金額
	SThzKakuteiFlg               string                       `json:"sThzKakuteiFlg"`               // 立替金発生状況.確定フラグ, 1:入金請求管理.計算日＜営業日, 0:上記外
	SGenbutuKabuKaituke          string                       `json:"sGenbutuKabuKaituke"`          // 株式現物買付可能額
	SSinyouSinkidate             string                       `json:"sSinyouSinkidate"`             // 信用新規建可能額
	SSinyouGenbiki               string                       `json:"sSinyouGenbiki"`               // 信用現引可能額
	SHosyouKinritu               string                       `json:"sHosyouKinritu"`               // 委託保証金率(%)
	SNseityouTousiKanougaku      string                       `json:"sNseityouTousiKanougaku"`      // NISA成長投資可能額
	STousinKaituke               string                       `json:"sTousinKaituke"`               // 投信買付可能額
	SRuitouKaituke               string                       `json:"sRuitouKaituke"`               // MMF・中国F買付
	SIPOKounyu                   string                       `json:"sIPOKounyu"`                   // IPO購入可能額
	SSyukkin                     string                       `json:"sSyukkin"`                     // 出金可能額
	SFusokugaku                  string                       `json:"sFusokugaku"`                  // 不足額(入金請求額）
	SLargeKaidateYoryoku         string                       `json:"sLargeKaidateYoryoku"`         // 先物買建
	SMiniKaidateYoryoku          string                       `json:"sMiniKaidateYoryoku"`          // OPプット売建(ミニ)
	SLargeUridateYoryoku         string                       `json:"sLargeUridateYoryoku"`         // 先物売建
	SMiniUridateYoryoku          string                       `json:"sMiniUridateYoryoku"`          // OPコール売建(ミニ)
	SOpKaidateYoryokyu           string                       `json:"sOpKaidateYoryokyu"`           // オプション新規買建
	SSyoukokinFusokugaku         string                       `json:"sSyoukokinFusokugaku"`         // 証拠金不足額（本日請求額）
	SGenbutuBaibaiDaikin         string                       `json:"sGenbutuBaibaiDaikin"`         // 現物売買代金
	SGenbutuOrderCount           string                       `json:"sGenbutuOrderCount"`           // 現物注文件数
	SGenbutuYakuzyouCount        string                       `json:"sGenbutuYakuzyouCount"`        // 現物約定件数
	SSinyouBaibaiDaikin          string                       `json:"sSinyouBaibaiDaikin"`          // 信用売買代金
	SSinyouOrderCount            string                       `json:"sSinyouOrderCount"`            // 信用注文件数
	SSinyouYakuzyouCount         string                       `json:"sSinyouYakuzyouCount"`         // 信用約定件数
	SSakiBaibaiDaikin            string                       `json:"sSakiBaibaiDaikin"`            // 先物売買代金
	SSakiOrderCount              string                       `json:"sSakiOrderCount"`              // 先物注文件数
	SSakiYakuzyouCount           string                       `json:"sSakiYakuzyouCount"`           // 先物約定件数
	SOpBaibaiDaikin              string                       `json:"sOpBaibaiDaikin"`              // オプション売買代金
	SOpOrderCount                string                       `json:"sOpOrderCount"`                // オプション注文件数
	SOpYakuzyouCount             string                       `json:"sOpYakuzyouCount"`             // オプション約定件数
	AHikazeiKouzaList            []ResHikazeiKouza            `json:"aHikazeiKouzaList"`            // 非課税口座リスト
	AOisyouHasseiZyoukyouList    []ResOisyouHasseiZyoukyou    `json:"aOisyouHasseiZyoukyouList"`    // 追証発生状況リスト
	AHosyoukinSeikyuZyoukyouList []ResHosyoukinSeikyuZyoukyou `json:"aHosyoukinSeikyuZyoukyouList"` // 保証金請求発生状況リスト
}

// ResHikazeiKouza 非課税口座リストの要素
type ResHikazeiKouza struct {
	SHikazeiTekiyouYear    string `json:"sHikazeiTekiyouYear"`    // 適用年（対象年）, YYYY
	SSeityouTousiKanougaku string `json:"sSeityouTousiKanougaku"` // 成長投資可能額
}

// ResOisyouHasseiZyoukyou 追証発生状況リストの要素
type ResOisyouHasseiZyoukyou struct {
	SOhzHasseiDay           string `json:"sOhzHasseiDay"`           // 発生日 YYYYMMDD
	SOhzHosyoukinRitu       string `json:"sOhzHosyoukinRitu"`       // 保証金率(%)
	SOhzNyukinKigenDay      string `json:"sOhzNyukinKigenDay"`      // 入金期限 YYYYMMDDHHMM
	SOhzOisyouKingaku       string `json:"sOhzOisyouKingaku"`       // 追証金額
	SOhzKakuteiFlg          string `json:"sOhzKakuteiFlg"`          // 確定フラグ 1:入金請求管理.計算日＜営業日 0:上記外
	SOhzHosyoukinZougen     string `json:"sOhzHosyoukinZougen"`     // 保証金増減
	SOhzNyukin              string `json:"sOhzNyukin"`              // 入金
	SOhzTategyokuKessai     string `json:"sOhzTategyokuKessai"`     // 建玉決済
	SOhzKessaisonNyukin     string `json:"sOhzKessaisonNyukin"`     // 決済損入金
	SOhzMikaisyouKingaku    string `json:"sOhzMikaisyouKingaku"`    // 未解消金額
	SOhzMikaisyouKingakuFlg string `json:"sOhzMikaisyouKingakuFlg"` // 未解消金額フラグ 未使用
}

// ResHosyoukinSeikyuZyoukyou 保証金請求発生状況リストの要素
type ResHosyoukinSeikyuZyoukyou struct {
	SHshzNyukinKigenDay           string `json:"sHshzNyukinKigenDay"`           // 入金期限 YYYYMMDDHHMM
	SHshzHosyoukinHasseiDay       string `json:"sHshzHosyoukinHasseiDay"`       // 保証金発生日 YYYYMMDDHHMM
	SHshzHosyoukin                string `json:"sHshzHosyoukin"`                // 保証金
	SHshzGenkinHosyoukinHasseiDay string `json:"sHshzGenkinHosyoukinHasseiDay"` // 現金保証金発生日 YYYYMMDDHHMM
	SHshzGenkinHosyoukin          string `json:"sHshzGenkinHosyoukin"`          // 現金保証金
}
