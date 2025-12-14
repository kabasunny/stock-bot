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
            Attribute("is_margin", Boolean, "信用取引かどうか", func() {
                Default(false)
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
