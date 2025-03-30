// response/genbutu_kabu_list.go
package response

// ResGenbutuKabuList は現物保有銘柄一覧のレスポンスを表すDTO
type ResGenbutuKabuList struct {
	P_no                             string           `json:"p_no"`                              // p_no
	CLMID                            string           `json:"sCLMID"`                            // 機能ID, CLMGenbutuKabuList
	ResultCode                       string           `json:"sResultCode"`                       // 結果コード, CLMKabuNewOrder.sResultCode 参照
	ResultText                       string           `json:"sResultText"`                       // 結果テキスト, CLMKabuNewOrder.sResultText 参照
	WarningCode                      string           `json:"sWarningCode"`                      // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	WarningText                      string           `json:"sWarningText"`                      // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	IssueCode                        string           `json:"sIssueCode"`                        // 銘柄コード, 要求設定値
	IppanGaisanHyoukagakuGoukei      string           `json:"sIppanGaisanHyoukagakuGoukei"`      // 概算評価額合計　(一般口座残高)
	IppanGaisanHyoukaSonekiGoukei    string           `json:"sIppanGaisanHyoukaSonekiGoukei"`    // 概算評価損益合計(一般口座残高)
	NisaGaisanHyoukagakuGoukei       string           `json:"sNisaGaisanHyoukagakuGoukei"`       // 概算評価額合計　(NISA口座残高)
	NisaGaisanHyoukaSonekiGoukei     string           `json:"sNisaGaisanHyoukaSonekiGoukei"`     // 概算評価損益合計(NISA口座残高)
	NseityouGaisanHyoukagakuGoukei   string           `json:"sNseityouGaisanHyoukagakuGoukei"`   // 概算評価額合計　(N成長口座残高)
	NseityouGaisanHyoukaSonekiGoukei string           `json:"sNseityouGaisanHyoukaSonekiGoukei"` // 概算評価損益合計(N成長口座残高)
	TokuteiGaisanHyoukagakuGoukei    string           `json:"sTokuteiGaisanHyoukagakuGoukei"`    // 概算評価額合計　(特定口座残高)
	TokuteiGaisanHyoukaSonekiGoukei  string           `json:"sTokuteiGaisanHyoukaSonekiGoukei"`  // 概算評価損益合計(特定口座残高)
	TotalGaisanHyoukagakuGoukei      string           `json:"sTotalGaisanHyoukagakuGoukei"`      // 概算評価額合計　(残高合計)
	TotalGaisanHyoukaSonekiGoukei    string           `json:"sTotalGaisanHyoukaSonekiGoukei"`    // 概算評価損益合計(残高合計)
	GenbutuKabuList                  []ResGenbutuKabu `json:"aGenbutuKabuList"`                  // 現物保有リスト
}

// ResGenbutuKabu 現物保有リストの要素
type ResGenbutuKabu struct {
	UriOrderWarningCode            string `json:"sUriOrderWarningCode"`            // 警告コード, CLMKabuNewOrder.sWarningCode 参照
	UriOrderWarningText            string `json:"sUriOrderWarningText"`            // 警告テキスト, CLMKabuNewOrder.sWarningTexts 参照
	UriOrderIssueCode              string `json:"sUriOrderIssueCode"`              // 銘柄コード, 保有銘柄コード
	UriOrderZyoutoekiKazeiC        string `json:"sUriOrderZyoutoekiKazeiC"`        // 譲渡益課税区分, CLMKabuNewOrder.sZyoutoekiKazeiC 参照
	UriOrderZanKabuSuryou          string `json:"sUriOrderZanKabuSuryou"`          // 残高株数
	UriOrderUritukeKanouSuryou     string `json:"sUriOrderUritukeKanouSuryou"`     // 売付可能株数
	UriOrderGaisanBokaTanka        string `json:"sUriOrderGaisanBokaTanka"`        // 概算簿価単価
	UriOrderHyoukaTanka            string `json:"sUriOrderHyoukaTanka"`            // 評価単価
	UriOrderGaisanHyoukagaku       string `json:"sUriOrderGaisanHyoukagaku"`       // 評価金額
	UriOrderGaisanHyoukaSoneki     string `json:"sUriOrderGaisanHyoukaSoneki"`     // 評価損益
	UriOrderGaisanHyoukaSonekiRitu string `json:"sUriOrderGaisanHyoukaSonekiRitu"` // 評価損益率(%)
	SyuzituOwarine                 string `json:"sSyuzituOwarine"`                 // 前日終値, 該当値取得不可時:""
	ZenzituHi                      string `json:"sZenzituHi"`                      // 前日比, 該当値取得不可時:""
	ZenzituHiPer                   string `json:"sZenzituHiPer"`                   // 前日比(%), 該当値取得不可時:""
	UpDownFlag                     string `json:"sUpDownFlag"`                     // 騰落率Flag(%), 前日比(%)取得不可時:""
	NissyoukinKasikabuZan          string `json:"sNissyoukinKasikabuZan"`          // 証金貸株残, 該当値取得不可時:""
}
