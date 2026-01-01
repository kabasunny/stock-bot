$body = @{
    symbol = "7203"
    trade_type = "BUY"
    order_type = "MARKET"
    quantity = 100
    position_account_type = "CASH"
} | ConvertTo-Json

Write-Host "リクエスト: $body"

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $body -ContentType "application/json"
    Write-Host "成功: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
}