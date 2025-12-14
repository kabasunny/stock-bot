// domain/model/stock_master.go
package model

// StockMaster は、株式銘柄マスタの情報を表すモデル
// 株式銘柄の基本情報 CLMIssueMstKabu の情報に対応
type StockMaster struct {
	MasterBase
	IssueCode               string  `gorm:"primaryKey;size:255"` // 銘柄コード (主キー)
	IssueName               string  `gorm:"size:255"`             // 銘柄名称
	IssueNameShort          string  `gorm:"size:255"`             // 銘柄名略称
	IssueNameKana           string  `gorm:"size:255"`             // 銘柄名（カナ）
	IssueNameEnglish        string  `gorm:"size:255"`             // 銘柄名（英語表記）
	MarketCode              string  `gorm:"size:255"`             // 優先市場コード
	IndustryCode            string  `gorm:"size:255"`             // 業種コード
	IndustryName            string  `gorm:"size:255"`             // 業種コード名
	TradingUnit             int     // 売買単位
	ListedSharesOutstanding int64   // 上場発行株数
	UpperLimit              float64 // 値幅上限
	LowerLimit              float64 // 値幅下限
	TickRules               []TickRule `gorm:"foreignKey:IssueCode;references:IssueCode"` // 呼値 (1対多の関係)
}
