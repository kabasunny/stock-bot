package web

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/gen/trade"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/tradeservice"
	"time"
)

// TradeService はTradeServiceのHTTP実装
type TradeService struct {
	tradeService service.TradeService
	logger       *slog.Logger
	session      *client.Session
}

// NewTradeService は新しいTradeServiceを作成する
func NewTradeService(tradeService service.TradeService, logger *slog.Logger, session *client.Session) *TradeService {
	return &TradeService{
		tradeService: tradeService,
		logger:       logger,
		session:      session,
	}
}

// GetSession はセッション情報を取得する
func (s *TradeService) GetSession(ctx context.Context) (*trade.GetSessionResult, error) {
	s.logger.Info("TradeService.GetSession called")

	session := s.tradeService.GetSession()
	if session == nil {
		return nil, fmt.Errorf("no active session")
	}

	return &trade.GetSessionResult{
		SessionID: session.SessionID,
		UserID:    session.UserID,
		LoginTime: session.LoginTime.Format(time.RFC3339),
	}, nil
}

// GetPositions は保有ポジションを取得する
func (s *TradeService) GetPositions(ctx context.Context) (*trade.GetPositionsResult, error) {
	s.logger.Info("TradeService.GetPositions called")

	positions, err := s.tradeService.GetPositions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	tradePositions := make([]*trade.TradePositionResult, len(positions))
	for i, pos := range positions {
		tradePositions[i] = &trade.TradePositionResult{
			Symbol:              pos.Symbol,
			PositionType:        convertPositionType(pos.PositionType),
			PositionAccountType: convertPositionAccountType(pos.PositionAccountType),
			AveragePrice:        pos.AveragePrice,
			Quantity:            uint(pos.Quantity),
		}
	}

	return &trade.GetPositionsResult{
		Positions: tradePositions,
	}, nil
}

// GetOrders は注文一覧を取得する
func (s *TradeService) GetOrders(ctx context.Context) (*trade.GetOrdersResult, error) {
	s.logger.Info("TradeService.GetOrders called")

	orders, err := s.tradeService.GetOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	tradeOrders := make([]*trade.TradeOrderResult, len(orders))
	for i, order := range orders {
		convertedAccountType := convertPositionAccountType(order.PositionAccountType)
		tradeOrders[i] = &trade.TradeOrderResult{
			OrderID:             order.OrderID,
			Symbol:              order.Symbol,
			TradeType:           convertTradeType(order.TradeType),
			OrderType:           convertOrderType(order.OrderType),
			Quantity:            uint(order.Quantity),
			Price:               order.Price,
			OrderStatus:         convertOrderStatus(order.OrderStatus),
			PositionAccountType: &convertedAccountType,
		}
	}

	return &trade.GetOrdersResult{
		Orders: tradeOrders,
	}, nil
}

// GetBalance は残高情報を取得する
func (s *TradeService) GetBalance(ctx context.Context) (*trade.TradeBalanceResult, error) {
	s.logger.Info("TradeService.GetBalance called")

	balance, err := s.tradeService.GetBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &trade.TradeBalanceResult{
		Cash:        balance.Cash,
		BuyingPower: balance.BuyingPower,
	}, nil
}

// GetPriceHistory は価格履歴を取得する
func (s *TradeService) GetPriceHistory(ctx context.Context, p *trade.GetPriceHistoryPayload) (*trade.GetPriceHistoryResult, error) {
	s.logger.Info("TradeService.GetPriceHistory called", "symbol", p.Symbol, "days", p.Days)

	days := int(p.Days)
	if days == 0 {
		days = 30 // デフォルト値
	}

	history, err := s.tradeService.GetPriceHistory(ctx, p.Symbol, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}

	tradeHistory := make([]*trade.TradePriceHistoryItem, len(history))
	for i, item := range history {
		tradeHistory[i] = &trade.TradePriceHistoryItem{
			Date:   item.Date.Format(time.RFC3339),
			Open:   item.Open,
			High:   item.High,
			Low:    item.Low,
			Close:  item.Close,
			Volume: uint64(item.Volume),
		}
	}

	return &trade.GetPriceHistoryResult{
		Symbol:  p.Symbol,
		History: tradeHistory,
	}, nil
}

// PlaceOrder は注文を発行する
func (s *TradeService) PlaceOrder(ctx context.Context, p *trade.PlaceOrderPayload) (*trade.TradeOrderResult, error) {
	s.logger.Info("TradeService.PlaceOrder called", "symbol", p.Symbol, "trade_type", p.TradeType, "quantity", p.Quantity)

	req := &service.PlaceOrderRequest{
		Symbol:              p.Symbol,
		TradeType:           convertTradeTypeFromAPI(p.TradeType),
		OrderType:           convertOrderTypeFromAPI(p.OrderType),
		Quantity:            int(p.Quantity),
		Price:               p.Price,
		TriggerPrice:        &p.TriggerPrice,
		PositionAccountType: convertPositionAccountTypeFromAPI(p.PositionAccountType),
	}

	order, err := s.tradeService.PlaceOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	convertedAccountType := convertPositionAccountType(order.PositionAccountType)
	return &trade.TradeOrderResult{
		OrderID:             order.OrderID,
		Symbol:              order.Symbol,
		TradeType:           convertTradeType(order.TradeType),
		OrderType:           convertOrderType(order.OrderType),
		Quantity:            uint(order.Quantity),
		Price:               order.Price,
		OrderStatus:         convertOrderStatus(order.OrderStatus),
		PositionAccountType: &convertedAccountType,
	}, nil
}

// CancelOrder は注文をキャンセルする
func (s *TradeService) CancelOrder(ctx context.Context, p *trade.CancelOrderPayload) error {
	s.logger.Info("TradeService.CancelOrder called", "order_id", p.OrderID)

	return s.tradeService.CancelOrder(ctx, p.OrderID)
}

// ValidateSymbol は銘柄の妥当性をチェックし、取引情報を返す
// ValidateSymbol は銘柄の妥当性をチェックし、取引情報を返す
func (s *TradeService) ValidateSymbol(ctx context.Context, p *trade.ValidateSymbolPayload) (*trade.ValidateSymbolResult, error) {
	s.logger.Info("TradeService.ValidateSymbol called", "symbol", p.Symbol)

	// TradeServiceから銘柄情報を取得
	// 型アサーションでGoaTradeServiceにアクセス
	goaTradeService, ok := s.tradeService.(*tradeservice.GoaTradeService)
	if !ok {
		// フォールバック: 簡易実装
		return &trade.ValidateSymbolResult{
			Valid:  true,
			Symbol: p.Symbol,
			Name:   func() *string { s := "銘柄名（取得不可）"; return &s }(),
		}, nil
	}

	// マスターデータから銘柄情報を取得
	stockInfo, err := goaTradeService.GetStockInfo(ctx, p.Symbol)
	if err != nil {
		// 銘柄が見つからない場合
		return &trade.ValidateSymbolResult{
			Valid:  false,
			Symbol: p.Symbol,
		}, nil
	}

	return &trade.ValidateSymbolResult{
		Valid:       true,
		Symbol:      p.Symbol,
		Name:        &stockInfo.Name,
		TradingUnit: func() *uint { u := uint(stockInfo.TradingUnit); return &u }(),
		Market:      &stockInfo.Market,
	}, nil
}

// 型変換ヘルパー関数

func convertPositionType(pt model.PositionType) string {
	switch pt {
	case model.PositionTypeLong:
		return "LONG"
	case model.PositionTypeShort:
		return "SHORT"
	default:
		return "LONG"
	}
}

func convertPositionAccountType(pat model.PositionAccountType) string {
	switch pat {
	case model.PositionAccountTypeCash:
		return "CASH"
	case model.PositionAccountTypeMarginNew:
		return "MARGIN_NEW"
	case model.PositionAccountTypeMarginRepay:
		return "MARGIN_REPAY"
	default:
		return "CASH"
	}
}

func convertTradeType(tt model.TradeType) string {
	switch tt {
	case model.TradeTypeBuy:
		return "BUY"
	case model.TradeTypeSell:
		return "SELL"
	default:
		return "BUY"
	}
}

func convertOrderType(ot model.OrderType) string {
	switch ot {
	case model.OrderTypeMarket:
		return "MARKET"
	case model.OrderTypeLimit:
		return "LIMIT"
	case model.OrderTypeStop:
		return "STOP"
	case model.OrderTypeStopLimit:
		return "STOP_LIMIT"
	default:
		return "MARKET"
	}
}

func convertOrderStatus(os model.OrderStatus) string {
	switch os {
	case model.OrderStatusNew:
		return "NEW"
	case model.OrderStatusPartiallyFilled:
		return "PARTIALLY_FILLED"
	case model.OrderStatusFilled:
		return "FILLED"
	case model.OrderStatusCanceled:
		return "CANCELLED"
	case model.OrderStatusRejected:
		return "REJECTED"
	default:
		return "NEW"
	}
}

func convertTradeTypeFromAPI(tt string) model.TradeType {
	switch tt {
	case "BUY":
		return model.TradeTypeBuy
	case "SELL":
		return model.TradeTypeSell
	default:
		return model.TradeTypeBuy
	}
}

func convertOrderTypeFromAPI(ot string) model.OrderType {
	switch ot {
	case "MARKET":
		return model.OrderTypeMarket
	case "LIMIT":
		return model.OrderTypeLimit
	case "STOP":
		return model.OrderTypeStop
	case "STOP_LIMIT":
		return model.OrderTypeStopLimit
	default:
		return model.OrderTypeMarket
	}
}

func convertPositionAccountTypeFromAPI(pat string) model.PositionAccountType {
	switch pat {
	case "CASH":
		return model.PositionAccountTypeCash
	case "MARGIN_NEW":
		return model.PositionAccountTypeMarginNew
	case "MARGIN_REPAY":
		return model.PositionAccountTypeMarginRepay
	default:
		return model.PositionAccountTypeCash
	}
}

// GetOrderHistory は注文履歴を取得する
func (s *TradeService) GetOrderHistory(ctx context.Context, p *trade.GetOrderHistoryPayload) (*trade.GetOrderHistoryResult, error) {
	s.logger.Info("TradeService.GetOrderHistory called", "status", p.Status, "symbol", p.Symbol, "limit", p.Limit)

	// パラメータの変換
	var status *model.OrderStatus
	if p.Status != nil {
		orderStatus := convertOrderStatusFromAPI(*p.Status)
		status = &orderStatus
	}

	var symbol *string
	if p.Symbol != nil && *p.Symbol != "" {
		symbol = p.Symbol
	}

	limit := int(p.Limit)
	if limit == 0 {
		limit = 100 // デフォルト値
	}

	// TradeServiceから履歴を取得
	orders, err := s.tradeService.GetOrderHistory(ctx, status, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	// レスポンス形式に変換
	historyOrders := make([]*trade.TradeOrderHistoryResult, len(orders))
	for i, order := range orders {
		// 約定履歴の変換
		executions := make([]*trade.TradeExecutionResult, len(order.Executions))
		for j, exec := range order.Executions {
			executions[j] = &trade.TradeExecutionResult{
				ExecutionID:      exec.ExecutionID,
				ExecutedQuantity: uint(exec.Quantity),
				ExecutedPrice:    exec.Price,
				ExecutedAt:       exec.ExecutedAt.Format(time.RFC3339),
			}
		}

		convertedAccountType := convertPositionAccountType(order.PositionAccountType)
		historyOrders[i] = &trade.TradeOrderHistoryResult{
			OrderID:             order.OrderID,
			Symbol:              order.Symbol,
			TradeType:           convertTradeType(order.TradeType),
			OrderType:           convertOrderType(order.OrderType),
			Quantity:            uint(order.Quantity),
			Price:               order.Price,
			OrderStatus:         convertOrderStatus(order.OrderStatus),
			PositionAccountType: &convertedAccountType,
			CreatedAt:           order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:           func() *string { s := order.UpdatedAt.Format(time.RFC3339); return &s }(),
			Executions:          executions,
		}
	}

	return &trade.GetOrderHistoryResult{
		Orders: historyOrders,
	}, nil
}

// convertOrderStatusFromAPI はAPI文字列をmodel.OrderStatusに変換する
func convertOrderStatusFromAPI(status string) model.OrderStatus {
	switch status {
	case "NEW":
		return model.OrderStatusNew
	case "PARTIALLY_FILLED":
		return model.OrderStatusPartiallyFilled
	case "FILLED":
		return model.OrderStatusFilled
	case "CANCELLED":
		return model.OrderStatusCanceled
	case "REJECTED":
		return model.OrderStatusRejected
	default:
		return model.OrderStatusNew
	}
}

// CorrectOrder は注文を訂正する
func (s *TradeService) CorrectOrder(ctx context.Context, p *trade.CorrectOrderPayload) (*trade.TradeOrderResult, error) {
	s.logger.Info("TradeService.CorrectOrder called", "order_id", p.OrderID, "price", p.Price, "quantity", p.Quantity)

	// パラメータの準備
	var newPrice *float64
	var newQuantity *int

	if p.Price != nil {
		newPrice = p.Price
	}

	if p.Quantity != nil {
		qty := int(*p.Quantity)
		newQuantity = &qty
	}

	// TradeServiceで訂正実行
	order, err := s.tradeService.CorrectOrder(ctx, p.OrderID, newPrice, newQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to correct order: %w", err)
	}

	// レスポンス形式に変換
	convertedAccountType := convertPositionAccountType(order.PositionAccountType)
	return &trade.TradeOrderResult{
		OrderID:             order.OrderID,
		Symbol:              order.Symbol,
		TradeType:           convertTradeType(order.TradeType),
		OrderType:           convertOrderType(order.OrderType),
		Quantity:            uint(order.Quantity),
		Price:               order.Price,
		OrderStatus:         convertOrderStatus(order.OrderStatus),
		PositionAccountType: &convertedAccountType,
	}, nil
}

// CancelAllOrders は全ての未約定注文をキャンセルする
func (s *TradeService) CancelAllOrders(ctx context.Context) (*trade.CancelAllOrdersResult, error) {
	s.logger.Info("TradeService.CancelAllOrders called")

	cancelledCount, err := s.tradeService.CancelAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel all orders: %w", err)
	}

	return &trade.CancelAllOrdersResult{
		CancelledCount: uint(cancelledCount),
	}, nil
}

// HealthCheck はサービスの健康状態をチェックする
func (s *TradeService) HealthCheck(ctx context.Context) (*trade.HealthCheckResult, error) {
	s.logger.Debug("TradeService.HealthCheck called")

	// TradeServiceから健康状態を取得
	goaTradeService, ok := s.tradeService.(*tradeservice.GoaTradeService)
	if !ok {
		// フォールバック
		return &trade.HealthCheckResult{
			Status:    "unhealthy",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	healthStatus, err := goaTradeService.HealthCheck(ctx)
	if err != nil {
		return &trade.HealthCheckResult{
			Status:    "unhealthy",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	return &trade.HealthCheckResult{
		Status:             healthStatus.Status,
		Timestamp:          healthStatus.Timestamp.Format(time.RFC3339),
		SessionValid:       &healthStatus.SessionValid,
		DatabaseConnected:  &healthStatus.DatabaseConnected,
		WebsocketConnected: &healthStatus.WebSocketConnected,
	}, nil
}
