package tradeservice

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/adapter"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// GoaTradeService は service.TradeService インターフェースのGoaクライアント実装
type GoaTradeService struct {
	balanceClient  client.BalanceClient
	orderClient    client.OrderClient
	priceClient    client.PriceInfoClient
	orderRepo      repository.OrderRepository
	masterRepo     repository.MasterRepository
	appSession     *client.Session
	sessionAdapter *adapter.SessionAdapter
	logger         *slog.Logger
}

// NewGoaTradeService は GoaTradeService の新しいインスタンスを作成する
func NewGoaTradeService(
	balanceClient client.BalanceClient,
	orderClient client.OrderClient,
	priceClient client.PriceInfoClient,
	orderRepo repository.OrderRepository,
	masterRepo repository.MasterRepository,
	appSession *client.Session,
	logger *slog.Logger,
) *GoaTradeService {
	return &GoaTradeService{
		balanceClient:  balanceClient,
		orderClient:    orderClient,
		priceClient:    priceClient,
		orderRepo:      orderRepo,
		masterRepo:     masterRepo,
		appSession:     appSession,
		sessionAdapter: adapter.NewSessionAdapter(),
		logger:         logger,
	}
}

// GetSession は現在のAPIセッション情報を取得する
func (s *GoaTradeService) GetSession() *model.Session {
	return s.sessionAdapter.ToDomainSession(s.appSession)
}

// GetPositions は現在の保有ポジションを取得する
func (s *GoaTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	s.logger.Info("GoaTradeService.GetPositions called")
	// 簡易実装 - 実際の実装は後で追加
	return []*model.Position{}, nil
}

// GetOrders は発注中の注文を取得する
func (s *GoaTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	s.logger.Info("GoaTradeService.GetOrders called")
	// 簡易実装 - 実際の実装は後で追加
	return []*model.Order{}, nil
}

// GetBalance は口座残高を取得する
func (s *GoaTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	s.logger.Info("GoaTradeService.GetBalance called")
	// 簡易実装 - 実際の実装は後で追加
	return &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (s *GoaTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	s.logger.Info("GoaTradeService.GetPriceHistory called", "symbol", symbol, "days", days)
	// 簡易実装 - 実際の実装は後で追加
	return []*service.HistoricalPrice{}, nil
}

// PlaceOrder は注文を発行する
func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	s.logger.Info("GoaTradeService.PlaceOrder called", "symbol", req.Symbol, "trade_type", req.TradeType)
	// 簡易実装 - 実際の実装は後で追加
	return &model.Order{
		OrderID:             "test-order-123",
		Symbol:              req.Symbol,
		TradeType:           req.TradeType,
		OrderType:           req.OrderType,
		Quantity:            req.Quantity,
		Price:               req.Price,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: req.PositionAccountType,
	}, nil
}

// CancelOrder は注文をキャンセルする
func (s *GoaTradeService) CancelOrder(ctx context.Context, orderID string) error {
	s.logger.Info("GoaTradeService.CancelOrder called", "order_id", orderID)

	// 注文の存在確認
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// 注文状態の確認
	if !order.IsUnexecuted() {
		return fmt.Errorf("order %s cannot be cancelled (status: %s)", orderID, order.OrderStatus)
	}

	// 簡易実装 - 実際のキャンセル処理は後で追加
	return nil
}

// CorrectOrder は注文を訂正する
func (s *GoaTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	s.logger.Info("GoaTradeService.CorrectOrder called", "order_id", orderID)
	// 簡易実装 - 実際の実装は後で追加
	return &model.Order{
		OrderID:     orderID,
		OrderStatus: model.OrderStatusNew,
	}, nil
}

// CancelAllOrders は全ての未約定注文をキャンセルする
func (s *GoaTradeService) CancelAllOrders(ctx context.Context) (int, error) {
	s.logger.Info("GoaTradeService.CancelAllOrders called")
	// 簡易実装 - 実際の実装は後で追加
	return 0, nil
}

// GetOrderHistory は注文履歴を取得する
func (s *GoaTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	s.logger.Info("GoaTradeService.GetOrderHistory called")
	// 簡易実装 - 実際の実装は後で追加
	return []*model.Order{}, nil
}

// HealthCheck はサービスの健康状態をチェックする
func (s *GoaTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	s.logger.Debug("GoaTradeService.HealthCheck called")

	return &service.HealthStatus{
		Status:             "healthy",
		Timestamp:          time.Now(),
		SessionValid:       s.appSession != nil && s.appSession.ResultCode == "0",
		DatabaseConnected:  true, // 簡易実装
		WebSocketConnected: true, // 簡易実装
	}, nil
}

// GetStockInfo はマスターデータから銘柄情報を取得する（ValidateSymbolで使用）
func (s *GoaTradeService) GetStockInfo(ctx context.Context, symbol string) (*StockInfo, error) {
	s.logger.Info("GoaTradeService.GetStockInfo called", "symbol", symbol)
	// 簡易実装 - 実際の実装は後で追加
	return &StockInfo{
		Symbol:      symbol,
		Name:        "テスト銘柄",
		TradingUnit: 100,
		Market:      "東証プライム",
	}, nil
}

// StockInfo は銘柄情報を表す
type StockInfo struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	TradingUnit int    `json:"trading_unit"`
	Market      string `json:"market"`
}
