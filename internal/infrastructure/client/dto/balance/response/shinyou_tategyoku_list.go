// response/shinyou_tategyoku_list.go
package response

// ResShinyouTategyokuList は信用建玉一覧のレスポンスを表すDTO
type ResShinyouTategyokuList struct {
	P_no                      string                `json:"p_no"`                       // p_no
	CLMID                     string                `json:"sCLMID"`                     // 機能ID, CLMShinyouTategyokuList
	ResultCode                string                `json:"sResultCode"`                // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText                string                `json:"sResultText"`                // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode               string                `json:"sWarningCode"`               // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText               string                `json:"sWarningText"`               // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	IssueCode                 string                `json:"sIssueCode"`                 // 銘柄コード, 要求設定値
	UritateDaikin             string                `json:"sUritateDaikin"`             // 売建代金合計
	KaitateDaikin             string                `json:"sKaitateDaikin"`             // 買建代金合計
	TotalDaikin               string                `json:"sTotalDaikin"`               // 総代金合計
	HyoukaSonekiGoukeiUridate string                `json:"sHyoukaSonekiGoukeiUridate"` // 評価損益合計_売建
	HyoukaSonekiGoukeiKaidate string                `json:"sHyoukaSonekiGoukeiKaidate"` // 評価損益合計_買建
	TotalHyoukaSonekiGoukei   string                `json:"sTotalHyoukaSonekiGoukei"`   // 総評価損益合計
	TokuteiHyoukaSonekiGoukei string                `json:"sTokuteiHyoukaSonekiGoukei"` // 特定口座残高評価損益合計
	IppanHyoukaSonekiGoukei   string                `json:"sIppanHyoukaSonekiGoukei"`   // 一般口座残高評価損益合計
	SinyouTategyokuList       []ResShinyouTategyoku `json:"aShinyouTategyokuList"`      // 信用建玉リスト
}

// ResShinyouTategyoku 信用建玉リストの要素
type ResShinyouTategyoku struct {
	OrderWarningCode            string `json:"sOrderWarningCode"`            // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	OrderWarningText            string `json:"sOrderWarningText"`            // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	OrderTategyokuNumber        string `json:"sOrderTategyokuNumber"`        // 建玉番号, 保有建玉番号
	OrderIssueCode              string `json:"sOrderIssueCode"`              // 銘柄コード, 保有銘柄コード
	OrderSizyouC                string `json:"sOrderSizyouC"`                // 市場, 00：東証
	OrderBaibaiKubun            string `json:"sOrderBaibaiKubun"`            // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	OrderBensaiKubun            string `json:"sOrderBensaiKubun"`            // 弁済区分, 00：なし, 26：制度信用6ヶ月, 29：制度信用無期限, 36：一般信用6ヶ月, 39：一般信用無期限
	OrderZyoutoekiKazeiC        string `json:"sOrderZyoutoekiKazeiC"`        // 譲渡益課税区分, 1：特定, 3：一般, 5：NISA, 9：法人
	OrderTategyokuSuryou        string `json:"sOrderTategyokuSuryou"`        // 建株数
	OrderTategyokuTanka         string `json:"sOrderTategyokuTanka"`         // 建単価
	OrderHyoukaTanka            string `json:"sOrderHyoukaTanka"`            // 評価単価
	OrderGaisanHyoukaSoneki     string `json:"sOrderGaisanHyoukaSoneki"`     // 評価損益
	OrderGaisanHyoukaSonekiRitu string `json:"sOrderGaisanHyoukaSonekiRitu"` // 評価損益率(%)
	TategyokuDaikin             string `json:"sTategyokuDaikin"`             // 建玉代金
	OrderTateTesuryou           string `json:"sOrderTateTesuryou"`           // 建手数料
	OrderZyunHibu               string `json:"sOrderZyunHibu"`               // 順日歩
	OrderGyakuhibu              string `json:"sOrderGyakuhibu"`              // 逆日歩
	OrderKakikaeryou            string `json:"sOrderKakikaeryou"`            // 書換料
	OrderKanrihi                string `json:"sOrderKanrihi"`                // 管理費
	OrderKasikaburyou           string `json:"sOrderKasikaburyou"`           // 貸株料
	OrderSonota                 string `json:"sOrderSonota"`                 // その他
	OrderTategyokuDay           string `json:"sOrderTategyokuDay"`           // 建日, YYYYMMDD, 00000000
	OrderTategyokuKizituDay     string `json:"sOrderTategyokuKizituDay"`     // 建玉期日日, YYYYMMDD, 00000000：無期限
	TategyokuSuryou             string `json:"sTategyokuSuryou"`             // 建玉数量
	OrderYakuzyouHensaiKabusu   string `json:"sOrderYakuzyouHensaiKabusu"`   // 約定返済株数
	OrderGenbikiGenwatasiKabusu string `json:"sOrderGenbikiGenwatasiKabusu"` // 現引現渡株数
	OrderOrderSuryou            string `json:"sOrderOrderSuryou"`            // 注文中数量
	OrderHensaiKanouSuryou      string `json:"sOrderHensaiKanouSuryou"`      // 返済可能数量
	SyuzituOwarine              string `json:"sSyuzituOwarine"`              // 前日終値, 該当値取得不可時:""
	ZenzituHi                   string `json:"sZenzituHi"`                   // 前日比, 該当値取得不可時:""
	ZenzituHiPer                string `json:"sZenzituHiPer"`                // 前日比(%), 該当値取得不可時:""
	UpDownFlag                  string `json:"sUpDownFlag"`                  // 騰落率Flag(%), 前日比(%)取得不可時:""
}
