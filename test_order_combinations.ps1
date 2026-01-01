# 注文タイプ × 口座区分の組み合わせテスト

$baseUrl = "http://localhost:8080"
$symbol = "7203"  # トヨタ自動車

Write-Host "=== 注文タイプ × 口座区分 組み合わせテスト ===" -ForegroundColor Green

# 1. 成行注文 × 現物
Write-Host "`n1. 成行注文 × 現物" -ForegroundColor Yellow
$body1 = @{
    symbol = $symbol
    trade_type = "BUY"
    order_type = "MARKET"
    quantity = 100
    position_account_type = "CASH"
} | ConvertTo-Json

try {
    $response1 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body1 -ContentType "application/json"
    Write-Host "✅ 成功: $($response1 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 2. 成行注文 × 信用新規
Write-Host "`n2. 成行注文 × 信用新規" -ForegroundColor Yellow
$body2 = @{
    symbol = $symbol
    trade_type = "BUY"
    order_type = "MARKET"
    quantity = 100
    position_account_type = "MARGIN_NEW"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body2 -ContentType "application/json"
    Write-Host "✅ 成功: $($response2 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. 指値注文 × 現物
Write-Host "`n3. 指値注文 × 現物" -ForegroundColor Yellow
$body3 = @{
    symbol = $symbol
    trade_type = "BUY"
    order_type = "LIMIT"
    quantity = 100
    price = 2800.0
    position_account_type = "CASH"
} | ConvertTo-Json

try {
    $response3 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body3 -ContentType "application/json"
    Write-Host "✅ 成功: $($response3 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. 指値注文 × 信用新規
Write-Host "`n4. 指値注文 × 信用新規" -ForegroundColor Yellow
$body4 = @{
    symbol = $symbol
    trade_type = "SELL"
    order_type = "LIMIT"
    quantity = 100
    price = 2900.0
    position_account_type = "MARGIN_NEW"
} | ConvertTo-Json

try {
    $response4 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body4 -ContentType "application/json"
    Write-Host "✅ 成功: $($response4 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. 逆指値注文 × 現物
Write-Host "`n5. 逆指値注文 × 現物" -ForegroundColor Yellow
$body5 = @{
    symbol = $symbol
    trade_type = "SELL"
    order_type = "STOP"
    quantity = 100
    trigger_price = 2700.0
    position_account_type = "CASH"
} | ConvertTo-Json

try {
    $response5 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body5 -ContentType "application/json"
    Write-Host "✅ 成功: $($response5 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. 逆指値注文 × 信用返済
Write-Host "`n6. 逆指値注文 × 信用返済" -ForegroundColor Yellow
$body6 = @{
    symbol = $symbol
    trade_type = "BUY"
    order_type = "STOP"
    quantity = 100
    trigger_price = 3000.0
    position_account_type = "MARGIN_REPAY"
} | ConvertTo-Json

try {
    $response6 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body6 -ContentType "application/json"
    Write-Host "✅ 成功: $($response6 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 7. 逆指値指値注文 × 現物
Write-Host "`n7. 逆指値指値注文 × 現物" -ForegroundColor Yellow
$body7 = @{
    symbol = $symbol
    trade_type = "BUY"
    order_type = "STOP_LIMIT"
    quantity = 100
    price = 2850.0
    trigger_price = 2800.0
    position_account_type = "CASH"
} | ConvertTo-Json

try {
    $response7 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body7 -ContentType "application/json"
    Write-Host "✅ 成功: $($response7 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

# 8. 逆指値指値注文 × 信用新規
Write-Host "`n8. 逆指値指値注文 × 信用新規" -ForegroundColor Yellow
$body8 = @{
    symbol = $symbol
    trade_type = "SELL"
    order_type = "STOP_LIMIT"
    quantity = 100
    price = 2750.0
    trigger_price = 2800.0
    position_account_type = "MARGIN_NEW"
} | ConvertTo-Json

try {
    $response8 = Invoke-RestMethod -Uri "$baseUrl/trade/orders" -Method POST -Body $body8 -ContentType "application/json"
    Write-Host "✅ 成功: $($response8 | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Host "❌ エラー: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== テスト完了 ===" -ForegroundColor Green