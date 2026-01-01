#!/bin/bash

# Goaサービステスト用スクリプト
BASE_URL="http://localhost:8080"

echo "=== Goa Service API Test ==="
echo "Base URL: $BASE_URL"
echo

# 1. ヘルスチェック
echo "1. Health Check"
curl -s -X GET "$BASE_URL/trade/health" | jq '.'
echo
echo

# 2. セッション情報取得
echo "2. Get Session"
curl -s -X GET "$BASE_URL/trade/session" | jq '.'
echo
echo

# 3. 残高取得
echo "3. Get Balance"
curl -s -X GET "$BASE_URL/trade/balance" | jq '.'
echo
echo

# 4. ポジション取得
echo "4. Get Positions"
curl -s -X GET "$BASE_URL/trade/positions" | jq '.'
echo
echo

# 5. 注文一覧取得
echo "5. Get Orders"
curl -s -X GET "$BASE_URL/trade/orders" | jq '.'
echo
echo

# 6. 注文発行（テスト用）
echo "6. Place Order (Test)"
curl -s -X POST "$BASE_URL/trade/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "7203",
    "trade_type": "BUY",
    "order_type": "MARKET",
    "quantity": 100,
    "price": 0,
    "position_account_type": "CASH"
  }' | jq '.'
echo
echo

# 7. 価格履歴取得
echo "7. Get Price History"
curl -s -X GET "$BASE_URL/trade/price-history/7203?days=5" | jq '.'
echo
echo

# 8. 銘柄妥当性チェック
echo "8. Validate Symbol"
curl -s -X GET "$BASE_URL/trade/symbols/7203/validate" | jq '.'
echo
echo

echo "=== Test Complete ==="