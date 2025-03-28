package response

// 時価情報問合取得 レスポンス
// internal/infrastructure/client/dto/price/response/get_price_info.go

type ResGetPriceInfo struct {
	CLMID              string                   `json:"sCLMID"`                        // 機能ID (CLMMfdsGetMarketPrice)
	CLMMfdsMarketPrice []ResMarketPriceInfoItem `json:"aCLMMfdsMarketPrice,omitempty"` // 取得リスト
}

// ResMarketPriceInfoItem は、時価情報問合取得のレスポンスの各項目を表す構造体
type ResMarketPriceInfoItem struct {
	IssueCode string `json:"sIssueCode"` // 対象銘柄コード
	// TargetColumn で指定した情報コードに対応するフィールド
	// 情報コードが可変なので、map[string]string で表現する
	Values map[string]string `json:"-"` // Values は、JSONに直接は含めない
}
