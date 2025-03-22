// response/genbutu_kabu_list.go
package response

// ResGenbutuKabuList は現物保有銘柄一覧のレスポンスを表すDTO
type ResGenbutuKabuList struct {
	P_no                              string           `json:"p_no"`                              // p_no
	SCLMID                            string           `json:"sCLMID"`                            // 機能ID, CLMGenbutuKabuList
	SResultCode                       string           `json:"sResultCode"`                       // 結果コード, CLMKabuNewOrder.sResultCode 参照
	SResultText                       string           `json:"sResultText"`                       // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	SWarningCode                      string           `json:"sWarningCode"`                      // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SWarningText                      string           `json:"sWarningText"`                      // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SIssueCode                        string           `json:"sIssueCode"`                        // 銘柄コード, 要求設定値
	SIppanGaisanHyoukagakuGoukei      string           `json:"sIppanGaisanHyoukagakuGoukei"`      // 概算評価額合計　(一般口座残高)
	SIppanGaisanHyoukaSonekiGoukei    string           `json:"sIppanGaisanHyoukaSonekiGoukei"`    // 概算評価損益合計(一般口座残高)
	SNisaGaisanHyoukagakuGoukei       string           `json:"sNisaGaisanHyoukagakuGoukei"`       // 概算評価額合計　(NISA口座残高)
	SNisaGaisanHyoukaSonekiGoukei     string           `json:"sNisaGaisanHyoukaSonekiGoukei"`     // 概算評価損益合計(NISA口座残高)
	SNseityouGaisanHyoukagakuGoukei   string           `json:"sNseityouGaisanHyoukagakuGoukei"`   // 概算評価額合計　(N成長口座残高)
	SNseityouGaisanHyoukaSonekiGoukei string           `json:"sNseityouGaisanHyoukaSonekiGoukei"` // 概算評価損益合計(N成長口座残高)
	STokuteiGaisanHyoukagakuGoukei    string           `json:"sTokuteiGaisanHyoukagakuGoukei"`    // 概算評価額合計　(特定口座残高)
	STokuteiGaisanHyoukaSonekiGoukei  string           `json:"sTokuteiGaisanHyoukaSonekiGoukei"`  // 概算評価損益合計(特定口座残高)
	STotalGaisanHyoukagakuGoukei      string           `json:"sTotalGaisanHyoukagakuGoukei"`      // 概算評価額合計　(残高合計)
	STotalGaisanHyoukaSonekiGoukei    string           `json:"sTotalGaisanHyoukaSonekiGoukei"`    // 概算評価損益合計(残高合計)
	AGenbutuKabuList                  []ResGenbutuKabu `json:"aGenbutuKabuList"`                  // 現物保有リスト
}

// ResGenbutuKabu 現物保有リストの要素
type ResGenbutuKabu struct {
	SUriOrderWarningCode            string `json:"sUriOrderWarningCode"`            // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	SUriOrderWarningText            string `json:"sUriOrderWarningText"`            // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	SUriOrderIssueCode              string `json:"sUriOrderIssueCode"`              // 銘柄コード, 保有銘柄コード
	SUriOrderZyoutoekiKazeiC        string `json:"sUriOrderZyoutoekiKazeiC"`        // 譲渡益課税区分, CLMKabuNewOrder.sZyoutoekiKazeiC 参照
	SUriOrderZanKabuSuryou          string `json:"sUriOrderZanKabuSuryou"`          // 残高株数
	SUriOrderUritukeKanouSuryou     string `json:"sUriOrderUritukeKanouSuryou"`     // 売付可能株数
	SUriOrderGaisanBokaTanka        string `json:"sUriOrderGaisanBokaTanka"`        // 概算簿価単価
	SUriOrderHyoukaTanka            string `json:"sUriOrderHyoukaTanka"`            // 評価単価
	SUriOrderGaisanHyoukagaku       string `json:"sUriOrderGaisanHyoukagaku"`       // 評価金額
	SUriOrderGaisanHyoukaSoneki     string `json:"sUriOrderGaisanHyoukaSoneki"`     // 評価損益
	SUriOrderGaisanHyoukaSonekiRitu string `json:"sUriOrderGaisanHyoukaSonekiRitu"` // 評価損益率(%)
	SSyuzituOwarine                 string `json:"sSyuzituOwarine"`                 // 前日終値, 該当値取得不可時:""
	SZenzituHi                      string `json:"sZenzituHi"`                      // 前日比, 該当値取得不可時:""
	SZenzituHiPer                   string `json:"sZenzituHiPer"`                   // 前日比(%), 該当値取得不可時:""
	SUpDownFlag                     string `json:"sUpDownFlag"`                     // 騰落率Flag(%), 前日比(%)取得不可時:""
	SNissyoukinKasikabuZan          string `json:"sNissyoukinKasikabuZan"`          // 証金貸株残, 該当値取得不可時:""
}
