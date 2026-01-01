package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 各注文タイプのリクエスト構造体
type OrderRequest struct {
	Symbol              string  `json:"symbol"`
	TradeType           string  `json:"trade_type"`
	OrderType           string  `json:"order_type"`
	Quantity            int     `json:"quantity"`
	Price               float64 `json:"price,omitempty"`
	TriggerPrice        float64 `json:"trigger_price,omitempty"`
	PositionAccountType string  `json:"position_account_type"`
}

type OrderResponse struct {
	OrderID             string  `json:"order_id"`
	Symbol              string  `json:"symbol"`
	TradeType           string  `json:"trade_type"`
	OrderType           string  `json:"order_type"`
	Quantity            int     `json:"quantity"`
	Price               float64 `json:"price"`
	TriggerPrice        float64 `json:"trigger_price"`
	PositionAccountType string  `json:"position_account_type"`
	OrderStatus         string  `json:"order_status"`
}

func main() {
	fmt.Println("=== 注文パラメータ検証テスト ===")

	// 1. 現在の株価を取得して適切な価格範囲を確認
	currentPrice := 2800.0 // トヨタ自動車の参考価格
	fmt.Printf("現在価格（参考）: %.0f円\n", currentPrice)

	// 価格範囲を計算（±10%程度）
	lowerPrice := currentPrice * 0.9
	upperPrice := currentPrice * 1.1
	triggerPriceBuy := currentPrice * 1.05  // 買い逆指値用
	triggerPriceSell := currentPrice * 0.95 // 売り逆指値用

	fmt.Printf("価格範囲: %.0f - %.0f円\n", lowerPrice, upperPrice)

	// 2. 各注文タイプのテストケース
	testCases := []struct {
		name    string
		request OrderRequest
	}{
		{
			name: "成行注文（現物）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "BUY",
				OrderType:           "MARKET",
				Quantity:            100,
				PositionAccountType: "CASH",
			},
		},
		{
			name: "成行注文（信用新規）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "BUY",
				OrderType:           "MARKET",
				Quantity:            100,
				PositionAccountType: "MARGIN_NEW",
			},
		},
		{
			name: "指値注文（現物）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "BUY",
				OrderType:           "LIMIT",
				Quantity:            100,
				Price:               lowerPrice,
				PositionAccountType: "CASH",
			},
		},
		{
			name: "指値注文（信用新規）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "SELL",
				OrderType:           "LIMIT",
				Quantity:            100,
				Price:               upperPrice,
				PositionAccountType: "MARGIN_NEW",
			},
		},
		{
			name: "逆指値注文（現物売り）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "SELL",
				OrderType:           "STOP",
				Quantity:            100,
				TriggerPrice:        triggerPriceSell,
				PositionAccountType: "CASH",
			},
		},
		{
			name: "逆指値注文（信用返済）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "BUY",
				OrderType:           "STOP",
				Quantity:            100,
				TriggerPrice:        triggerPriceBuy,
				PositionAccountType: "MARGIN_REPAY",
			},
		},
		{
			name: "逆指値指値注文（現物買い）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "BUY",
				OrderType:           "STOP_LIMIT",
				Quantity:            100,
				Price:               upperPrice,
				TriggerPrice:        triggerPriceBuy,
				PositionAccountType: "CASH",
			},
		},
		{
			name: "逆指値指値注文（信用新規売り）",
			request: OrderRequest{
				Symbol:              "7203",
				TradeType:           "SELL",
				OrderType:           "STOP_LIMIT",
				Quantity:            100,
				Price:               lowerPrice,
				TriggerPrice:        triggerPriceSell,
				PositionAccountType: "MARGIN_NEW",
			},
		},
	}

	// 3. 各テストケースを実行
	baseURL := "http://localhost:8080"

	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Printf("   リクエスト: %+v\n", tc.request)

		// HTTP リクエストを送信
		success, response := sendOrderRequest(baseURL, tc.request)
		if success {
			fmt.Printf("   ✅ 成功: %s\n", response)
		} else {
			fmt.Printf("   ❌ エラー: %s\n", response)
		}

		// API制限を考慮して少し待機
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n=== テスト完了 ===")
}

// sendOrderRequest はGoaサービスに注文リクエストを送信します
func sendOrderRequest(baseURL string, req OrderRequest) (bool, string) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return false, fmt.Sprintf("JSON変換エラー: %v", err)
	}

	resp, err := http.Post(
		baseURL+"/trade/orders",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return false, fmt.Sprintf("HTTP エラー: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Sprintf("レスポンス読み取りエラー: %v", err)
	}

	if resp.StatusCode == 201 {
		// 成功レスポンスをパース
		var orderResp OrderResponse
		if err := json.Unmarshal(body, &orderResp); err == nil {
			return true, fmt.Sprintf("注文ID: %s, ステータス: %s", orderResp.OrderID, orderResp.OrderStatus)
		}
		return true, string(body)
	} else {
		return false, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
	}
}
