// internal/infrastructure/client/dto/master/response/option_master.go
package response

// ResOptionMaster は、オプション銘柄マスタの情報を表すDTO
type ResOptionMaster struct {
	CLMID                    string `json:"sCLMID"`               // 機能ID (CLMIssueMstOp)
	IssueCode                string `json:"sIssueCode"`           // 銘柄コード
	IssueName                string `json:"sIssueName"`           // 銘柄名
	IssueNameEnglish         string `json:"sIssueNameEizi"`       // 銘柄名（英語表記）
	FutureOptionProduct      string `json:"sSakOpSyouhin"`        // 先物ＯＰ商品
	UnderlyingAssetType      string `json:"sGensisanKubun"`       // 原資産区分
	UnderlyingAssetCode      string `json:"sGensisanCode"`        // 原資産コード
	ContractMonth            string `json:"sGengetu"`             // 限月 (YYYYMM)
	ListingMarket            string `json:"sZyouzyouSizyou"`      // 上場市場
	StrikePrice              string `json:"sKousiPrice"`          // 行使価格
	PutCall                  string `json:"sPutCall"`             // プット・コール (5：プット, 7：コール)
	TradingStartDate         string `json:"sTorihikiStartDay"`    // 取引開始日 (YYYYMMDD)
	LastTradingDay           string `json:"sLastBaibaiDay"`       // 最終売買日 (YYYYMMDD)
	LastExerciseDate         string `json:"sKenrikousiLastDay"`   // 権利行使最終日 (YYYYMMDD)
	UnitQuantity             string `json:"sTaniSuryou"`          // 単位数量
	TickUnitNumber           string `json:"sYobineTaniNumber"`    // 呼値の単位番号
	InformationSource        string `json:"sZyouhouSource"`       // 情報系ソース
	InformationCode          string `json:"sZyouhouCode"`         // 情報系コード
	LowerLimit               string `json:"sNehabaMin"`           // 値幅下限
	UpperLimit               string `json:"sNehabaMax"`           // 値幅上限
	IssueRegulation1         string `json:"sIssueKisei1C"`        // 銘柄規制１区分
	PreviousClose            string `json:"sZenzituOwarine"`      // 前日終値
	PreviousTheoreticalPrice string `json:"sZenzituRironPrice"`   // 前日理論価格
	FloorSlipOutputFlag      string `json:"sBaDenpyouOutputUmuC"` // 場伝票出力有無区分
	CreateDate               string `json:"sCreateDate"`          // 作成日時
	UpdateDate               string `json:"sUpdateDate"`          // 更新日時
	UpdateNumber             string `json:"sUpdateNumber"`        // 更新通番
	ATMFlag                  string `json:"sATMFlag"`             // アットザマネーフラグ
}
