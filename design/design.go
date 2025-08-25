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
