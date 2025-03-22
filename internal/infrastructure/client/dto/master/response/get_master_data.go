// internal/infrastructure/client/dto/master/response/get_master_data.go
package response

// GetMasterDataResponse は、マスタ情報問合取得の応答を表すDTOです。
// 各マスタ情報の配列をフィールドとして持ちます。
type GetMasterDataResponse struct {
	CLMID                     string                   `json:"sCLMID"`
	StockMaster               []StockMaster            `json:"CLMIssueMstKabu,omitempty"`
	StockMarketMaster         []StockMarketMaster      `json:"CLMIssueSizyouMstKabu,omitempty"`
	FutureMaster              []FutureMaster           `json:"CLMIssueMstSak,omitempty"`
	OptionMaster              []OptionMaster           `json:"CLMIssueMstOp,omitempty"`
	StockIssueRegulation      []StockIssueRegulation   `json:"CLMIssueSizyouKiseiKabu,omitempty"`
	FutureOptionRegulation    []FutureOptionRegulation `json:"CLMIssueSizyouKiseiHasei,omitempty"`
	MarginRate                []MarginRate             `json:"CLMDaiyouKakeme,omitempty"`
	MarginMaster              []MarginMaster           `json:"CLMHosyoukinMst,omitempty"`
	ErrorReason               []ErrorReason            `json:"CLMOrderErrReason,omitempty"`
	DateInfo                  []DateInfo               `json:"CLMDateZyouhou,omitempty"`
	SystemStatus              []SystemStatus           `json:"CLMSystemStatus,omitempty"`
	TickRule                  []TickRule               `json:"CLMYobine,omitempty"`
	OperationStatus           []OperationStatus        `json:"CLMUnyouStatus,omitempty"`
	OperationStatusStock      []OperationStatus        `json:"CLMUnyouStatusKabu,omitempty"`  // 運用ステータス（株）
	OperationStatusDerivative []OperationStatus        `json:"CLMUnyouStatusHasei,omitempty"` // 運用ステータス（派生）
	// ... 他にも必要なマスタ情報があれば、ここに追加 ...
}
