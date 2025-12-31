package client

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
)

// HTTPTradeService はHTTP経由でTradeServiceにアクセスするクライアント（プレースホルダー実装）
type HTTPTradeService struct {
	endpoint string
}

// NewHTTPTradeService は新しいHTTPTradeServiceを作成する
func NewHTTPTradeService(endpoint string) *HTTPTradeService {
	return &HTTPTradeService{
		endpoint: endpoint,
	}
}

// GetSession は現在のAPIセッション情報を取得する
func (h *HTTPTradeService) GetSession() *client.Session {
	// HTTP経由の場合、セッション情報は直接取得できないため、
	// 実際の実装では認証トークンなどを管理する必要がある
	return nil
}

// GetPositions は現在の保有ポジションを取得する
func (h *HTTPTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	// TODO: 実際のHTTP APIクライアント実装
	return nil, fmt.Errorf("HTTP TradeService not implemented yet")
}

// GetOrders は発注中の注文を取得する
func (h *HTTPTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	// TODO: 実際のHTTP APIクライアント実装
	return nil, fmt.Errorf("HTTP TradeService not implemented yet")
}

// GetBalance は口座残高を取得する
func (h *HTTPTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	// TODO: 実際のHTTP APIクライアント実装
	return nil, fmt.Errorf("HTTP TradeService not implemented yet")
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (h *HTTPTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	// TODO: 実際のHTTP APIクライアント実装
	return nil, fmt.Errorf("HTTP TradeService not implemented yet")
}

// PlaceOrder は注文を発行する
func (h *HTTPTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	// TODO: 実際のHTTP APIクライアント実装
	return nil, fmt.Errorf("HTTP TradeService not implemented yet")
}

// CancelOrder は注文をキャンセルする
func (h *HTTPTradeService) CancelOrder(ctx context.Context, orderID string) error {
	// TODO: 実際のHTTP APIクライアント実装
	return fmt.Errorf("HTTP TradeService not implemented yet")
}
