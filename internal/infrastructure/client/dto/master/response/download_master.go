// internal/infrastructure/client/dto/master/response/get_master_data.go
package response

// GetMasterDataResponse は、マスタ情報問合取得の応答を表すDTO
// 各マスタ情報の配列をフィールドとして持つ
type ResDownloadMaster struct {
	CLMID                     string                      `json:"sCLMID"`
	SystemStatus              ResSystemStatus             `json:"CLMSystemStatus,omitempty"`
	DateInfo                  []ResDateInfo               `json:"CLMDateZyouhou,omitempty"`
	TickRule                  []ResTickRule               `json:"CLMYobine,omitempty"`
	OperationStatus           []ResOperationStatus        `json:"CLMUnyouStatus,omitempty"`      // 配列に変更
	OperationStatusStock      []ResOperationStatus        `json:"CLMUnyouStatusKabu,omitempty"`  // 配列に変更
	OperationStatusDerivative []ResOperationStatus        `json:"CLMUnyouStatusHasei,omitempty"` // 配列に変更
	StockMaster               []ResStockMaster            `json:"CLMIssueMstKabu,omitempty"`
	StockMarketMaster         []ResStockMarketMaster      `json:"CLMIssueSizyouMstKabu,omitempty"`
	StockIssueRegulation      []ResStockIssueRegulation   `json:"CLMIssueSizyouKiseiKabu,omitempty"`
	FutureMaster              []ResFutureMaster           `json:"CLMIssueMstSak,omitempty"`
	OptionMaster              []ResOptionMaster           `json:"CLMIssueMstOp,omitempty"`
	FutureOptionRegulation    []ResFutureOptionRegulation `json:"CLMIssueSizyouKiseiHasei,omitempty"`
	MarginRate                []ResMarginRate             `json:"CLMDaiyouKakeme,omitempty"`
	MarginMaster              []ResMarginMaster           `json:"CLMHosyoukinMst,omitempty"`
	ErrorReason               []ResErrorReason            `json:"CLMOrderErrReason,omitempty"`
	DownloadComplete          ResDownloadComplete         `json:"CLMEventDownloadComplete,omitempty"` // これがない場合は終了していない
	// ... 他にも必要なマスタ情報があれば、ここに追加 ...
}
