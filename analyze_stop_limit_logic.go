package main

import (
	"fmt"
	"stock-bot/domain/model"
)

func main() {
	fmt.Println("=== STOP_LIMIT注文の正しい理解 ===")

	// STOP_LIMIT注文の仕様分析
	fmt.Println("\n立花証券のSTOP_LIMIT注文パラメータ:")
	fmt.Println("- OrderPrice: 逆指値が発動した時の注文価格")
	fmt.Println("- GyakusasiZyouken: 逆指値の発動条件価格")
	fmt.Println("- GyakusasiPrice: 逆指値が発動した時の注文価格（OrderPriceと同じ）")
	fmt.Println("- GyakusasiOrderType: \"2\"（通常+逆指値）")

	fmt.Println("\nAPIスキーマとの対応:")
	fmt.Println("- price: 逆指値が発動した時の注文価格 → OrderPrice, GyakusasiPrice")
	fmt.Println("- trigger_price: 逆指値の発動条件価格 → GyakusasiZyouken")

	// テストケース1: 買い注文
	fmt.Println("\n=== 買い注文のケース ===")
	fmt.Println("シナリオ: 現在価格970円の株を、975円以上になったら980円で買い注文")

	buyPrice := 980.0
	buyTriggerPrice := 975.0
	buyTradeType := model.TradeTypeBuy

	fmt.Printf("入力パラメータ:\n")
	fmt.Printf("  trade_type: %s\n", buyTradeType)
	fmt.Printf("  price: %.0f (逆指値発動時の注文価格)\n", buyPrice)
	fmt.Printf("  trigger_price: %.0f (発動条件: この価格以上で発動)\n", buyTriggerPrice)

	fmt.Printf("期待される立花証券パラメータ:\n")
	fmt.Printf("  OrderPrice: %.0f\n", buyPrice)
	fmt.Printf("  GyakusasiZyouken: %.0f\n", buyTriggerPrice)
	fmt.Printf("  GyakusasiPrice: %.0f\n", buyPrice)
	fmt.Printf("  BaibaiKubun: 3 (買い)\n")

	// テストケース2: 売り注文
	fmt.Println("\n=== 売り注文のケース ===")
	fmt.Println("シナリオ: 現在価格1000円の株を、990円以下になったら985円で売り注文")

	sellPrice := 985.0
	sellTriggerPrice := 990.0
	sellTradeType := model.TradeTypeSell

	fmt.Printf("入力パラメータ:\n")
	fmt.Printf("  trade_type: %s\n", sellTradeType)
	fmt.Printf("  price: %.0f (逆指値発動時の注文価格)\n", sellPrice)
	fmt.Printf("  trigger_price: %.0f (発動条件: この価格以下で発動)\n", sellTriggerPrice)

	fmt.Printf("期待される立花証券パラメータ:\n")
	fmt.Printf("  OrderPrice: %.0f\n", sellPrice)
	fmt.Printf("  GyakusasiZyouken: %.0f\n", sellTriggerPrice)
	fmt.Printf("  GyakusasiPrice: %.0f\n", sellPrice)
	fmt.Printf("  BaibaiKubun: 1 (売り)\n")

	fmt.Println("\n=== 重要なポイント ===")
	fmt.Println("1. OrderPriceとGyakusasiPriceは同じ値（逆指値発動時の注文価格）")
	fmt.Println("2. GyakusasiZyoukenは発動条件価格（trigger_price）")
	fmt.Println("3. 売買方向は関係なく、パラメータの意味は同じ")
	fmt.Println("4. 発動条件の方向（以上/以下）は立花証券側で売買区分から自動判定")

	fmt.Println("\n=== 現在の変換ロジックの修正点 ===")
	fmt.Println("✅ formatOrderPrice: price をそのまま使用（修正済み）")
	fmt.Println("✅ formatGyakusasiPrice: price をそのまま使用（既に正しい）")
	fmt.Println("✅ formatTriggerPrice: trigger_price をそのまま使用（既に正しい）")
	fmt.Println("✅ convertGyakusasiOrderType: \"2\" を返す（既に正しい）")
}
