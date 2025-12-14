// domain/model/stock_master.go
package model

// StockMaster は、株式銘柄マスタの情報を表すモデル
// 株式銘柄の基本情報 CLMIssueMstKabu の情報に対応
type StockMaster struct {
	MasterBase
	IssueCode   string  `gorm:"primaryKey;size:255"` // 銘柄コード (主キー)
	IssueName   string  `gorm:"size:255"`             // 銘柄名称
	TradingUnit int     // 売買単位
	MarketCode  string  `gorm:"size:255"` // 市場コード
	UpperLimit  float64 // 値幅上限
	LowerLimit  float64 // 値幅下限
	//PreviousClose float64 // 前日終値 (必要に応じて)　削除可能
	TickRules []TickRule `gorm:"foreignKey:IssueCode;references:IssueCode"` // 呼値 (1対多の関係)
}
