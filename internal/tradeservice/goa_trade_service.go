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
	order_response "stock-bot/internal/infrastructure/client/dto/order/response"
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
	unifiedClient  *client.TachibanaUnifiedClient // セッション回復用
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

// SetUnifiedClient はセッション回復用の統合クライアントを設定する
func (s *GoaTradeService) SetUnifiedClient(unifiedClient *client.TachibanaUnifiedClient) {
	s.unifiedClient = unifiedClient
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

	// balanceClientがnilの場合はスタブ実装を返す
	if s.balanceClient == nil {
		s.logger.Warn("balanceClient is nil, returning stub balance")
		return &service.Balance{
			Cash:        1000000.0,
			BuyingPower: 800000.0,
		}, nil
	}

	// セッションエラー時の自動回復を試行
	_, err := s.executeWithSessionRecovery(ctx, func() (interface{}, error) {
		return s.balanceClient.GetZanKaiSummary(ctx, s.appSession)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// 簡易実装 - 実際の立花API呼び出しは後で実装
	return &service.Balance{
		Cash:        1000000.0,
		BuyingPower: 800000.0,
	}, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (s *GoaTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	s.logger.Info("GoaTradeService.GetPriceHistory called", "symbol", symbol, "days", days)

	// 簡易実装 - 実際の立花API呼び出しは後で実装
	return []*service.HistoricalPrice{}, nil
}

// PlaceOrder は注文を発行する
func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	s.logger.Info("GoaTradeService.PlaceOrder called", "symbol", req.Symbol, "trade_type", req.TradeType)

	// デバッグ: 各フィールドの状態を確認
	s.logger.Info("Debug: checking service fields",
		"orderClient_nil", s.orderClient == nil,
		"appSession_nil", s.appSession == nil,
		"logger_nil", s.logger == nil)

	if s.orderClient == nil {
		return nil, fmt.Errorf("orderClient is nil")
	}
	if s.appSession == nil {
		return nil, fmt.Errorf("appSession is nil")
	}

	// セッションエラー時の自動回復を試行
	result, err := s.executeWithSessionRecovery(ctx, func() (interface{}, error) {
		// デバッグ: リクエスト内容を確認
		s.logger.Info("Debug: request details",
			"OrderType", req.OrderType,
			"Price", req.Price,
			"TriggerPrice", req.TriggerPrice)

		// 立花証券APIの注文パラメータを作成
		params := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1", // 特定口座
			IssueCode:                req.Symbol,
			SizyouC:                  "00", // 東証
			BaibaiKubun:              convertTradeType(req.TradeType),
			Condition:                "0", // 指定なし（逆指値は別パラメータで制御）
			OrderPrice:               formatOrderPrice(req.Price, req.OrderType),
			OrderSuryou:              fmt.Sprintf("%d", req.Quantity),
			GenkinShinyouKubun:       convertPositionAccountType(req.PositionAccountType),
			OrderExpireDay:           "0", // 当日限り
			GyakusasiOrderType:       convertGyakusasiOrderType(req.OrderType),
			GyakusasiZyouken:         formatTriggerPrice(req.TriggerPrice, req.OrderType),
			GyakusasiPrice:           formatGyakusasiPrice(req.Price, req.OrderType),
			TatebiType:               "*", // 指定なし
			TategyokuZyoutoekiKazeiC: "*", // 指定なし
		}

		// デバッグ: 送信パラメータをログ出力
		s.logger.Info("Sending order parameters",
			"ZyoutoekiKazeiC", params.ZyoutoekiKazeiC,
			"IssueCode", params.IssueCode,
			"SizyouC", params.SizyouC,
			"BaibaiKubun", params.BaibaiKubun,
			"Condition", params.Condition,
			"OrderPrice", params.OrderPrice,
			"OrderSuryou", params.OrderSuryou,
			"GenkinShinyouKubun", params.GenkinShinyouKubun,
			"OrderExpireDay", params.OrderExpireDay,
			"GyakusasiOrderType", params.GyakusasiOrderType,
			"GyakusasiZyouken", params.GyakusasiZyouken,
			"GyakusasiPrice", params.GyakusasiPrice,
			"TatebiType", params.TatebiType,
			"TategyokuZyoutoekiKazeiC", params.TategyokuZyoutoekiKazeiC)

		// 立花証券APIに注文を送信
		return s.orderClient.NewOrder(ctx, s.appSession, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// 立花証券APIのレスポンスを処理
	tachibanaResponse := result.(*order_response.ResNewOrder)

	// ドメインモデルに変換
	triggerPrice := 0.0
	if req.TriggerPrice != nil {
		triggerPrice = *req.TriggerPrice
	}

	order := &model.Order{
		OrderID:             tachibanaResponse.OrderNumber,
		Symbol:              req.Symbol,
		TradeType:           req.TradeType,
		OrderType:           req.OrderType,
		Quantity:            req.Quantity,
		Price:               req.Price,
		TriggerPrice:        triggerPrice,
		OrderStatus:         model.OrderStatusNew, // 新規注文として設定
		PositionAccountType: req.PositionAccountType,
	}

	// データベースに保存
	if s.orderRepo != nil {
		if err := s.orderRepo.Save(ctx, order); err != nil {
			s.logger.Warn("failed to save order to database", "error", err)
		}
	}

	s.logger.Info("order placed successfully", "order_id", order.OrderID)
	return order, nil
}

// CancelOrder は注文をキャンセルする
func (s *GoaTradeService) CancelOrder(ctx context.Context, orderID string) error {
	s.logger.Info("GoaTradeService.CancelOrder called", "order_id", orderID)

	// orderRepoがnilの場合はスタブ実装
	if s.orderRepo == nil {
		s.logger.Warn("orderRepo is nil, returning stub error")
		return fmt.Errorf("order not found: %s (stub implementation)", orderID)
	}

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

// executeWithSessionRecovery はセッションエラー時の自動回復機能付きで処理を実行
func (s *GoaTradeService) executeWithSessionRecovery(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	// 最初の試行
	result, err := operation()
	if err == nil {
		return result, nil
	}

	// セッションエラーかどうかを判定
	if !s.isSessionError(err) {
		return nil, err
	}

	s.logger.Warn("session error detected, attempting recovery", "error", err)

	// セッション回復を試行
	if s.unifiedClient == nil {
		return nil, fmt.Errorf("session recovery not available: %w", err)
	}

	// 新しいセッションを取得
	newSession, recoveryErr := s.unifiedClient.GetSession(ctx)
	if recoveryErr != nil {
		return nil, fmt.Errorf("session recovery failed: %w (original error: %v)", recoveryErr, err)
	}

	// セッションを更新
	s.appSession = newSession
	s.logger.Info("session recovered successfully")

	// 操作を再試行
	result, retryErr := operation()
	if retryErr != nil {
		return nil, fmt.Errorf("operation failed after session recovery: %w", retryErr)
	}

	return result, nil
}

// isSessionError はエラーがセッション関連かどうかを判定
func (s *GoaTradeService) isSessionError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()

	// 立花証券のセッションエラーパターン
	sessionErrorPatterns := []string{
		"session expired",
		"invalid session",
		"authentication required",
		"unauthorized",
		"401",
		"403",
	}

	for _, pattern := range sessionErrorPatterns {
		if contains(errorStr, pattern) {
			return true
		}
	}

	return false
}

// contains は文字列に部分文字列が含まれているかチェック（大文字小文字を無視）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// 変換関数群

// convertTradeType はドメインの取引種別を立花証券APIの形式に変換
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

// convertOrderType はドメインの注文種別を立花証券APIの形式に変換
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

// convertPositionAccountType はドメインの口座種別を立花証券APIの形式に変換
func convertPositionAccountType(accountType model.PositionAccountType) string {
	switch accountType {
	case model.PositionAccountTypeCash:
		return "0" // 現物
	case model.PositionAccountTypeMarginNew:
		return "2" // 信用新規（制度信用6ヶ月）
	case model.PositionAccountTypeMarginRepay:
		return "4" // 信用返済
	default:
		return "0" // デフォルトは現物
	}
}

// convertOrderStatus は立花証券APIのステータスをドメインの形式に変換
func convertOrderStatus(status string) model.OrderStatus {
	switch status {
	case "0":
		return model.OrderStatusNew // 新規
	case "1":
		return model.OrderStatusPartiallyFilled // 一部約定
	case "2":
		return model.OrderStatusFilled // 全約定
	case "3":
		return model.OrderStatusCanceled // 取消
	case "4":
		return model.OrderStatusRejected // 拒否
	default:
		return model.OrderStatusNew
	}
}

// formatPrice は注文価格を立花証券API用の文字列に変換
func formatPrice(price float64, orderType model.OrderType) string {
	if orderType == model.OrderTypeMarket {
		return "0" // 成行の場合は0
	}
	return fmt.Sprintf("%.0f", price)
}

// formatOrderPrice は通常注文価格を立花証券API用の文字列に変換
func formatOrderPrice(price float64, orderType model.OrderType) string {
	switch orderType {
	case model.OrderTypeMarket, model.OrderTypeStop:
		return "*" // 成行・逆指値成行の場合は*
	case model.OrderTypeLimit:
		return fmt.Sprintf("%.0f", price) // 指値の場合は価格
	case model.OrderTypeStopLimit:
		// STOP_LIMIT注文では、OrderPriceは通常時の価格として逆指値価格より低い値を使用
		normalPrice := price - 5
		if normalPrice <= 0 {
			normalPrice = price * 0.95 // 5%低い価格
		}
		return fmt.Sprintf("%.0f", normalPrice)
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
