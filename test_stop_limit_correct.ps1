# STOP_LIMIT注文の正しいパラメータでのテスト

Write-Host "=== STOP_LIMIT注文テスト（正しいパラメータ） ===" -ForegroundColor Green

# 買い注文: trigger_price以上になったら、より高い価格で買い注文
Write-Host "`n1. 買い注文（正しいパラメータ）:" -ForegroundColor Yellow
Write-Host "シナリオ: 970円以上になったら975円で買い注文" -ForegroundColor Cyan

$buyOrderBody = @{
    symbol = "3668"
    trade_type = "BUY"
    order_type = "STOP_LIMIT"
    quantity = 100
    price = 975                # 逆指値発動時の注文価格（trigger_priceより高く設定）
    trigger_price = 970        # 発動条件価格
    position_account_type = "CASH"
}

$buyOrderJson = $buyOrderBody | ConvertTo-Json -Compress
Write-Host "送信データ: $buyOrderJson" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $buyOrderJson -ContentType "application/json"
    Write-Host "✅ 買い注文成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ 買い注文エラー" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 3

# 売り注文: trigger_price以下になったら、より低い価格で売り注文
Write-Host "`n2. 売り注文（正しいパラメータ）:" -ForegroundColor Yellow
Write-Host "シナリオ: 980円以下になったら975円で売り注文" -ForegroundColor Cyan

$sellOrderBody = @{
    symbol = "3668"
    trade_type = "SELL"
    order_type = "STOP_LIMIT"
    quantity = 100
    price = 975                # 逆指値発動時の注文価格（trigger_priceより低く設定）
    trigger_price = 980        # 発動条件価格
    position_account_type = "CASH"
}

$sellOrderJson = $sellOrderBody | ConvertTo-Json -Compress
Write-Host "送信データ: $sellOrderJson" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/trade/orders" -Method POST -Body $sellOrderJson -ContentType "application/json"
    Write-Host "✅ 売り注文成功!" -ForegroundColor Green
    Write-Host "レスポンス: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ 売り注文エラー" -ForegroundColor Red
    Write-Host "エラー: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "詳細: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host "`n=== テスト完了 ===" -ForegroundColor Green