package design

import (
	. "goa.design/goa/v3/dsl"
)

// API全体の定義
var _ = API("stockbot", func() {
	Title("Stock Bot Service")
	Description("Service for placing and managing stock orders")
	Server("stockbot", func() {
		Host("localhost", func() {
			// ポートは後ほど設定ファイルから読み込む
			URI("http://localhost:8080")
		})
	})
})

// 注文サービス(Order)の定義
var _ = Service("order", func() {
	Description("The order service handles placing stock orders.")

	// POST /order
	Method("create", func() {
		Description("Create a new stock order.")

		// リクエストのペイロード(JSONボディ)
		Payload(func() {
			Attribute("symbol", String, "銘柄コード (例: 7203)")
			Attribute("trade_type", String, "売買区分 (BUY/SELL)", func() {
				Enum("BUY", "SELL")
			})
			Attribute("order_type", String, "注文種別 (MARKET/LIMITなど)", func() {
				Enum("MARKET", "LIMIT", "STOP", "STOP_LIMIT")
			})
			Attribute("quantity", UInt64, "発注数量")
			Attribute("price", Float64, "発注価格 (LIMIT注文の場合)", func() {
				Default(0)
			})
			Attribute("position_account_type", String, "ポジションの口座区分 (CASH/MARGIN_NEW/MARGIN_REPAY)", func() {
				Enum("CASH", "MARGIN_NEW", "MARGIN_REPAY")
				Default("CASH")
			})
			Required("symbol", "trade_type", "order_type", "quantity")
		})

		// レスポンス
		Result(func() {
			Description("ID of the created order")
			Attribute("order_id", String, "受付済み注文ID")
			Required("order_id")
		})

		// HTTPプロトコルとのマッピング
		HTTP(func() {
			POST("/order")
			Response(StatusCreated)
		})
	})
})

// Goa Type for Balance Summary
var BalanceResult = ResultType("application/vnd.stockbot.balance", func() {
	Description("A summary of the account balance.")
	Attribute("available_cash_for_stock", Float64, "現物株式買付可能額")
	Attribute("available_margin_for_new_position", Float64, "信用新規建可能額")
	Attribute("margin_maintenance_rate", Float64, "委託保証金率(%)")
	Attribute("withdrawable_cash", Float64, "出金可能額")
	Attribute("has_margin_call", Boolean, "追証発生フラグ (1:発生, 0:未発生)")
	Required(
		"available_cash_for_stock",
		"available_margin_for_new_position",
		"margin_maintenance_rate",
		"withdrawable_cash",
		"has_margin_call",
	)
})

// 残高サービス(Balance)の定義
var _ = Service("balance", func() {
	Description("The balance service provides account balance information.")

	// GET /balance
	Method("get", func() {
		Description("Get the account balance summary.")
		Payload(Empty) // No request body
		Result(BalanceResult)

		HTTP(func() {
			GET("/balance")
			Response(StatusOK)
		})
	})
})

// Goa Type for Stock Price Result
var PriceResult = ResultType("application/vnd.stockbot.price", func() {
	Description("The current price information for a stock.")
	Attribute("symbol", String, "銘柄コード")
	Attribute("price", Float64, "現在値")
	Attribute("timestamp", String, "価格取得日時 (RFC3339)")
	Required("symbol", "price", "timestamp")
})

// 価格サービス(Price)の定義
var _ = Service("price", func() {
	Description("The price service provides current stock price information.")

	// GET /price/{symbol}
	Method("get", func() {
		Description("Get the current price for a specified stock symbol.")
		Payload(func() {
			Attribute("symbol", String, "Stock symbol to look up")
			Required("symbol")
		})
		Result(PriceResult)

		HTTP(func() {
			GET("/price/{symbol}")
			Response(StatusOK)
		})
	})

	// GET /price/{symbol}/history
	Method("get_history", func() {
		Description("Get historical price data for a specified stock symbol.")
		Payload(func() {
			Attribute("symbol", String, "Stock symbol to look up")
			Attribute("days", UInt, "Number of historical days to retrieve (optional)", func() {
				Default(0) // 0 means all available history
			})
			Required("symbol")
		})
		Result(HistoricalPriceResult)

		HTTP(func() {
			GET("/price/{symbol}/history")
			Param("days")
			Response(StatusOK)
		})
	})
})

// Goa Type for a single Historical Price data point
var HistoricalPriceItem = Type("HistoricalPriceItem", func() {
	Description("A single historical price data point.")
	Attribute("date", String, "日付 (YYYY-MM-DD)")
	Attribute("open", Float64, "始値")
	Attribute("high", Float64, "高値")
	Attribute("low", Float64, "安値")
	Attribute("close", Float64, "終値")
	Attribute("volume", UInt64, "出来高")
	Required("date", "open", "high", "low", "close")
})

// Goa Type for Historical Price Result (collection)
var HistoricalPriceResult = ResultType("application/vnd.stockbot.historical-price", func() {
	Description("Historical price information for a stock.")
	Attribute("symbol", String, "銘柄コード")
	Attribute("history", ArrayOf(HistoricalPriceItem), "過去の価格データ")
	Required("symbol", "history")
})

// Goa Type for a single Position
var PositionResult = Type("PositionResult", func() {
	Description("A single trading position.")
	Attribute("symbol", String, "銘柄コード")
	Attribute("position_type", String, "ポジション種別 (CASH, MARGIN_LONG, MARGIN_SHORT)", func() {
		Enum("CASH", "MARGIN_LONG", "MARGIN_SHORT")
	})
	Attribute("quantity", Float64, "保有数量")
	Attribute("average_cost", Float64, "平均取得単価")
	Attribute("current_price", Float64, "現在値")
	Attribute("unrealized_pl", Float64, "評価損益")
	Attribute("unrealized_pl_rate", Float64, "評価損益率(%)")
	Attribute("opened_date", String, "建日 (信用取引の場合 YYYYMMDD)")

	Required("symbol", "position_type", "quantity", "average_cost")
})

// Goa Type for a collection of Positions
var PositionCollection = ResultType("application/vnd.stockbot.position-collection", func() {
	Description("A collection of trading positions.")
	Attribute("positions", ArrayOf(PositionResult), "保有ポジションのリスト")
	Required("positions")
})

// 保有ポジションサービス(Position)の定義
var _ = Service("position", func() {
	Description("The position service provides information about current holdings.")

	// GET /positions
	Method("list", func() {
		Description("List current positions.")
		Payload(func() {
			Attribute("type", String, "取得するポジション種別 (all, cash, margin)", func() {
				Enum("all", "cash", "margin")
				Default("all")
			})
		})
		Result(PositionCollection)

		HTTP(func() {
			GET("/positions")
			Param("type") // Map the 'type' attribute to a query parameter
			Response(StatusOK)
		})
	})
})

// Goa Type for Stock Master Data (simplified)
var StockMasterResult = ResultType("application/vnd.stockbot.stock-master", func() {
	Description("Basic master data for a single stock.")
	Attribute("symbol", String, "銘柄コード")
	Attribute("name", String, "銘柄名")
	Attribute("name_kana", String, "銘柄名（カナ）")
	Attribute("market", String, "優先市場")
	Attribute("industry_code", String, "業種コード")
	Attribute("industry_name", String, "業種コード名")

	Required("symbol", "name", "market") // Minimal required fields
})

// マスタデータサービス(Master)の定義
var _ = Service("master", func() {
	Description("The master service provides master data.")

	// GET /master/stocks/{symbol}
	Method("get_stock", func() { // Renamed method
		Description("Get basic master data for a single stock.")
		Payload(func() {
			Attribute("symbol", String, "Stock symbol to look up")
			Required("symbol")
		})
		Result(StockMasterResult) // New result type

		HTTP(func() {
			GET("/master/stocks/{symbol}")
			Response(StatusOK)
		})
	})

	// POST /master/update
	Method("update", func() {
		Description("Trigger a manual update of the master data.")
		Payload(Empty)
		Result(Empty) // 成功したかどうかはHTTPステータスで判断

		HTTP(func() {
			POST("/master/update")
			Response(StatusAccepted) // 処理を受け付けたことを示す
		})
	})
})

// TradeServiceの定義
var _ = Service("trade", func() {
	Description("The trade service provides unified trading operations.")

	// GET /trade/session
	Method("get_session", func() {
		Description("Get current API session information.")
		Payload(Empty)
		Result(func() {
			Attribute("session_id", String, "セッションID")
			Attribute("user_id", String, "ユーザーID")
			Attribute("login_time", String, "ログイン時刻 (RFC3339)")
			Required("session_id", "user_id", "login_time")
		})

		HTTP(func() {
			GET("/trade/session")
			Response(StatusOK)
		})
	})

	// GET /trade/positions
	Method("get_positions", func() {
		Description("Get current trading positions.")
		Payload(Empty)
		Result(func() {
			Attribute("positions", ArrayOf(TradePositionResult), "保有ポジションのリスト")
			Required("positions")
		})

		HTTP(func() {
			GET("/trade/positions")
			Response(StatusOK)
		})
	})

	// GET /trade/orders
	Method("get_orders", func() {
		Description("Get current orders.")
		Payload(Empty)
		Result(func() {
			Attribute("orders", ArrayOf(TradeOrderResult), "注文のリスト")
			Required("orders")
		})

		HTTP(func() {
			GET("/trade/orders")
			Response(StatusOK)
		})
	})

	// GET /trade/balance
	Method("get_balance", func() {
		Description("Get account balance.")
		Payload(Empty)
		Result(TradeBalanceResult)

		HTTP(func() {
			GET("/trade/balance")
			Response(StatusOK)
		})
	})

	// GET /trade/price-history/{symbol}
	Method("get_price_history", func() {
		Description("Get price history for a symbol.")
		Payload(func() {
			Attribute("symbol", String, "銘柄コード")
			Attribute("days", UInt, "取得日数", func() {
				Default(30)
			})
			Required("symbol")
		})
		Result(func() {
			Attribute("symbol", String, "銘柄コード")
			Attribute("history", ArrayOf(TradePriceHistoryItem), "価格履歴")
			Required("symbol", "history")
		})

		HTTP(func() {
			GET("/trade/price-history/{symbol}")
			Param("days")
			Response(StatusOK)
		})
	})

	// POST /trade/orders
	Method("place_order", func() {
		Description("Place a new order.")
		Payload(func() {
			Attribute("symbol", String, "銘柄コード")
			Attribute("trade_type", String, "売買区分", func() {
				Enum("BUY", "SELL")
			})
			Attribute("order_type", String, "注文種別", func() {
				Enum("MARKET", "LIMIT", "STOP")
			})
			Attribute("quantity", UInt, "数量")
			Attribute("price", Float64, "価格 (指値の場合)", func() {
				Default(0)
			})
			Attribute("trigger_price", Float64, "トリガー価格 (逆指値の場合)", func() {
				Default(0)
			})
			Attribute("position_account_type", String, "口座区分", func() {
				Enum("CASH", "MARGIN_NEW", "MARGIN_REPAY")
				Default("CASH")
			})
			Required("symbol", "trade_type", "order_type", "quantity")
		})
		Result(TradeOrderResult)

		HTTP(func() {
			POST("/trade/orders")
			Response(StatusCreated)
		})
	})

	// DELETE /trade/orders/{order_id}
	Method("cancel_order", func() {
		Description("Cancel an existing order.")
		Payload(func() {
			Attribute("order_id", String, "注文ID")
			Required("order_id")
		})
		Result(Empty)

		HTTP(func() {
			DELETE("/trade/orders/{order_id}")
			Response(StatusNoContent)
		})
	})

	// PUT /trade/orders/{order_id}
	Method("correct_order", func() {
		Description("Correct an existing order.")
		Payload(func() {
			Attribute("order_id", String, "注文ID")
			Attribute("price", Float64, "新しい価格")
			Attribute("quantity", UInt, "新しい数量")
			Required("order_id")
		})
		Result(TradeOrderResult)

		HTTP(func() {
			PUT("/trade/orders/{order_id}")
			Response(StatusOK)
		})
	})

	// DELETE /trade/orders
	Method("cancel_all_orders", func() {
		Description("Cancel all pending orders.")
		Payload(Empty)
		Result(func() {
			Attribute("cancelled_count", UInt, "キャンセルされた注文数")
			Required("cancelled_count")
		})

		HTTP(func() {
			DELETE("/trade/orders")
			Response(StatusOK)
		})
	})

	// GET /trade/symbols/{symbol}/validate
	Method("validate_symbol", func() {
		Description("Validate if a symbol is tradable and get trading information.")
		Payload(func() {
			Attribute("symbol", String, "銘柄コード")
			Required("symbol")
		})
		Result(func() {
			Attribute("valid", Boolean, "取引可能かどうか")
			Attribute("symbol", String, "銘柄コード")
			Attribute("name", String, "銘柄名")
			Attribute("trading_unit", UInt, "売買単位")
			Attribute("market", String, "市場")
			Required("valid", "symbol")
		})

		HTTP(func() {
			GET("/trade/symbols/{symbol}/validate")
			Response(StatusOK)
		})
	})

	// GET /trade/orders/history
	Method("get_order_history", func() {
		Description("Get order history with optional filtering.")
		Payload(func() {
			Attribute("status", String, "注文状態でフィルタ (NEW/FILLED/CANCELLED)", func() {
				Enum("NEW", "PARTIALLY_FILLED", "FILLED", "CANCELLED", "REJECTED")
			})
			Attribute("symbol", String, "銘柄コードでフィルタ")
			Attribute("limit", UInt, "取得件数制限", func() {
				Default(100)
			})
		})
		Result(func() {
			Attribute("orders", ArrayOf(TradeOrderHistoryResult), "注文履歴")
			Required("orders")
		})

		HTTP(func() {
			GET("/trade/orders/history")
			Param("status")
			Param("symbol")
			Param("limit")
			Response(StatusOK)
		})
	})

	// GET /trade/health
	Method("health_check", func() {
		Description("Check service health status.")
		Payload(Empty)
		Result(func() {
			Attribute("status", String, "サービス状態", func() {
				Enum("healthy", "degraded", "unhealthy")
			})
			Attribute("timestamp", String, "チェック時刻 (RFC3339)")
			Attribute("session_valid", Boolean, "セッション有効性")
			Attribute("database_connected", Boolean, "データベース接続状態")
			Attribute("websocket_connected", Boolean, "WebSocket接続状態")
			Required("status", "timestamp")
		})

		HTTP(func() {
			GET("/trade/health")
			Response(StatusOK)
		})
	})
})

// TradeService用の型定義
var TradePositionResult = Type("TradePositionResult", func() {
	Description("Trading position information.")
	Attribute("symbol", String, "銘柄コード")
	Attribute("position_type", String, "ポジション種別", func() {
		Enum("LONG", "SHORT")
	})
	Attribute("position_account_type", String, "口座区分", func() {
		Enum("CASH", "MARGIN_NEW", "MARGIN_REPAY")
	})
	Attribute("average_price", Float64, "平均取得価格")
	Attribute("quantity", UInt, "数量")
	Required("symbol", "position_type", "position_account_type", "average_price", "quantity")
})

var TradeOrderResult = Type("TradeOrderResult", func() {
	Description("Trading order information.")
	Attribute("order_id", String, "注文ID")
	Attribute("symbol", String, "銘柄コード")
	Attribute("trade_type", String, "売買区分", func() {
		Enum("BUY", "SELL")
	})
	Attribute("order_type", String, "注文種別", func() {
		Enum("MARKET", "LIMIT", "STOP")
	})
	Attribute("quantity", UInt, "数量")
	Attribute("price", Float64, "価格")
	Attribute("order_status", String, "注文状態", func() {
		Enum("NEW", "PARTIALLY_FILLED", "FILLED", "CANCELLED", "REJECTED")
	})
	Attribute("position_account_type", String, "口座区分", func() {
		Enum("CASH", "MARGIN_NEW", "MARGIN_REPAY")
	})
	Required("order_id", "symbol", "trade_type", "order_type", "quantity", "price", "order_status")
})

var TradeBalanceResult = Type("TradeBalanceResult", func() {
	Description("Account balance information.")
	Attribute("cash", Float64, "現金残高")
	Attribute("buying_power", Float64, "買付余力")
	Required("cash", "buying_power")
})

var TradePriceHistoryItem = Type("TradePriceHistoryItem", func() {
	Description("Historical price data point.")
	Attribute("date", String, "日付 (RFC3339)")
	Attribute("open", Float64, "始値")
	Attribute("high", Float64, "高値")
	Attribute("low", Float64, "安値")
	Attribute("close", Float64, "終値")
	Attribute("volume", UInt64, "出来高")
	Required("date", "open", "high", "low", "close", "volume")
})

var TradeOrderHistoryResult = Type("TradeOrderHistoryResult", func() {
	Description("Order history information with execution details.")
	Attribute("order_id", String, "注文ID")
	Attribute("symbol", String, "銘柄コード")
	Attribute("trade_type", String, "売買区分", func() {
		Enum("BUY", "SELL")
	})
	Attribute("order_type", String, "注文種別", func() {
		Enum("MARKET", "LIMIT", "STOP")
	})
	Attribute("quantity", UInt, "数量")
	Attribute("price", Float64, "価格")
	Attribute("order_status", String, "注文状態", func() {
		Enum("NEW", "PARTIALLY_FILLED", "FILLED", "CANCELLED", "REJECTED")
	})
	Attribute("position_account_type", String, "口座区分", func() {
		Enum("CASH", "MARGIN_NEW", "MARGIN_REPAY")
	})
	Attribute("created_at", String, "注文日時 (RFC3339)")
	Attribute("updated_at", String, "更新日時 (RFC3339)")
	Attribute("executions", ArrayOf(TradeExecutionResult), "約定履歴")
	Required("order_id", "symbol", "trade_type", "order_type", "quantity", "price", "order_status", "created_at")
})

var TradeExecutionResult = Type("TradeExecutionResult", func() {
	Description("Execution information.")
	Attribute("execution_id", String, "約定ID")
	Attribute("executed_quantity", UInt, "約定数量")
	Attribute("executed_price", Float64, "約定価格")
	Attribute("executed_at", String, "約定日時 (RFC3339)")
	Required("execution_id", "executed_quantity", "executed_price", "executed_at")
})
