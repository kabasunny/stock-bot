// response/zan_kai_summary.go
package response

// ResZanKaiSummary は可能額サマリーのレスポンスを表すDTO
type ResZanKaiSummary struct {
	P_no                        string                       `json:"p_no"`                         // p_no
	CLMID                       string                       `json:"sCLMID"`                       // 機能ID, CLMZanKaiSummary
	ResultCode                  string                       `json:"sResultCode"`                  // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText                  string                       `json:"sResultText"`                  // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode                 string                       `json:"sWarningCode"`                 // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText                 string                       `json:"sWarningText"`                 // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	UpdateDate                  string                       `json:"sUpdateDate"`                  // 更新日時, YYYYMMDDHHMM
	OisyouHasseiFlg             string                       `json:"sOisyouHasseiFlg"`             // 追証発生フラグ, 1:追証発生, 0:追証未発生
	OhzsKeisanDay               string                       `json:"sOhzsKeisanDay"`               // 追証発生状況詳細.計算日, YYYYMMDD
	OhzsGenkinHosyoukin         string                       `json:"sOhzsGenkinHosyoukin"`         // 追証発生状況詳細.現金保証金
	OhzsDaiyouHyoukagaku        string                       `json:"sOhzsDaiyouHyoukagaku"`        // 追証発生状況詳細.代用証券評価額
	OhzsSasiireHosyoukin        string                       `json:"sOhzsSasiireHosyoukin"`        // 追証発生状況詳細.差入保証金
	OhzsHyoukaSoneki            string                       `json:"sOhzsHyoukaSoneki"`            // 追証発生状況詳細.評価損益
	OhzsSyokeihi                string                       `json:"sOhzsSyokeihi"`                // 追証発生状況詳細.諸経費
	OhzsMiukeKessaiSon          string                       `json:"sOhzsMiukeKessaiSon"`          // 追証発生状況詳細.未受渡決済損
	OhzsMiukeKessaiEki          string                       `json:"sOhzsMiukeKessaiEki"`          // 追証発生状況詳細.未受渡決済益
	OhzsUkeireHosyoukin         string                       `json:"sOhzsUkeireHosyoukin"`         // 追証発生状況詳細.受入保証金
	OhzsTatekabuDaikin          string                       `json:"sOhzsTatekabuDaikin"`          // 追証発生状況詳細.建株代金
	OhzsItakuHosyoukinRitu      string                       `json:"sOhzsItakuHosyoukinRitu"`      // 追証発生状況詳細.委託保証金率(%)
	TatekaekinHasseiFlg         string                       `json:"sTatekaekinHasseiFlg"`         // 立替金発生フラグ, 1:立替金発生, 0:立替金未発生
	ThzNyukinKigenDay           string                       `json:"sThzNyukinKigenDay"`           // 立替金発生状況.入金期限, YYYYMMDD
	ThzSeisangaku               string                       `json:"sThzSeisangaku"`               // 立替金発生状況.精算額
	ThzHibakariKousokukin       string                       `json:"sThzHibakariKousokukin"`       // 立替金発生状況.日計り拘束金
	ThzHurikaegaku              string                       `json:"sThzHurikaegaku"`              // 立替金発生状況.振替額
	ThzHituyouNyukingaku        string                       `json:"sThzHituyouNyukingaku"`        // 立替金発生状況.必要入金額
	ThzKakuteiFlg               string                       `json:"sThzKakuteiFlg"`               // 立替金発生状況.確定フラグ, 1:入金請求管理.計算日＜営業日, 0:上記外
	GenbutuKabuKaituke          string                       `json:"sGenbutuKabuKaituke"`          // 株式現物買付可能額
	SinyouSinkidate             string                       `json:"sSinyouSinkidate"`             // 信用新規建可能額
	SinyouGenbiki               string                       `json:"sSinyouGenbiki"`               // 信用現引可能額
	HosyouKinritu               string                       `json:"sHosyouKinritu"`               // 委託保証金率(%)
	NseityouTousiKanougaku      string                       `json:"sNseityouTousiKanougaku"`      // NISA成長投資可能額
	TousinKaituke               string                       `json:"sTousinKaituke"`               // 投信買付可能額
	RuitouKaituke               string                       `json:"sRuitouKaituke"`               // MMF・中国F買付
	IPOKounyu                   string                       `json:"sIPOKounyu"`                   // IPO購入可能額
	Syukkin                     string                       `json:"sSyukkin"`                     // 出金可能額
	Fusokugaku                  string                       `json:"sFusokugaku"`                  // 不足額(入金請求額）
	LargeKaidateYoryoku         string                       `json:"sLargeKaidateYoryoku"`         // 先物買建
	MiniKaidateYoryoku          string                       `json:"sMiniKaidateYoryoku"`          // OPプット売建(ミニ)
	LargeUridateYoryoku         string                       `json:"sLargeUridateYoryoku"`         // 先物売建
	MiniUridateYoryoku          string                       `json:"sMiniUridateYoryoku"`          // OPコール売建(ミニ)
	OpKaidateYoryokyu           string                       `json:"sOpKaidateYoryokyu"`           // オプション新規買建
	SyoukokinFusokugaku         string                       `json:"sSyoukokinFusokugaku"`         // 証拠金不足額（本日請求額）
	GenbutuBaibaiDaikin         string                       `json:"sGenbutuBaibaiDaikin"`         // 現物売買代金
	GenbutuOrderCount           string                       `json:"sGenbutuOrderCount"`           // 現物注文件数
	GenbutuYakuzyouCount        string                       `json:"sGenbutuYakuzyouCount"`        // 現物約定件数
	SinyouBaibaiDaikin          string                       `json:"sSinyouBaibaiDaikin"`          // 信用売買代金
	SinyouOrderCount            string                       `json:"sSinyouOrderCount"`            // 信用注文件数
	SinyouYakuzyouCount         string                       `json:"sSinyouYakuzyouCount"`         // 信用約定件数
	SakiBaibaiDaikin            string                       `json:"sSakiBaibaiDaikin"`            // 先物売買代金
	SakiOrderCount              string                       `json:"sSakiOrderCount"`              // 先物注文件数
	SakiYakuzyouCount           string                       `json:"sSakiYakuzyouCount"`           // 先物約定件数
	OpBaibaiDaikin              string                       `json:"sOpBaibaiDaikin"`              // オプション売買代金
	OpOrderCount                string                       `json:"sOpOrderCount"`                // オプション注文件数
	OpYakuzyouCount             string                       `json:"sOpYakuzyouCount"`             // オプション約定件数
	HikazeiKouzaList            []ResHikazeiKouza            `json:"aHikazeiKouzaList"`            // 非課税口座リスト
	OisyouHasseiZyoukyouList    []ResOisyouHasseiZyoukyou    `json:"aOisyouHasseiZyoukyouList"`    // 追証発生状況リスト
	HosyoukinSeikyuZyoukyouList []ResHosyoukinSeikyuZyoukyou `json:"aHosyoukinSeikyuZyoukyouList"` // 保証金請求発生状況リスト
}

// ResHikazeiKouza 非課税口座リストの要素
type ResHikazeiKouza struct {
	HikazeiTekiyouYear    string `json:"sHikazeiTekiyouYear"`    // 適用年（対象年）, YYYY
	SeityouTousiKanougaku string `json:"sSeityouTousiKanougaku"` // 成長投資可能額
}

// ResOisyouHasseiZyoukyou 追証発生状況リストの要素
type ResOisyouHasseiZyoukyou struct {
	OhzHasseiDay           string `json:"sOhzHasseiDay"`           // 発生日 YYYYMMDD
	OhzHosyoukinRitu       string `json:"sOhzHosyoukinRitu"`       // 保証金率(%)
	OhzNyukinKigenDay      string `json:"sOhzNyukinKigenDay"`      // 入金期限 YYYYMMDDHHMM
	OhzOisyouKingaku       string `json:"sOhzOisyouKingaku"`       // 追証金額
	OhzKakuteiFlg          string `json:"sOhzKakuteiFlg"`          // 確定フラグ 1:入金請求管理.計算日＜営業日 0:上記外
	OhzHosyoukinZougen     string `json:"sOhzHosyoukinZougen"`     // 保証金増減
	OhzNyukin              string `json:"sOhzNyukin"`              // 入金
	OhzTategyokuKessai     string `json:"sOhzTategyokuKessai"`     // 建玉決済
	OhzKessaisonNyukin     string `json:"sOhzKessaisonNyukin"`     // 決済損入金
	OhzMikaisyouKingaku    string `json:"sOhzMikaisyouKingaku"`    // 未解消金額
	OhzMikaisyouKingakuFlg string `json:"sOhzMikaisyouKingakuFlg"` // 未解消金額フラグ 未使用
}

// ResHosyoukinSeikyuZyoukyou 保証金請求発生状況リストの要素
type ResHosyoukinSeikyuZyoukyou struct {
	HshzNyukinKigenDay           string `json:"sHshzNyukinKigenDay"`           // 入金期限 YYYYMMDDHHMM
	HshzHosyoukinHasseiDay       string `json:"sHshzHosyoukinHasseiDay"`       // 保証金発生日 YYYYMMDDHHMM
	HshzHosyoukin                string `json:"sHshzHosyoukin"`                // 保証金
	HshzGenkinHosyoukinHasseiDay string `json:"sHshzGenkinHosyoukinHasseiDay"` // 現金保証金発生日 YYYYMMDDHHMM
	HshzGenkinHosyoukin          string `json:"sHshzGenkinHosyoukin"`          // 現金保証金
}
