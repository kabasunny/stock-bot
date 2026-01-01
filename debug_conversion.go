package main

import (
	"fmt"
	"stock-bot/domain/model"
)

// 変換関数のテスト用コピー（goa_trade_service.goから）

// formatOrderPrice は通常注文価格を立花証券API用の文字列に変換
func formatOrderPrice(price float64, orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeMarket, model.OrderTypeStop:
		return "*" // 成行・逆指値成行の場合は*
	case model.OrderTypeLimit:
		return fmt.Sprintf("%.0f", price) // 指値の場合は価格
	case model.OrderTypeStopLimit:
		// STOP_LIMIT注文では、OrderPriceは逆指値時の注文価格（price）を使用
		return fmt.Sprintf("%.0f", price)
	default:
		return "*"
	}
}

// formatGyakusasiPrice は逆指値価格を立花証券API用の文字列に変換
func formatGyakusasiPrice(price float64, orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeStop:
		return "0" // 逆指値成行の場合は0
	case model.OrderTypeStopLimit:
		return fmt.Sprintf("%.0f", price) // 逆指値指値の場合は注文価格
	default:
		return "*" // 通常注文は*
	}
}

// formatTriggerPrice は逆指値価格を立花証券API用の文字列に変換
func formatTriggerPrice(triggerPrice *float64, orderType model.OrderType) string {
	if triggerPrice != nil && (orderType == model.OrderTypeStop || orderType == model.OrderTypeStopLimit) {
		return fmt.Sprintf("%.0f", *triggerPrice)
	}
	return "" // 逆指値以外は空文字
}

// convertGyakusasiOrderType は逆指値注文タイプを変換
func convertGyakusasiOrderType(orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeStop:
		return "1" // 逆指値
	case model.OrderTypeStopLimit:
		return "2" // 通常＋逆指値
	default:
		return "0" // 通常注文
	}
}

func main() {
	fmt.Println("=== 変換関数デバッグテスト ===")

	// テストケース1: STOP_LIMIT注文（修正版）
	fmt.Println("\n1. STOP_LIMIT注文の変換テスト（修正版）:")
	price := 972.0
	triggerPrice := 974.0
	orderType := model.OrderTypeStopLimit

	fmt.Printf("入力パラメータ:\n")
	fmt.Printf("  price: %.0f (逆指値発動時の注文価格)\n", price)
	fmt.Printf("  trigger_price: %.0f (発動条件価格)\n", triggerPrice)
	fmt.Printf("  order_type: %s\n", orderType)

	// 変換結果
	orderPriceResult := formatOrderPrice(price, orderType)
	gyakusasiPriceResult := formatGyakusasiPrice(price, orderType)
	triggerPriceResult := formatTriggerPrice(&triggerPrice, orderType)
	gyakusasiOrderTypeResult := convertGyakusasiOrderType(orderType)

	fmt.Printf("\n変換結果:\n")
	fmt.Printf("  OrderPrice: %s\n", orderPriceResult)
	fmt.Printf("  GyakusasiPrice: %s\n", gyakusasiPriceResult)
	fmt.Printf("  GyakusasiZyouken: %s\n", triggerPriceResult)
	fmt.Printf("  GyakusasiOrderType: %s\n", gyakusasiOrderTypeResult)

	// 期待される結果との比較
	fmt.Printf("\n期待される結果:\n")
	fmt.Printf("  OrderPrice: 972 (price と同じ)\n")
	fmt.Printf("  GyakusasiPrice: 972 (price と同じ)\n")
	fmt.Printf("  GyakusasiZyouken: 974 (trigger_price と同じ)\n")
	fmt.Printf("  GyakusasiOrderType: 2\n")

	// 問題の分析
	fmt.Printf("\n結果の検証:\n")
	if orderPriceResult == "972" {
		fmt.Printf("  ✅ OrderPrice正しい: %s\n", orderPriceResult)
	} else {
		fmt.Printf("  ❌ OrderPrice不正: 期待値=972, 実際=%s\n", orderPriceResult)
	}

	if gyakusasiPriceResult == "972" {
		fmt.Printf("  ✅ GyakusasiPrice正しい: %s\n", gyakusasiPriceResult)
	} else {
		fmt.Printf("  ❌ GyakusasiPrice不正: 期待値=972, 実際=%s\n", gyakusasiPriceResult)
	}

	if triggerPriceResult == "974" {
		fmt.Printf("  ✅ GyakusasiZyouken正しい: %s\n", triggerPriceResult)
	} else {
		fmt.Printf("  ❌ GyakusasiZyouken不正: 期待値=974, 実際=%s\n", triggerPriceResult)
	}

	if gyakusasiOrderTypeResult == "2" {
		fmt.Printf("  ✅ GyakusasiOrderType正しい: %s\n", gyakusasiOrderTypeResult)
	} else {
		fmt.Printf("  ❌ GyakusasiOrderType不正: 期待値=2, 実際=%s\n", gyakusasiOrderTypeResult)
	}

	// 修正案の提案
	fmt.Printf("\n修正案:\n")
	fmt.Printf("STOP_LIMIT注文では:\n")
	fmt.Printf("  - OrderPrice: 逆指値時の注文価格（price）を使用\n")
	fmt.Printf("  - GyakusasiPrice: 逆指値時の注文価格（price）を使用\n")
	fmt.Printf("  - GyakusasiZyouken: 逆指値条件価格（trigger_price）を使用\n")
	fmt.Printf("  - GyakusasiOrderType: \"2\"（通常+逆指値）を使用\n")

	fmt.Println("\n=== テスト完了 ===")
}
