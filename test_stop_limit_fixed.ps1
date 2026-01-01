# STOP_LIMIT注文のテスト（修正版）

Write-Host "=== STOP_LIMIT注文テスト（修正版） ===" -ForegroundColor Green

# 1. STOP注文（参考）
Write-Host "`n1. STOP注文テスト:" -ForegroundColor Yellow

$stopOrderBody = @{
    symbol = "3632"
    trade_type = "BUY"
    order_type = "STOP"
    quantity = 100
    trigger_price = 460
    position_account_type = "CASH"
}

$stopOrderJson = $stopOrderBody | ConvertTo-Json -Compress
Write-Host "送信データ: $stopOrderJson" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $stopOrderJson -ContentType "application/json"
    Write-Host "✅ STOP注文成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ STOP注文エラー" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 3

# 2. STOP_LIMIT注文（修正版）
Write-Host "`n2. STOP_LIMIT注文テスト（修正版）:" -ForegroundColor Yellow

$stopLimitOrderBody = @{
    symbol = "3668"
    trade_type = "BUY"
    order_type = "STOP_LIMIT"
    quantity = 100
    price = 972
    trigger_price = 974
    position_account_type = "CASH"
}

$stopLimitOrderJson = $stopLimitOrderBody | ConvertTo-Json -Compress
Write-Host "送信データ: $stopLimitOrderJson" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $stopLimitOrderJson -ContentType "application/json"
    Write-Host "✅ STOP_LIMIT注文成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ STOP_LIMIT注文エラー" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host "`n=== テスト完了 ===" -ForegroundColor Green