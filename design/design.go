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
