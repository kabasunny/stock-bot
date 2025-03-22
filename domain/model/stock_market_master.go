// domain/model/stock_market_master.go
package model

// StockMarketMaster は、株式銘柄市場マスタの情報を表すモデルです。
// (必要に応じて)
type StockMarketMaster struct {
	MasterBase            // 共通フィールド
	IssueCode     string  `gorm:"index;size:255"` // 銘柄コード (複合ユニークキーの一部)
	ListingMarket string  `gorm:"index;size:255"` // 上場市場 (複合ユニークキーの一部)
	PreviousClose float64 // 前日終値
}
