// internal/infrastructure/client/dto/master/response/get_master_data_query.go
package response

// GetMasterDataQueryResponse は、マスタ情報問合取得のレスポンスを表すDTO
type ResGetMasterData struct {
	CLMID             string                 `json:"sCLMID"`                          // 機能ID (CLMMfdsGetMasterData)
	StockMaster       []ResStockMaster       `json:"CLMIssueMstKabu,omitempty"`       // 株式銘柄マスタ (CLMIssueMstKabu) の配列
	StockMarketMaster []ResStockMarketMaster `json:"CLMIssueSizyouMstKabu,omitempty"` // 株式市場マスタ (CLMIssueSizyouMstKabu) の配列
	FutureMaster      []ResFutureMaster      `json:"CLMIssueMstSak,omitempty"`        // 先物銘柄マスタ (CLMIssueMstSak) の配列
	OptionMaster      []ResOptionMaster      `json:"CLMIssueMstOp,omitempty"`         // オプション銘柄マスタ (CLMIssueMstOp) の配列
	IssueMstOther     []ResIssueMstOtherItem `json:"CLMIssueMstOther,omitempty"`      // その他銘柄マスタ (CLMIssueMstOther) の配列
	IssueMstIndex     []ResIssueMstIndexItem `json:"CLMIssueMstIndex,omitempty"`      // 指数マスタ (CLMIssueMstIndex) の配列
	IssueMstFx        []ResIssueMstFxItem    `json:"CLMIssueMstFx,omitempty"`         // 為替マスタ (CLMIssueMstFx) の配列
	ErrorReason       []ResErrorReason       `json:"CLMOrderErrReason,omitempty"`     // 注文エラー理由マスタ (CLMOrderErrReason) の配列
	DateInfo          []ResDateInfo          `json:"CLMDateZyouhou,omitempty"`        // 日付情報マスタ (CLMDateZyouhou) の配列
}

// ResIssueMstOtherItem は、その他銘柄マスタ (CLMIssueMstOther) の項目を表す構造体
type ResIssueMstOtherItem struct {
	IssueCode string `json:"sIssueCode,omitempty"` // 銘柄コード
	IssueName string `json:"sIssueName,omitempty"` // 銘柄名称
}

// ResIssueMstIndexItem は、指数マスタ (CLMIssueMstIndex) の項目を表す構造体
type ResIssueMstIndexItem struct {
	IssueCode string `json:"sIssueCode,omitempty"` // 指数コード
	IssueName string `json:"sIssueName,omitempty"` // 指数名称
}

// ResIssueMstFxItem は、為替マスタ (CLMIssueMstFx) の項目を表す構造体
type ResIssueMstFxItem struct {
	IssueCode string `json:"sIssueCode,omitempty"` // 通貨コード (例: USDJPY)
	IssueName string `json:"sIssueName,omitempty"` // 通貨名称 (例: ドル円)
}
