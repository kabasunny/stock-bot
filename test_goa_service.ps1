# Goaサービステスト用PowerShellスクリプト
$BASE_URL = "http://localhost:8080"

Write-Host "=== Goa Service API Test ===" -ForegroundColor Green
Write-Host "Base URL: $BASE_URL"
Write-Host ""

# 1. ヘルスチェック
Write-Host "1. Health Check" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/health" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 2. セッション情報取得
Write-Host "2. Get Session" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/session" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 3. 残高取得
Write-Host "3. Get Balance" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/balance" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 4. ポジション取得
Write-Host "4. Get Positions" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/positions" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 5. 注文一覧取得
Write-Host "5. Get Orders" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/orders" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 6. 注文発行（テスト用）
Write-Host "6. Place Order (Test)" -ForegroundColor Yellow
try {
    $orderData = @{
        symbol = "7203"
        trade_type = "BUY"
        order_type = "MARKET"
        quantity = 100
        price = 0
        position_account_type = "CASH"
    }
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/orders" -Method POST -Body ($orderData | ConvertTo-Json) -ContentType "application/json"
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 7. 価格履歴取得
Write-Host "7. Get Price History" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/price-history/7203?days=5" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 8. 銘柄妥当性チェック
Write-Host "8. Validate Symbol" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/trade/symbols/7203/validate" -Method GET
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Green