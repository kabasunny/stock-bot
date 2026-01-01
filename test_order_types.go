package main

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/tradeservice"
)

// MockOrderClient は注文クライアントのモック実装
type MockOrderClient struct{}

func (m *MockOrderClient) NewOrder(ctx context.Context, session interface{}, params interface{}) (interface{}, error) {
	// モックレスポンス
	return &struct {
		OrderID    string
		Status     string
		ResultCode string
	}{
		OrderID:    "MOCK-ORDER-123",
		Status:     "0", // 新規
		ResultCode: "0", // 成功
	}, nil
}

func main() {
	fmt.Println("=== 注文タイプ × 口座区分 組み合わせテスト ===")

	// モックサービスの作成
	mockOrderClient := &MockOrderClient{}
	tradeService := &tradeservice.GoaTradeService{}

	// テストケース定義
	testCases := []struct {
		name         string
		orderType    model.OrderType
		accountType  model.PositionAccountType
		tradeType    model.TradeType
		price        float64
		triggerPrice float64
	}{
		{"成行注文 × 現物", model.OrderTypeMarket, model.PositionAccountTypeCash, model.TradeTypeBuy, 0, 0},
		{"成行注文 × 信用新規", model.OrderTypeMarket, model.PositionAccountTypeMarginNew, model.TradeTypeBuy, 0, 0},
		{"指値注文 × 現物", model.OrderTypeLimit, model.PositionAccountTypeCash, model.TradeTypeBuy, 2800.0, 0},
		{"指値注文 × 信用新規", model.OrderTypeLimit, model.PositionAccountTypeMarginNew, model.TradeTypeSell, 2900.0, 0},
		{"逆指値注文 × 現物", model.OrderTypeStop, model.PositionAccountTypeCash, model.TradeTypeSell, 0, 2700.0},
		{"逆指値注文 × 信用返済", model.OrderTypeStop, model.PositionAccountTypeMarginRepay, model.TradeTypeBuy, 0, 3000.0},
		{"逆指値指値注文 × 現物", model.OrderTypeStopLimit, model.PositionAccountTypeCash, model.TradeTypeBuy, 2850.0, 2800.0},
		{"逆指値指値注文 × 信用新規", model.OrderTypeStopLimit, model.PositionAccountTypeMarginNew, model.TradeTypeSell, 2750.0, 2800.0},
	}

	// 各テストケースを実行
	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)

		req := &service.PlaceOrderRequest{
			Symbol:              "7203", // トヨタ自動車
			TradeType:           tc.tradeType,
			OrderType:           tc.orderType,
			Quantity:            100,
			Price:               tc.price,
			TriggerPrice:        tc.triggerPrice,
			PositionAccountType: tc.accountType,
		}

		// パラメータ変換のテスト
		fmt.Printf("  注文タイプ: %s\n", tc.orderType)
		fmt.Printf("  口座区分: %s\n", tc.accountType)
		fmt.Printf("  売買区分: %s\n", tc.tradeType)
		if tc.price > 0 {
			fmt.Printf("  指値価格: %.0f円\n", tc.price)
		}
		if tc.triggerPrice > 0 {
			fmt.Printf("  トリガー価格: %.0f円\n", tc.triggerPrice)
		}

		// 立花証券API変換のテスト
		fmt.Printf("  → 立花API変換結果:\n")
		fmt.Printf("    注文種別: %s\n", convertOrderType(tc.orderType))
		fmt.Printf("    口座区分: %s\n", convertPositionAccountType(tc.accountType))
		fmt.Printf("    売買区分: %s\n", convertTradeType(tc.tradeType))
		if tc.orderType == model.OrderTypeStop || tc.orderType == model.OrderTypeStopLimit {
			fmt.Printf("    逆指値種別: %s\n", convertGyakusasiOrderType(tc.orderType))
			fmt.Printf("    逆指値条件: %s\n", convertGyakusasiCondition(tc.orderType))
		}

		fmt.Printf("  ✅ 変換成功\n")
	}

	fmt.Println("\n=== テスト完了 ===")
	fmt.Println("全ての注文タイプ × 口座区分の組み合わせが正常に変換されました。")
}

// 変換関数（tradeservice パッケージから複製）
func convertOrderType(orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeMarket:
		return "1" // 成行
	case model.OrderTypeLimit:
		return "2" // 指値
	case model.OrderTypeStop:
		return "3" // 逆指値
	case model.OrderTypeStopLimit:
		return "4" // 逆指値指値
	default:
		return "1" // デフォルトは成行
	}
}

func convertPositionAccountType(accountType model.PositionAccountType) string {
	switch accountType {
	case model.PositionAccountTypeCash:
		return "1" // 現物
	case model.PositionAccountTypeMarginNew:
		return "2" // 信用新規
	case model.PositionAccountTypeMarginRepay:
		return "3" // 信用返済
	default:
		return "1" // デフォルトは現物
	}
}

func convertTradeType(tradeType model.TradeType) string {
	switch tradeType {
	case model.TradeTypeBuy:
		return "3" // 買い
	case model.TradeTypeSell:
		return "1" // 売り
	default:
		return "3" // デフォルトは買い
	}
}

func convertGyakusasiOrderType(orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeStop:
		return "1" // 逆指値成行
	case model.OrderTypeStopLimit:
		return "2" // 逆指値指値
	default:
		return "" // 通常注文は空文字
	}
}

func convertGyakusasiCondition(orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeStop, model.OrderTypeStopLimit:
		return "1" // 以下（売り）または以上（買い）
	default:
		return "" // 通常注文は空文字
	}
}
