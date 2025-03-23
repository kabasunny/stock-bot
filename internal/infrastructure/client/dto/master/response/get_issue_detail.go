package response

// 銘柄詳細情報問合取得 レスポンス
// internal/infrastructure/client/dto/master/response/get_issue_detail.go

type ResGetIssueDetail struct {
	CLMID              string                   `json:"sCLMID"`                        // 機能ＩＤ (CLMMfdsGetIssueDetail)
	CLMMfdsIssueDetail []ResIssueDetailListItem `json:"aCLMMfdsIssueDetail,omitempty"` // 取得リスト
}

type ResIssueDetailListItem struct {
	IssueCode string `json:"sIssueCode"`      // 対象銘柄コード
	PBPSB     string `json:"pBPSB,omitempty"` // BPS（実績）／一株資産(実績最新・連結)
	PCLOE     string `json:"pCLOE,omitempty"` // 落日（本決算）／配当権利落日 YYYY/MM/DD
	PEPSF     string `json:"pEPSF,omitempty"` // EPS（予想）／一株利益(予想・通期連結)
	PEXRD     string `json:"pEXRD,omitempty"` // 最終落日（決算期以外） YYYY/MM/DD
	PIDVE     string `json:"pIDVE,omitempty"` // 落日（中間決算）／中間配当権利落日 YYYY/MM/DD
	PROEL     string `json:"pROEL,omitempty"` // ROE（予想）
	PRPER     string `json:"pRPER,omitempty"` // PER（予想）／連結優先　PER
	PSPBR     string `json:"pSPBR,omitempty"` // PBR（実績）／PBR(単純)
	PSPRO     string `json:"pSPRO,omitempty"` // 株式益回り（予想）／益回り(単純)
	PSYIE     string `json:"pSYIE,omitempty"` // 配当利回り（予想）／利回り(単純)
	PYHPD     string `json:"pYHPD,omitempty"` // 年初来高値：更新日 YYYY/MM/DD
	PYHPR     string `json:"pYHPR,omitempty"` // 年初来高値
	PYLPD     string `json:"pYLPD,omitempty"` // 年初来安値：更新日 YYYY/MM/DD
	PYLPR     string `json:"pYLPR,omitempty"` // 年初来安値
}
