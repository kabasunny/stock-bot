package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// LightweightAgent は軽量なHTTPクライアントベースのエージェント
type LightweightAgent struct {
	httpClient *http.Client
	baseURL    string
	strategy   Strategy
	logger     *slog.Logger
}

// NewLightweightAgent は新しい軽量エージェントを作成する
func NewLightweightAgent(baseURL string, strategy Strategy, logger *slog.Logger) *LightweightAgent {
	return &LightweightAgent{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  baseURL,
		strategy: strategy,
		logger:   logger,
	}
}

// Run はエージェントのメインループを実行する
func (a *LightweightAgent) Run(ctx context.Context) error {
	a.logger.Info("Starting lightweight agent", "strategy", a.strategy.Name())

	ticker := time.NewTicker(a.strategy.GetExecutionInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Agent stopped by context")
			return ctx.Err()
		case <-ticker.C:
			if err := a.executeStrategy(ctx); err != nil {
				a.logger.Error("Strategy execution failed", "error", err)
				// エラーが発生してもループを継続
			}
		}
	}
}

// executeStrategy は戦略を実行する
func (a *LightweightAgent) executeStrategy(ctx context.Context) error {
	// 1. ヘルスチェック
	if !a.checkHealth(ctx) {
		return fmt.Errorf("service is not healthy")
	}

	// 2. 市場データ収集
	marketData, err := a.collectMarketData(ctx)
	if err != nil {
		return fmt.Errorf("failed to collect market data: %w", err)
	}

	// 3. 戦略評価
	signal, err := a.strategy.Evaluate(ctx, marketData)
	if err != nil {
		return fmt.Errorf("strategy evaluation failed: %w", err)
	}

	// 4. 注文実行
	if signal.ShouldTrade {
		order, err := a.placeOrder(ctx, signal.OrderRequest)
		if err != nil {
			return fmt.Errorf("failed to place order: %w", err)
		}
		a.logger.Info("Order placed successfully",
			"strategy", a.strategy.Name(),
			"order_id", order.OrderID,
			"symbol", order.Symbol,
			"trade_type", order.TradeType,
			"quantity", order.Quantity)
	} else {
		a.logger.Debug("No trading signal", "strategy", a.strategy.Name())
	}

	return nil
}

// collectMarketData は市場データを収集する
func (a *LightweightAgent) collectMarketData(ctx context.Context) (*MarketData, error) {
	// 残高取得
	balance, err := a.getBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// ポジション取得
	positions, err := a.getPositions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// 注文一覧取得
	orders, err := a.getOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return &MarketData{
		Balance:   balance,
		Positions: positions,
		Orders:    orders,
		Timestamp: time.Now(),
	}, nil
}

// PlaceOrderRequest は注文リクエスト
type PlaceOrderRequest struct {
	Symbol              string  `json:"symbol"`
	TradeType           string  `json:"trade_type"`
	OrderType           string  `json:"order_type"`
	Quantity            uint    `json:"quantity"`
	Price               float64 `json:"price"`
	PositionAccountType string  `json:"position_account_type"`
}

// Balance は残高情報
type Balance struct {
	Cash        float64 `json:"cash"`
	BuyingPower float64 `json:"buying_power"`
}

// Position はポジション情報
type Position struct {
	Symbol              string  `json:"symbol"`
	PositionType        string  `json:"position_type"`
	PositionAccountType string  `json:"position_account_type"`
	AveragePrice        float64 `json:"average_price"`
	Quantity            uint    `json:"quantity"`
}

// Order は注文情報
type Order struct {
	OrderID             string  `json:"order_id"`
	Symbol              string  `json:"symbol"`
	TradeType           string  `json:"trade_type"`
	OrderType           string  `json:"order_type"`
	Quantity            uint    `json:"quantity"`
	Price               float64 `json:"price"`
	OrderStatus         string  `json:"order_status"`
	PositionAccountType string  `json:"position_account_type"`
}

// checkHealth はサービスの健康状態をチェックする
func (a *LightweightAgent) checkHealth(ctx context.Context) bool {
	resp, err := a.httpClient.Get(a.baseURL + "/trade/health")
	if err != nil {
		a.logger.Error("Health check failed", "error", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// getBalance は残高を取得する
func (a *LightweightAgent) getBalance(ctx context.Context) (*Balance, error) {
	resp, err := a.httpClient.Get(a.baseURL + "/trade/balance")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var balance Balance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// getPositions はポジションを取得する
func (a *LightweightAgent) getPositions(ctx context.Context) ([]*Position, error) {
	resp, err := a.httpClient.Get(a.baseURL + "/trade/positions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result struct {
		Positions []*Position `json:"positions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Positions, nil
}

// getOrders は注文一覧を取得する
func (a *LightweightAgent) getOrders(ctx context.Context) ([]*Order, error) {
	resp, err := a.httpClient.Get(a.baseURL + "/trade/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result struct {
		Orders []*Order `json:"orders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Orders, nil
}

// placeOrder は注文を発行する
func (a *LightweightAgent) placeOrder(ctx context.Context, req *PlaceOrderRequest) (*Order, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Post(
		a.baseURL+"/trade/orders",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
