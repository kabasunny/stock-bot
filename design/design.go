package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("stock", func() {
	Title("株式取引ボット API")
	Description("株式取引ボットの機能を提供するAPI")
	Server("stock-bot", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

// BalanceSummary は口座サマリーのデータ構造を定義します。
var BalanceSummary = ResultType("application/vnd.stock.balance.summary", func() {
	Description("口座の残高サマリー情報")
	Attributes(func() {
		Attribute("total_assets", Float64, "総資産 (円)", func() {
			Example(3000000.50)
		})
		Attribute("cash_buying_power", Float64, "現物買付可能額 (円)", func() {
			Example(1000000)
		})
		Attribute("margin_buying_power", Float64, "信用新規建可能額 (円)", func() {
			Example(5000000)
		})
		Attribute("withdrawal_possible_amount", Float64, "出金可能額 (円)", func() {
			Example(500000)
		})
		Attribute("margin_rate", Float64, "委託保証金率 (%)", func() {
			Example(60.5)
		})
		Attribute("updated_at", String, "最終更新日時", func() {
			Format(FormatDateTime)
			Example("2023-08-23T10:00:00Z")
		})
	})
	View("default", func() {
		Attribute("total_assets")
		Attribute("cash_buying_power")
		Attribute("margin_buying_power")
		Attribute("withdrawal_possible_amount")
		Attribute("margin_rate")
		Attribute("updated_at")
	})
	Required("total_assets", "cash_buying_power", "margin_buying_power", "withdrawal_possible_amount", "margin_rate", "updated_at")
})

// CanEntryResult はエントリー可否の判断結果を定義します。
var CanEntryResult = ResultType("application/vnd.stock.balance.can.entry", func() {
	Description("エントリー可否の判断結果")
	Attributes(func() {
		Attribute("can_entry", Boolean, "エントリー可能かどうかのフラグ", func() {
			Example(true)
		})
		Attribute("buying_power", Float64, "エントリー判断時点の買付余力", func() {
			Example(1234567.89)
		})
	})
	View("default", func() {
		Attribute("can_entry")
		Attribute("buying_power")
	})
	Required("can_entry", "buying_power")
})

var _ = Service("balance", func() {
	Description("残高サービスは口座の残高情報を提供します。")

	Method("summary", func() {
		Description("口座の残高サマリーを取得します。")
		Result(BalanceSummary)
		HTTP(func() {
			GET("/balance/summary")
			Response(StatusOK)
		})
	})

	Method("canEntry", func() {
		Description("指定した銘柄にエントリー可能か判断します。")
		Payload(func() {
			Attribute("issue_code", String, "銘柄コード")
			Required("issue_code")
		})
		Result(CanEntryResult)
		HTTP(func() {
			GET("/balance/can_entry/{issue_code}")
			Response(StatusOK)
		})
	})
})

// NewOrderPayload は新規注文エンドポイントのペイロードを定義します。
var NewOrderPayload = Type("NewOrderPayload", func() {
	Description("新規注文を作成するためのペイロード")
	Attribute("symbol", String, "銘柄コード", func() {
		Example("9432")
	})
	Attribute("trade_type", String, "取引種別", func() {
		Enum("BUY", "SELL")
		Example("BUY")
	})
	Attribute("order_type", String, "注文種別", func() {
		Enum("MARKET", "LIMIT", "STOP", "STOP_LIMIT")
		Example("LIMIT")
	})
	Attribute("quantity", Int, "注文数量", func() {
		Example(100)
	})
	Attribute("price", Float64, "指値価格", func() {
		Example(3000.5)
	})
	Attribute("trigger_price", Float64, "逆指値トリガー価格", func() {
		Example(3100)
	})
	Attribute("time_in_force", String, "注文の有効期間", func() {
		Enum("DAY")
		Default("DAY")
	})
	Attribute("is_margin", Boolean, "信用取引かどうか", func() {
		Default(false)
	})
	Required("symbol", "trade_type", "order_type", "quantity")
})

// OrderResult は注文の結果型を定義します。
var OrderResult = ResultType("application/vnd.stock.order", func() {
	Description("注文操作の結果")
	Attributes(func() {
		Attribute("order_id", String, "一意の注文ID", func() {
			Example("20230824-123456")
		})
		Attribute("symbol", String, "銘柄コード", func() {
			Example("9432")
		})
		Attribute("trade_type", String, "取引種別")
		Attribute("order_type", String, "注文種別")
		Attribute("quantity", Int, "注文数量")
		Attribute("price", Float64, "指値価格")
		Attribute("trigger_price", Float64, "逆指値トリガー価格")
		Attribute("time_in_force", String, "注文の有効期間")
		Attribute("order_status", String, "注文ステータス")
		Attribute("is_margin", Boolean, "信用取引かどうか")
		Attribute("created_at", String, "作成タイムスタンプ", func() {
			Format(FormatDateTime)
		})
		Attribute("updated_at", String, "更新タイムスタンプ", func() {
			Format(FormatDateTime)
		})
	})
	View("default", func() {
		Attribute("order_id")
		Attribute("symbol")
		Attribute("trade_type")
		Attribute("order_type")
		Attribute("quantity")
		Attribute("price")
		Attribute("trigger_price")
		Attribute("time_in_force")
		Attribute("order_status")
		Attribute("is_margin")
		Attribute("created_at")
		Attribute("updated_at")
	})
	Required("order_id", "symbol", "trade_type", "order_type", "quantity", "order_status", "is_margin", "created_at", "updated_at")
})


var _ = Service("order", func() {
	Description("注文サービスは株式の注文操作を提供します。")

	Method("newOrder", func() {
		Description("新しい株式注文を作成します。")
		Payload(NewOrderPayload)
		Result(OrderResult)
		HTTP(func() {
			POST("/orders")
			Response(StatusCreated)
		})
	})
})
