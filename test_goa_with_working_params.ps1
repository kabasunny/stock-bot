# 立花クライアントで成功したパラメータでGoaサービスをテスト

Write-Host "=== 立花クライアント成功パラメータでGoaサービステスト ===" -ForegroundColor Green

$body = @{
    symbol = "6658"
    trade_type = "BUY"
    order_type = "MARKET"
    quantity = 100
    position_account_type = "CASH"
} | ConvertTo-Json

Write-Host "リクエスト内容:" -ForegroundColor Yellow
Write-Host $body

Write-Host "Goaサービスにリクエスト送信中..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $body -ContentType "application/json"
    Write-Host "成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "エラー発生" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host "テスト完了" -ForegroundColor Green