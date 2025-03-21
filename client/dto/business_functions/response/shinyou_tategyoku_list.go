// business_functions/res_shinyou_tategyoku_list.go
package business_functions

// ResShinyouTategyokuList は信用建玉一覧のレスポンスを表すDTO
type ResShinyouTategyokuList struct {
	P_no                       string                `json:"p_no"`                       // p_no
	SCLMID                     string                `json:"sCLMID"`                     // 機能ID, CLMShinyouTategyokuList
	SResultCode                string                `json:"sResultCode"`                // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                string                `json:"sResultText"`                // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode               string                `json:"sWarningCode"`               // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText               string                `json:"sWarningText"`               // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode                 string                `json:"sIssueCode"`                 // 銘柄コード, 要求設定値
	SUritateDaikin             string                `json:"sUritateDaikin"`             // 売建代金合計
	SKaitateDaikin             string                `json:"sKaitateDaikin"`             // 買建代金合計
	STotalDaikin               string                `json:"sTotalDaikin"`               // 総代金合計
	SHyoukaSonekiGoukeiUridate string                `json:"sHyoukaSonekiGoukeiUridate"` // 評価損益合計_売建
	SHyoukaSonekiGoukeiKaidate string                `json:"sHyoukaSonekiGoukeiKaidate"` // 評価損益合計_買建
	STotalHyoukaSonekiGoukei   string                `json:"sTotalHyoukaSonekiGoukei"`   // 総評価損益合計
	STokuteiHyoukaSonekiGoukei string                `json:"sTokuteiHyoukaSonekiGoukei"` // 特定口座残高評価損益合計
	SIppanHyoukaSonekiGoukei   string                `json:"sIppanHyoukaSonekiGoukei"`   // 一般口座残高評価損益合計
	ASinyouTategyokuList       []ResShinyouTategyoku `json:"aShinyouTategyokuList"`      // 信用建玉リスト
}

// ResShinyouTategyoku 信用建玉リストの要素
type ResShinyouTategyoku struct {
	SOrderWarningCode            string `json:"sOrderWarningCode"`            // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SOrderWarningText            string `json:"sOrderWarningText"`            // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SOrderTategyokuNumber        string `json:"sOrderTategyokuNumber"`        // 建玉番号, 保有建玉番号
	SOrderIssueCode              string `json:"sOrderIssueCode"`              // 銘柄コード, 保有銘柄コード
	SOrderSizyouC                string `json:"sOrderSizyouC"`                // 市場, 00：東証
	SOrderBaibaiKubun            string `json:"sOrderBaibaiKubun"`            // 売買区分, CLMKabuNewOrder.sBaibaiKubun 参照
	SOrderBensaiKubun            string `json:"sOrderBensaiKubun"`            // 弁済区分, 00：なし, 26：制度信用6ヶ月, 29：制度信用無期限, 36：一般信用6ヶ月, 39：一般信用無期限
	SOrderZyoutoekiKazeiC        string `json:"sOrderZyoutoekiKazeiC"`        // 譲渡益課税区分, 1：特定, 3：一般, 5：NISA, 9：法人
	SOrderTategyokuSuryou        string `json:"sOrderTategyokuSuryou"`        // 建株数
	SOrderTategyokuTanka         string `json:"sOrderTategyokuTanka"`         // 建単価
	SOrderHyoukaTanka            string `json:"sOrderHyoukaTanka"`            // 評価単価
	SOrderGaisanHyoukaSoneki     string `json:"sOrderGaisanHyoukaSoneki"`     // 評価損益
	SOrderGaisanHyoukaSonekiRitu string `json:"sOrderGaisanHyoukaSonekiRitu"` // 評価損益率(%)
	STategyokuDaikin             string `json:"sTategyokuDaikin"`             // 建玉代金
	SOrderTateTesuryou           string `json:"sOrderTateTesuryou"`           // 建手数料
	SOrderZyunHibu               string `json:"sOrderZyunHibu"`               // 順日歩
	SOrderGyakuhibu              string `json:"sOrderGyakuhibu"`              // 逆日歩
	SOrderKakikaeryou            string `json:"sOrderKakikaeryou"`            // 書換料
	SOrderKanrihi                string `json:"sOrderKanrihi"`                // 管理費
	SOrderKasikaburyou           string `json:"sOrderKasikaburyou"`           // 貸株料
	SOrderSonota                 string `json:"sOrderSonota"`                 // その他
	SOrderTategyokuDay           string `json:"sOrderTategyokuDay"`           // 建日, YYYYMMDD, 00000000
	SOrderTategyokuKizituDay     string `json:"sOrderTategyokuKizituDay"`     // 建玉期日日, YYYYMMDD, 00000000：無期限
	STategyokuSuryou             string `json:"sTategyokuSuryou"`             // 建玉数量
	SOrderYakuzyouHensaiKabusu   string `json:"sOrderYakuzyouHensaiKabusu"`   // 約定返済株数
	SOrderGenbikiGenwatasiKabusu string `json:"sOrderGenbikiGenwatasiKabusu"` // 現引現渡株数
	SOrderOrderSuryou            string `json:"sOrderOrderSuryou"`            // 注文中数量
	SOrderHensaiKanouSuryou      string `json:"sOrderHensaiKanouSuryou"`      // 返済可能数量
	SSyuzituOwarine              string `json:"sSyuzituOwarine"`              // 前日終値, 該当値取得不可時:""
	SZenzituHi                   string `json:"sZenzituHi"`                   // 前日比, 該当値取得不可時:""
	SZenzituHiPer                string `json:"sZenzituHiPer"`                // 前日比(%), 該当値取得不可時:""
	SUpDownFlag                  string `json:"sUpDownFlag"`                  // 騰落率Flag(%), 前日比(%)取得不可時:""
}
