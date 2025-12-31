package response

// GetCreditInfoResponse は、信用残情報問合取得のレスポンスを表すDTO
// internal/infrastructure/client/dto/master/response/get_credit_info.go

type ResGetCreditInfo struct {
	CLMID             string                  `json:"sCLMID"`                       // 機能ID
	CLMMfdsShinyouZan []ResCreditInfoListItem `json:"aCLMMfdsShinyouZan,omitempty"` // 取得リスト
}

type ResCreditInfoListItem struct {
	IssueCode string `json:"sIssueCode"`      // 対象銘柄コード
	PMBB3     string `json:"pMBB3,omitempty"` // 信用残買残(一般)
	PMBB6     string `json:"pMBB6,omitempty"` // 信用残買残(制度)
	PMBBQ     string `json:"pMBBQ,omitempty"` // 信用残買残(合算)
	PMBC3     string `json:"pMBC3,omitempty"` // 信用残売残前週比(一般)
	PMBC6     string `json:"pMBC6,omitempty"` // 信用残売残前週比(制度)
	PMBCQ     string `json:"pMBCQ,omitempty"` // 信用残売残前週比(合算)
	PMBD      string `json:"pMBD,omitempty"`  // 信用残日付 YYYY/MM/DD
	PMBN3     string `json:"pMBN3,omitempty"` // 信用残買残前週比(一般)
	PMBN6     string `json:"pMBN6,omitempty"` // 信用残買残前週比(制度)
	PMBNQ     string `json:"pMBNQ,omitempty"` // 信用残買残前週比(合算)
	PMBR3     string `json:"pMBR3,omitempty"` // 信用倍率(一般)
	PMBR6     string `json:"pMBR6,omitempty"` // 信用倍率(制度)
	PMBRQ     string `json:"pMBRQ,omitempty"` // 信用倍率(合算)
	PMBS3     string `json:"pMBS3,omitempty"` // 信用残売残(一般)
	PMBS6     string `json:"pMBS6,omitempty"` // 信用残売残(制度)
	PMBSQ     string `json:"pMBSQ,omitempty"` // 信用残売残(合算)
}
