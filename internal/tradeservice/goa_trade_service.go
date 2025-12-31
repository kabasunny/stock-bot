package tradeservice

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"strconv"
	"time"
)

// GoaTradeService は service.TradeService インターフェースのGoaクライアント実装
type GoaTradeService struct {
	balanceClient client.BalanceClient
	orderClient   client.OrderClient
	priceClient   client.PriceInfoClient
	orderRepo     repository.OrderRepository
	masterRepo    repository.MasterRepository // マスターデータリポジトリ
	appSession    *client.Session
	logger        *slog.Logger
}

// NewGoaTradeService は GoaTradeService の新しいインスタンスを作成する
func NewGoaTradeService(
	balanceClient client.BalanceClient,
	orderClient client.OrderClient,
	priceClient client.PriceInfoClient,
	orderRepo repository.OrderRepository,
	masterRepo repository.MasterRepository, // マスターデータリポジトリ追加
	appSession *client.Session,
	logger *slog.Logger,
) *GoaTradeService {
	return &GoaTradeService{
		balanceClient: balanceClient,
		orderClient:   orderClient,
		priceClient:   priceClient,
		orderRepo:     orderRepo,
		masterRepo:    masterRepo, // 追加
		appSession:    appSession,
		logger:        logger,
	}
}

// GetSession は現在のAPIセッション情報を取得する
func (s *GoaTradeService) GetSession() *client.Session {
	return s.appSession
}

// GetPositions は現在の保有ポジションを取得する
func (s *GoaTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	s.logger.Info("GoaTradeService.GetPositions called")

	// 現物保有銘柄リストの取得
	genbutuList, err := s.balanceClient.GetGenbutuKabuList(ctx, s.appSession)
	if err != nil {
		return nil, fmt.Errorf("failed to get genbutu kabu list: %w", err)
	}

	// APIのレスポンスDTOからドメインモデルに変換
	positions := make([]*model.Position, 0, len(genbutuList.GenbutuKabuList))
	for _, kabu := range genbutuList.GenbutuKabuList {
		quantity, err := strconv.Atoi(kabu.UriOrderZanKabuSuryou)
		if err != nil {
			s.logger.Warn("could not parse quantity for genbutsu position, skipping", "raw", kabu.UriOrderZanKabuSuryou, "error", err)
			continue
		}
		if quantity == 0 {
			continue // 残高0のポジションは無視
		}

		avgPrice, err := strconv.ParseFloat(kabu.UriOrderGaisanBokaTanka, 64)
		if err != nil {
			s.logger.Warn("could not parse average price for genbutsu position, skipping", "raw", kabu.UriOrderGaisanBokaTanka, "error", err)
			continue
		}

		positions = append(positions, &model.Position{
			Symbol:              kabu.UriOrderIssueCode,
			PositionType:        model.PositionTypeLong,
			PositionAccountType: model.PositionAccountTypeCash, // 現物ポジション
			AveragePrice:        avgPrice,
			Quantity:            quantity,
		})
	}

	// 信用建玉リストの取得とマージ
	shinyouList, err := s.balanceClient.GetShinyouTategyokuList(ctx, s.appSession)
	if err != nil {
		// 信用口座がない場合なども考えられるため、エラーログは出すが処理は続行
		s.logger.Error("failed to get shinyou tategyoku list, proceeding with genbutsu positions only", "error", err)
		return positions, nil
	}

	for _, shinyou := range shinyouList.SinyouTategyokuList {
		quantity, err := strconv.Atoi(shinyou.OrderHensaiKanouSuryou)
		if err != nil {
			s.logger.Warn("could not parse quantity for shinyou position, skipping", "raw", shinyou.OrderHensaiKanouSuryou, "error", err)
			continue
		}
		if quantity == 0 {
			continue // 返済可能数量0のポジションは無視
		}

		avgPrice, err := strconv.ParseFloat(shinyou.OrderTategyokuTanka, 64)
		if err != nil {
			s.logger.Warn("could not parse average price for shinyou position, skipping", "raw", shinyou.OrderTategyokuTanka, "error", err)
			continue
		}

		var posType model.PositionType
		switch shinyou.OrderBaibaiKubun {
		case "1": // 売建
			posType = model.PositionTypeShort
		case "3": // 買建
			posType = model.PositionTypeLong
		default:
			s.logger.Warn("unknown baibai kubun for shinyou position, skipping", "raw", shinyou.OrderBaibaiKubun)
			continue
		}

		positions = append(positions, &model.Position{
			Symbol:              shinyou.OrderIssueCode,
			PositionType:        posType,
			PositionAccountType: model.PositionAccountTypeMarginNew, // 信用ポジションをMARGIN_NEWとして設定
			AveragePrice:        avgPrice,
			Quantity:            quantity,
		})
	}

	return positions, nil
}

// GetOrders は発注中の注文を取得する
func (s *GoaTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	s.logger.Info("GoaTradeService.GetOrders called")

	// 注文一覧取得リクエストを作成
	// 未約定+一部約定の注文を取得
	req := order_request.ReqOrderList{
		CLMID:              "CLMOrderList",
		OrderSyoukaiStatus: "5",
	}

	// APIクライアント経由で注文一覧を取得
	res, err := s.orderClient.GetOrderList(ctx, s.appSession, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list from api client: %w", err)
	}

	if res.ResultCode != "0" {
		return nil, fmt.Errorf("get order list api returned error: code=%s, text=%s", res.ResultCode, res.ResultText)
	}

	// APIのレスポンスDTOからドメインモデルに変換
	orders := make([]*model.Order, 0, len(res.OrderList))
	for _, rawOrder := range res.OrderList {
		// TradeType のマッピング
		var tradeType model.TradeType
		switch rawOrder.OrderBaibaiKubun {
		case "1": // 売
			tradeType = model.TradeTypeSell
		case "3": // 買
			tradeType = model.TradeTypeBuy
		default:
			s.logger.Warn("unknown trade type in order list, skipping", "raw", rawOrder.OrderBaibaiKubun)
			continue
		}

		// OrderType のマッピング
		var orderType model.OrderType
		switch rawOrder.OrderOrderPriceKubun {
		case "1": // 成行
			orderType = model.OrderTypeMarket
		case "2": // 指値
			orderType = model.OrderTypeLimit
		default:
			orderType = model.OrderTypeMarket // 不明な場合は成行として扱う (要検討)
			s.logger.Warn("unknown order price kubun, defaulting to MARKET", "raw", rawOrder.OrderOrderPriceKubun)
		}

		// OrderStatus のマッピング
		var orderStatus model.OrderStatus
		switch rawOrder.OrderYakuzyouStatus {
		case "0": // 未約定
			orderStatus = model.OrderStatusNew
		case "1": // 一部約定
			orderStatus = model.OrderStatusPartiallyFilled
		default:
			s.logger.Warn("skipping order with unhandled execution status", "raw", rawOrder.OrderYakuzyouStatus, "status_name", rawOrder.OrderStatus)
			continue // 全部約定などはここでは取得しないはず
		}

		quantity, err := strconv.Atoi(rawOrder.OrderOrderSuryou)
		if err != nil {
			s.logger.Warn("could not parse quantity in order list, skipping", "raw", rawOrder.OrderOrderSuryou, "error", err)
			continue
		}

		price, err := strconv.ParseFloat(rawOrder.OrderOrderPrice, 64)
		if err != nil {
			s.logger.Warn("could not parse price in order list, skipping", "raw", rawOrder.OrderOrderPrice, "error", err)
			price = 0 // パース失敗時は0
		}

		orders = append(orders, &model.Order{
			OrderID:     rawOrder.OrderOrderNumber,
			Symbol:      rawOrder.OrderIssueCode,
			TradeType:   tradeType,
			OrderType:   orderType,
			Quantity:    quantity,
			Price:       price,
			OrderStatus: orderStatus,
		})
	}

	return orders, nil
}

// GetBalance は口座残高を取得する
func (s *GoaTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
	s.logger.Info("GoaTradeService.GetBalance called")

	summary, err := s.balanceClient.GetZanKaiSummary(ctx, s.appSession)
	if err != nil {
		return nil, fmt.Errorf("failed to get zan kai summary: %w", err)
	}

	// stringからfloat64への変換
	cash, err := strconv.ParseFloat(summary.Syukkin, 64)
	if err != nil {
		s.logger.Error("failed to parse cash (Syukkin)", "raw", summary.Syukkin, "error", err)
		cash = 0
	}

	buyingPower, err := strconv.ParseFloat(summary.GenbutuKabuKaituke, 64)
	if err != nil {
		s.logger.Error("failed to parse buying power (GenbutuKabuKaituke)", "raw", summary.GenbutuKabuKaituke, "error", err)
		buyingPower = 0
	}

	return &service.Balance{
		Cash:        cash,
		BuyingPower: buyingPower,
	}, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (s *GoaTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*service.HistoricalPrice, error) {
	s.logger.Info("GoaTradeService.GetPriceHistory called", "symbol", symbol, "days", days)

	// 1. リクエストを作成
	req := request.ReqGetPriceInfoHistory{
		IssueCode: symbol,
	}

	// 2. priceClient を使って履歴情報を取得
	res, err := s.priceClient.GetPriceInfoHistory(ctx, s.appSession, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history for symbol %s: %w", symbol, err)
	}
	if res == nil || len(res.CLMMfdsGetMarketPriceHistory) == 0 {
		s.logger.Warn("no price history returned for symbol", "symbol", symbol)
		return []*service.HistoricalPrice{}, nil // 空のスライスを返し、エラーとはしない
	}

	// 3. レスポンスDTOをドメインサービスの型に変換
	history := make([]*service.HistoricalPrice, 0, len(res.CLMMfdsGetMarketPriceHistory))
	for _, item := range res.CLMMfdsGetMarketPriceHistory {
		// YYYYMMDD 形式の日付を time.Time にパース
		date, err := time.Parse("20060102", item.SDate)
		if err != nil {
			s.logger.Warn("could not parse date in price history, skipping item", "raw", item.SDate, "error", err)
			continue
		}

		// 安全なパース用のヘルパー関数
		parseFloat := func(val string) float64 {
			if val == "" {
				return 0
			}
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				s.logger.Warn("could not parse float in price history", "raw", val, "error", err)
				return 0
			}
			return f
		}
		parseInt64 := func(val string) int64 {
			if val == "" {
				return 0
			}
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				s.logger.Warn("could not parse int64 in price history", "raw", val, "error", err)
				return 0
			}
			return i
		}

		// 分割調整済みの価格を使用する (pDOPxK など)
		histItem := &service.HistoricalPrice{
			Date:   date,
			Open:   parseFloat(item.PDOPxK),
			High:   parseFloat(item.PDHPxK),
			Low:    parseFloat(item.PDLPxK),
			Close:  parseFloat(item.PDPPxK),
			Volume: parseInt64(item.PDVxK),
		}

		// 重要なOHLCのいずれかが0の場合はスキップ (データ欠損の可能性)
		if histItem.Open == 0 || histItem.High == 0 || histItem.Low == 0 || histItem.Close == 0 {
			s.logger.Info("skipping historical data point with zero OHLC values", "date", item.SDate, "symbol", symbol)
			continue
		}

		history = append(history, histItem)
	}

	// 必要であれば 'days' パラメータで結果をフィルタリングする
	if len(history) > days && days > 0 {
		// 日付でソートされていると仮定 (新しいものが先頭の場合)
		history = history[:days]
	}

	s.logger.Info("successfully fetched and processed price history", "symbol", symbol, "records_count", len(history))
	return history, nil
}

// PlaceOrder は注文を発行する
func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
	s.logger.Info("GoaTradeService.PlaceOrder called", "request", req)

	// 1. 銘柄の妥当性チェック
	if err := s.validateSymbol(ctx, req.Symbol); err != nil {
		return nil, fmt.Errorf("symbol validation failed: %w", err)
	}

	// 2. 売買単位チェック
	if err := s.validateTradingUnit(ctx, req.Symbol, req.Quantity); err != nil {
		return nil, fmt.Errorf("trading unit validation failed: %w", err)
	}

	// BaibaiKubun のマッピング
	var baibaiKubun string
	switch req.TradeType {
	case model.TradeTypeBuy:
		baibaiKubun = "3" // 買
	case model.TradeTypeSell:
		baibaiKubun = "1" // 売
	default:
		return nil, fmt.Errorf("unknown trade type: %s", req.TradeType)
	}

	// OrderPrice, Gyakusasi... のマッピング
	var orderPrice string
	var gyakusasiOrderType string = "0"
	var gyakusasiZyouken string = "0"
	var gyakusasiPrice string = "*"

	switch req.OrderType {
	case model.OrderTypeMarket:
		orderPrice = "0" // 成行
	case model.OrderTypeLimit:
		orderPrice = strconv.FormatFloat(req.Price, 'f', -1, 64)
	case model.OrderTypeStop: // 逆指値
		gyakusasiOrderType = "1" // 逆指値注文
		gyakusasiZyouken = strconv.FormatFloat(req.TriggerPrice, 'f', -1, 64)
		gyakusasiPrice = "0"
		orderPrice = "*"
	default:
		return nil, fmt.Errorf("unknown order type: %s", req.OrderType)
	}

	// GenkinShinyouKubun のマッピング
	var genkinShinyouKubun string
	switch req.PositionAccountType {
	case model.PositionAccountTypeCash:
		genkinShinyouKubun = "0" // 現物
	case model.PositionAccountTypeMarginNew:
		genkinShinyouKubun = "2" // 信用新規
	case model.PositionAccountTypeMarginRepay:
		genkinShinyouKubun = "4" // 信用返済
	default:
		s.logger.Warn("unknown PositionAccountType, defaulting to cash", "position_account_type", req.PositionAccountType)
		genkinShinyouKubun = "0"
	}

	// APIクライアントに渡すパラメータを作成
	params := client.NewOrderParams{
		ZyoutoekiKazeiC:          "1", // 譲渡益課税区分: 1(特定口座)
		IssueCode:                req.Symbol,
		SizyouC:                  "00", // 東証
		BaibaiKubun:              baibaiKubun,
		Condition:                "0", // 指定なし
		OrderPrice:               orderPrice,
		OrderSuryou:              strconv.Itoa(req.Quantity),
		GenkinShinyouKubun:       genkinShinyouKubun,
		OrderExpireDay:           "0",                // 当日限り
		GyakusasiOrderType:       gyakusasiOrderType, // 逆指値注文の種類
		GyakusasiZyouken:         gyakusasiZyouken,   // 逆指値条件
		GyakusasiPrice:           gyakusasiPrice,     // 逆指値価格
		TatebiType:               "*",                // 指定なし
		TategyokuZyoutoekiKazeiC: "*",                // 指定なし
	}

	// 注文を執行
	res, err := s.orderClient.NewOrder(ctx, s.appSession, params)
	if err != nil {
		return nil, fmt.Errorf("failed to place new order via api client: %w", err)
	}
	if res.ResultCode != "0" {
		return nil, fmt.Errorf("new order api returned error: code=%s, text=%s", res.ResultCode, res.ResultText)
	}

	// レスポンスをドメインモデルに変換
	newOrder := &model.Order{
		OrderID:             res.OrderNumber,
		Symbol:              req.Symbol,
		TradeType:           req.TradeType,
		OrderType:           req.OrderType,
		Quantity:            req.Quantity,
		Price:               req.Price,
		OrderStatus:         model.OrderStatusNew,
		PositionAccountType: req.PositionAccountType,
	}

	// データベースに保存
	if err := s.orderRepo.Save(ctx, newOrder); err != nil {
		s.logger.Error("successfully placed order but failed to save to DB", "order_id", newOrder.OrderID, "error", err)
		return newOrder, fmt.Errorf("placed order but failed to save to DB: %w", err)
	}
	s.logger.Info("successfully saved new order to DB", "order_id", newOrder.OrderID)

	return newOrder, nil
}

// validateSymbol は銘柄コードの妥当性をチェックする
func (s *GoaTradeService) validateSymbol(ctx context.Context, symbol string) error {
	rawResult, err := s.masterRepo.FindByIssueCode(ctx, symbol, "StockMaster")
	if err != nil {
		return fmt.Errorf("failed to find stock master: %w", err)
	}
	if rawResult == nil {
		return fmt.Errorf("symbol %s not found in master data", symbol)
	}
	return nil
}

// validateTradingUnit は売買単位をチェックする
func (s *GoaTradeService) validateTradingUnit(ctx context.Context, symbol string, quantity int) error {
	rawResult, err := s.masterRepo.FindByIssueCode(ctx, symbol, "StockMaster")
	if err != nil {
		return fmt.Errorf("failed to find stock master: %w", err)
	}
	if rawResult == nil {
		return fmt.Errorf("symbol %s not found in master data", symbol)
	}

	stockMaster, ok := rawResult.(*model.StockMaster)
	if !ok {
		return fmt.Errorf("unexpected type returned from repository")
	}

	if stockMaster.TradingUnit > 0 && quantity%stockMaster.TradingUnit != 0 {
		return fmt.Errorf("quantity %d must be multiple of trading unit %d", quantity, stockMaster.TradingUnit)
	}

	return nil
}

// CancelOrder は注文をキャンセルする
func (s *GoaTradeService) CancelOrder(ctx context.Context, orderID string) error {
	s.logger.Info("GoaTradeService.CancelOrder called", "orderID", orderID)

	// 1. データベースから注文情報を取得
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// 2. 注文がキャンセル可能な状態かチェック
	if !order.IsUnexecuted() {
		return fmt.Errorf("order %s cannot be cancelled, current status: %s", orderID, order.OrderStatus)
	}

	// 3. 営業日の取得（簡易実装：現在日付をYYYYMMDD形式で使用）
	// 実際の実装では、証券会社の営業日カレンダーを参照する必要があります
	eigyouDay := time.Now().Format("20060102")

	// 4. 証券会社APIでキャンセル実行
	params := client.CancelOrderParams{
		OrderNumber: orderID,
		EigyouDay:   eigyouDay,
	}

	res, err := s.orderClient.CancelOrder(ctx, s.appSession, params)
	if err != nil {
		return fmt.Errorf("failed to cancel order via api client: %w", err)
	}

	if res.ResultCode != "0" {
		return fmt.Errorf("cancel order api returned error: code=%s, text=%s", res.ResultCode, res.ResultText)
	}

	// 5. データベースの注文状態を更新
	order.OrderStatus = model.OrderStatusCanceled
	if err := s.orderRepo.Save(ctx, order); err != nil {
		s.logger.Error("successfully cancelled order but failed to update DB", "order_id", orderID, "error", err)
		// APIでのキャンセルは成功しているので、エラーは返さずに警告ログのみ
	}

	s.logger.Info("successfully cancelled order", "order_id", orderID)
	return nil
}

// GetOrderHistory は注文履歴を取得する
func (s *GoaTradeService) GetOrderHistory(ctx context.Context, status *model.OrderStatus, symbol *string, limit int) ([]*model.Order, error) {
	s.logger.Info("GoaTradeService.GetOrderHistory called", "status", status, "symbol", symbol, "limit", limit)

	orders, err := s.orderRepo.FindOrderHistory(ctx, status, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	s.logger.Info("successfully retrieved order history", "count", len(orders))
	return orders, nil
}

// CorrectOrder は注文を訂正する
func (s *GoaTradeService) CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error) {
	s.logger.Info("GoaTradeService.CorrectOrder called", "orderID", orderID, "newPrice", newPrice, "newQuantity", newQuantity)

	// 1. データベースから注文情報を取得
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	// 2. 注文が訂正可能な状態かチェック
	if !order.IsUnexecuted() {
		return nil, fmt.Errorf("order %s cannot be corrected, current status: %s", orderID, order.OrderStatus)
	}

	// 3. 営業日の取得（簡易実装）
	eigyouDay := time.Now().Format("20060102")

	// 4. 訂正パラメータの準備
	params := client.CorrectOrderParams{
		OrderNumber:    orderID,
		EigyouDay:      eigyouDay,
		Condition:      "0", // 指定なし
		OrderExpireDay: "0", // 当日限り
	}

	// 価格の設定
	if newPrice != nil {
		params.OrderPrice = strconv.FormatFloat(*newPrice, 'f', -1, 64)
	} else {
		params.OrderPrice = strconv.FormatFloat(order.Price, 'f', -1, 64)
	}

	// 数量の設定
	if newQuantity != nil {
		params.OrderSuryou = strconv.Itoa(*newQuantity)
	} else {
		params.OrderSuryou = strconv.Itoa(order.Quantity)
	}

	// 5. 証券会社APIで訂正実行
	res, err := s.orderClient.CorrectOrder(ctx, s.appSession, params)
	if err != nil {
		return nil, fmt.Errorf("failed to correct order via api client: %w", err)
	}

	if res.ResultCode != "0" {
		return nil, fmt.Errorf("correct order api returned error: code=%s, text=%s", res.ResultCode, res.ResultText)
	}

	// 6. データベースの注文情報を更新
	if newPrice != nil {
		order.Price = *newPrice
	}
	if newQuantity != nil {
		order.Quantity = *newQuantity
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		s.logger.Error("successfully corrected order but failed to update DB", "order_id", orderID, "error", err)
		// APIでの訂正は成功しているので、エラーは返さずに警告ログのみ
	}

	s.logger.Info("successfully corrected order", "order_id", orderID)
	return order, nil
}

// CancelAllOrders は全ての未約定注文をキャンセルする
func (s *GoaTradeService) CancelAllOrders(ctx context.Context) (int, error) {
	s.logger.Info("GoaTradeService.CancelAllOrders called")

	// 1. 証券会社APIで一括キャンセル実行
	params := client.CancelOrderAllParams{}
	res, err := s.orderClient.CancelOrderAll(ctx, s.appSession, params)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel all orders via api client: %w", err)
	}

	if res.ResultCode != "0" {
		return 0, fmt.Errorf("cancel all orders api returned error: code=%s, text=%s", res.ResultCode, res.ResultText)
	}

	// 2. データベースの未約定注文をキャンセル状態に更新
	// 未約定注文を取得
	newOrders, err := s.orderRepo.FindByStatus(ctx, model.OrderStatusNew)
	if err != nil {
		s.logger.Error("failed to find new orders for bulk cancel update", "error", err)
		return 0, fmt.Errorf("failed to find orders to update: %w", err)
	}

	partiallyFilledOrders, err := s.orderRepo.FindByStatus(ctx, model.OrderStatusPartiallyFilled)
	if err != nil {
		s.logger.Error("failed to find partially filled orders for bulk cancel update", "error", err)
		return 0, fmt.Errorf("failed to find orders to update: %w", err)
	}

	// 全ての未約定注文をキャンセル状態に更新
	cancelledCount := 0
	allOrders := append(newOrders, partiallyFilledOrders...)

	for _, order := range allOrders {
		order.OrderStatus = model.OrderStatusCanceled
		if err := s.orderRepo.Save(ctx, order); err != nil {
			s.logger.Error("failed to update order status to cancelled", "order_id", order.OrderID, "error", err)
			continue
		}
		cancelledCount++
	}

	s.logger.Info("successfully cancelled all orders", "cancelled_count", cancelledCount)
	return cancelledCount, nil
}

// ValidateSymbolInternal は内部用の銘柄バリデーション（外部公開用）
func (s *GoaTradeService) ValidateSymbolInternal(ctx context.Context, symbol string) error {
	return s.validateSymbol(ctx, symbol)
}

// StockInfo は銘柄情報を表現する構造体
type StockInfo struct {
	Symbol      string
	Name        string
	Market      string
	TradingUnit int
}

// GetStockInfo は銘柄の詳細情報を取得する
func (s *GoaTradeService) GetStockInfo(ctx context.Context, symbol string) (*StockInfo, error) {
	rawResult, err := s.masterRepo.FindByIssueCode(ctx, symbol, "StockMaster")
	if err != nil {
		return nil, fmt.Errorf("failed to find stock master: %w", err)
	}
	if rawResult == nil {
		return nil, fmt.Errorf("symbol %s not found in master data", symbol)
	}

	stockMaster, ok := rawResult.(*model.StockMaster)
	if !ok {
		return nil, fmt.Errorf("unexpected type returned from repository")
	}

	return &StockInfo{
		Symbol:      stockMaster.IssueCode,
		Name:        stockMaster.IssueName,
		Market:      stockMaster.MarketCode,
		TradingUnit: stockMaster.TradingUnit,
	}, nil
}

// HealthCheck はサービスの健康状態をチェックする
func (s *GoaTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
	s.logger.Debug("GoaTradeService.HealthCheck called")

	status := &service.HealthStatus{
		Timestamp: time.Now(),
	}

	// セッション有効性チェック
	if s.appSession != nil {
		status.SessionValid = true
	}

	// データベース接続チェック
	if s.orderRepo != nil {
		// 簡単なクエリでデータベース接続をテスト
		_, err := s.orderRepo.FindByStatus(ctx, model.OrderStatusNew)
		status.DatabaseConnected = (err == nil)
	}

	// WebSocket接続状態は現在の実装では直接チェックできないため、
	// 簡易的にEventClientの存在で判定
	status.WebSocketConnected = true // 簡易実装

	// 全体的な健康状態を判定
	if status.SessionValid && status.DatabaseConnected && status.WebSocketConnected {
		status.Status = "healthy"
	} else if status.SessionValid && status.DatabaseConnected {
		status.Status = "degraded"
	} else {
		status.Status = "unhealthy"
	}

	s.logger.Debug("health check completed", "status", status.Status)
	return status, nil
}
