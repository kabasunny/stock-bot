package response

// 蓄積情報問合取得 レスポンス
// internal/infrastructure/client/dto/price/response/get_price_info_history.go

type ResGetPriceInfoHistory struct {
	SCLMID                        string                          `json:"sCLMID"`                                  // 機能ID (CLMMfdsGetMarketPriceHistory)
	SIssueCode                    string                          `json:"sIssueCode"`                              // 銘柄コード
	SSizyouC                      string                          `json:"sSizyouC"`                                // 市場コード
	ACLMMfdsGetMarketPriceHistory []ResMarketPriceHistoryInfoItem `json:"aCLMMfdsGetMarketPriceHistory,omitempty"` // 取得リスト
}

type ResMarketPriceHistoryInfoItem struct {
	SDate  string `json:"sDate"`            // 日付 YYYYMMDD
	PDOP   string `json:"pDOP,omitempty"`   // 始値
	PDHP   string `json:"pDHP,omitempty"`   // 高値
	PDLP   string `json:"pDLP,omitempty"`   // 安値
	PDPP   string `json:"pDPP,omitempty"`   // 終値
	PDV    string `json:"pDV,omitempty"`    // 出来高
	PDOPxK string `json:"pDOPxK,omitempty"` // 始値ｘ分割係数
	PDHPxK string `json:"pDHPxK,omitempty"` // 高値ｘ分割係数
	PDLPxK string `json:"pDLPxK,omitempty"` // 安値ｘ分割係数
	PDPPxK string `json:"pDPPxK,omitempty"` // 終値ｘ分割係数
	PDVxK  string `json:"pDVxK,omitempty"`  // 出来高÷分割係数
	PSPUO  string `json:"pSPUO,omitempty"`  // 分割前単位
	PSPUC  string `json:"pSPUC,omitempty"`  // 分割後単位
	PSPUK  string `json:"pSPUK,omitempty"`  // 分割換算係数
}
