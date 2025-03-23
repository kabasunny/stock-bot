package response

// ニュースヘッダー問合取得 レスポンス
// internal/infrastructure/client/dto/master/response/get_news_head.go

type ResGetNewsHeader struct {
	CLMID           string                `json:"sCLMID"`                     // 機能ＩＤ (CLMMfdsGetNewsHead)
	PRECMAX         string                `json:"p_REC_MAX"`                  // レコード件数
	CLMMfdsNewsHead []ResNewsHeadListItem `json:"aCLMMfdsNewsHead,omitempty"` // 取得リスト
}

type ResNewsHeadListItem struct {
	PID  string `json:"p_ID"`  // ニュースＩＤ
	PDT  string `json:"p_DT"`  // ニュース日付 YYYYMMDD
	PTM  string `json:"p_TM"`  // ニュース時刻 HHMMSS
	PCGL string `json:"p_CGL"` // ニュースカテゴリリスト (複数設定時は「|」区切り)
	PGNL string `json:"p_GNL"` // ニュースジャンルリスト (複数設定時は「|」区切り)
	PISL string `json:"p_ISL"` // ニュース関連銘柄コードリスト (複数設定時は「|」区切り)
	PHDL string `json:"p_HDL"` // ニュースヘッドライン（タイトル）(BASE64エンコード)
}
