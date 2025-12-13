# API cURL Command Examples

This file contains example `curl` commands for interacting with the Stock Bot API.

## Pre-requisites

- The application server must be running (`go run ./cmd/myapp/main.go`).
- The server is expected to be on `http://localhost:8080`.

# 注意点
- Windowsでは JSON文字列の中のダブルクォート " をエスケープする必要があります。
-  → \"key\":\"value\" のように書く。
- PowerShell では curl は Invoke-WebRequest のエイリアスですが、通常の curl.exe がインストールされていればそのまま使えます。
- 長い JSON を扱う場合は ファイルに保存して -d @file.json とする方が安全で読みやすいです。

---

## Balance Service

### Get Account Balance Summary

Retrieves a summary of the account's balance, including available cash and margin information.

```sh
# Get account balance summary
curl -i -X GET http://localhost:8080/balance
```

---

## Order Service

### Create a Market Buy Order

Places a "market" buy order for a specified quantity of a stock.

```sh
# Create a new MARKET BUY order for 100 shares of symbol 7203 (Toyota)
curl -i -X POST -H "Content-Type: application/json" -d "{\"symbol\":\"7203\",\"trade_type\":\"BUY\",\"order_type\":\"MARKET\",\"quantity\":100}" http://localhost:8080/order
```

### Create a Limit Sell Order

Places a "limit" sell order for a specified quantity of a stock at a specific price or better.

```sh
# Create a new LIMIT SELL order for 50 shares of symbol 6758 (Sony) at a price of 13000
curl -i -X POST -H "Content-Type: application/json" -d "{\"symbol\":\"7203\",\"trade_type\":\"SELL\",\"order_type\":\"LIMIT\",\"quantity\":100,\"price\":3500.0}" http://localhost:8080/order
```

### Create a Market Buy Order (Margin)

Places a market buy order using margin.

```sh
# Create a new MARKET BUY order for 100 shares of symbol 9984 (SoftBank) using margin
curl -i -X POST -H "Content-Type: application/json" -d "{\"symbol\":\"9984\",\"trade_type\":\"BUY\",\"order_type\":\"MARKET\",\"quantity\":100,\"is_margin\":true}" http://localhost:8080/order
```

---

## Position Service

### List Current Positions

Retrieves a list of all currently held positions (cash and margin).

```sh
# Get all positions
curl -i -X GET http://localhost:8080/positions

# Get only cash positions
curl -i -X GET "http://localhost:8080/positions?type=cash"

# Get only margin positions
curl -i -X GET "http://localhost:8080/positions?type=margin"
```

---

## Master Service

### Get Stock Detail

Retrieves detailed master data for a specific stock symbol.

```sh
# Get details for Toyota (symbol 7203)
curl -i -X GET http://localhost:8080/master/stocks/7203

# Get details for Sony (symbol 6758)
curl -i -X GET http://localhost:8080/master/stocks/6758
```

