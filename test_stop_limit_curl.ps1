# STOP_LIMIT注文のテスト（curl使用）

Write-Host "=== STOP_LIMIT注文テスト（修正版） ===" -ForegroundColor Green

# 1. STOP注文（参考）
Write-Host "`n1. STOP注文テスト:" -ForegroundColor Yellow
$stopJson = '{"symbol":"3632","trade_type":"BUY","order_type":"STOP","quantity":100,"trigger_price":460,"position_account_type":"CASH"}'

try {
    $response = curl -X POST "http://localhost:8080/trade/orders" -H "Content-Type: application/json" -d $stopJson
    Write-Host "STOP注文レスポンス:" -ForegroundColor Green
    Write-Host $response
} catch {
    Write-Host "STOP注文エラー: $($_.Exception.Message)" -ForegroundColor Red
}

Start-Sleep -Seconds 2

# 2. STOP_LIMIT注文（修正版）
Write-Host "`n2. STOP_LIMIT注文テスト（修正版）:" -ForegroundColor Yellow
$stopLimitJson = '{"symbol":"3668","trade_type":"BUY","order_type":"STOP_LIMIT","quantity":100,"price":972,"trigger_price":974,"position_account_type":"CASH"}'

Write-Host "送信データ: $stopLimitJson" -ForegroundColor Cyan

try {
    $response = curl -X POST "http://localhost:8080/trade/orders" -H "Content-Type: application/json" -d $stopLimitJson
    Write-Host "STOP_LIMIT注文レスポンス:" -ForegroundColor Green
    Write-Host $response
} catch {
    Write-Host "STOP_LIMIT注文エラー: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== テスト完了 ===" -ForegroundColor Green