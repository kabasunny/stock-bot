# 修正されたGoaサービスでSTOP_LIMIT注文をテスト

Write-Host "=== Goa STOP_LIMIT注文テスト ===" -ForegroundColor Green

# 1. 逆指値注文（STOP）
Write-Host "1. 逆指値注文（STOP）テスト:" -ForegroundColor Yellow

$stopOrderJson = @"
{
    "symbol": "3632",
    "trade_type": "BUY",
    "order_type": "STOP",
    "quantity": 100,
    "trigger_price": 460,
    "position_account_type": "CASH"
}
"@

Write-Host $stopOrderJson

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $stopOrderJson -ContentType "application/json"
    Write-Host "✅ 成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー発生" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 2

# 2. 逆指値指値注文（STOP_LIMIT）
Write-Host "`n2. 逆指値指値注文（STOP_LIMIT）テスト:" -ForegroundColor Yellow

$stopLimitOrderJson = @"
{
    "symbol": "3668",
    "trade_type": "BUY",
    "order_type": "STOP_LIMIT",
    "quantity": 100,
    "price": 972,
    "trigger_price": 974,
    "position_account_type": "CASH"
}
"@

Write-Host $stopLimitOrderJson

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $stopLimitOrderJson -ContentType "application/json"
    Write-Host "✅ 成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー発生" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host "`n=== テスト完了 ===" -ForegroundColor Green