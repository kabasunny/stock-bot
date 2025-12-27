package agent

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/infrastructure/client"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"strconv"
	"time"
	// "stock-bot/internal/infrastructure/client/dto/balance/request"
)

// GoaTradeService は TradeService インターフェースのGoaクライアント実装
type GoaTradeService struct {
	balanceClient client.BalanceClient
	orderClient   client.OrderClient
	priceClient   client.PriceInfoClient
	orderRepo     repository.OrderRepository
	appSession    *client.Session
	logger        *slog.Logger
}

// NewGoaTradeService は GoaTradeService の新しいインスタンスを作成する
func NewGoaTradeService(
	balanceClient client.BalanceClient,
	orderClient client.OrderClient,
	priceClient client.PriceInfoClient,
	orderRepo repository.OrderRepository,
	appSession *client.Session,
	logger *slog.Logger,
) *GoaTradeService {
	return &GoaTradeService{
		balanceClient: balanceClient,
		orderClient:   orderClient,
		priceClient:   priceClient,
		orderRepo:     orderRepo,
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
			Symbol:       kabu.UriOrderIssueCode,
			PositionType: model.PositionTypeLong, // 現物はLONG
			AveragePrice: avgPrice,
			Quantity:     quantity,
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
			Symbol:       shinyou.OrderIssueCode,
			PositionType: posType,
			AveragePrice: avgPrice,
			Quantity:     quantity,
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
func (s *GoaTradeService) GetBalance(ctx context.Context) (*Balance, error) {
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

	agentBalance := &Balance{
		Cash:        cash,
		BuyingPower: buyingPower,
	}
	return agentBalance, nil
}

// GetPrice は指定した銘柄の現在価格を取得する
func (s *GoaTradeService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	s.logger.Info("GoaTradeService.GetPrice called", "symbol", symbol)

	// リクエストを作成
	req := request.ReqGetPriceInfo{
		CLMID:           "CLMMfdsGetMarketPrice",
		TargetIssueCode: symbol,
		TargetColumn:    "CurrentPrice", // 現在値のみ取得
	}

	// priceClient を使って価格情報を取得
	res, err := s.priceClient.GetPriceInfo(ctx, s.appSession, req)
	if err != nil {
		return 0, fmt.Errorf("failed to get price info for symbol %s: %w", symbol, err)
	}

	// レスポンスをパースして価格を取得
	if res == nil || len(res.CLMMfdsMarketPrice) == 0 {
		return 0, fmt.Errorf("no price info returned for symbol %s", symbol)
	}

	item := res.CLMMfdsMarketPrice[0] // 最初のアイテムを使用

	priceStr, ok := item.Values["CurrentPrice"]
	if !ok {
		return 0, fmt.Errorf("CurrentPrice not found in response for symbol %s", symbol)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price '%s' for symbol %s: %w", priceStr, symbol, err)
	}

	return price, nil
}

// GetPriceHistory は指定した銘柄の過去の価格情報を取得する
func (s *GoaTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error) {
	s.logger.Info("GoaTradeService.GetPriceHistory called", "symbol", symbol, "days", days)

	// 1. リクエストを作成
	req := request.ReqGetPriceInfoHistory{
		IssueCode: symbol,
		// SizyouC: "00", // 東証 (デフォルトのため省略可能)
	}

	// 2. priceClient を使って履歴情報を取得
	res, err := s.priceClient.GetPriceInfoHistory(ctx, s.appSession, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history for symbol %s: %w", symbol, err)
	}
	if res == nil || len(res.CLMMfdsGetMarketPriceHistory) == 0 {
		s.logger.Warn("no price history returned for symbol", "symbol", symbol)
		return []*HistoricalPrice{}, nil // 空のスライスを返し、エラーとはしない
	}

	// 3. レスポンスDTOをエージェントの型に変換
	history := make([]*HistoricalPrice, 0, len(res.CLMMfdsGetMarketPriceHistory))
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
		parseInt := func(val string) int {
			if val == "" {
				return 0
			}
			i, err := strconv.Atoi(val)
			if err != nil {
				s.logger.Warn("could not parse int in price history", "raw", val, "error", err)
				return 0
			}
			return i
		}

		// 分割調整済みの価格を使用する (pDOPxK など)
		histItem := &HistoricalPrice{
			Date:   date,
			Open:   parseFloat(item.PDOPxK),
			High:   parseFloat(item.PDHPxK),
			Low:    parseFloat(item.PDLPxK),
			Close:  parseFloat(item.PDPPxK),
			Volume: parseInt(item.PDVxK),
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
func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error) {
	s.logger.Info("GoaTradeService.PlaceOrder called", "request", req)

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
		// ドキュメントとテストケースに基づき、GyakusasiZyouken がトリガー価格
		gyakusasiZyouken = strconv.FormatFloat(req.TriggerPrice, 'f', -1, 64)
		// 逆指値成行注文なので、GyakusasiPrice は "0" (成行)
		gyakusasiPrice = "0"
		// この種の注文では、メインの OrderPrice は "*" に設定する必要がある
		orderPrice = "*"
	default:
		return nil, fmt.Errorf("unknown order type: %s", req.OrderType)
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
		GenkinShinyouKubun:       "0",                // 現物
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
		OrderID:     res.OrderNumber,
		Symbol:      req.Symbol,
		TradeType:   req.TradeType,
		OrderType:   req.OrderType,
		Quantity:    req.Quantity,
		Price:       req.Price,
		OrderStatus: model.OrderStatusNew,
		// TimeInForce はgormのデフォルト値'DAY'に任せる
	}

	// データベースに保存
	if err := s.orderRepo.Save(ctx, newOrder); err != nil {
		// APIでの注文は成功しているがDB保存に失敗した場合。
		// 本来はより堅牢なエラーハンドリング（リトライ、ログ、通知など）が必要。
		s.logger.Error("successfully placed order but failed to save to DB", "order_id", newOrder.OrderID, "error", err)
		return newOrder, fmt.Errorf("placed order but failed to save to DB: %w", err)
	}
	s.logger.Info("successfully saved new order to DB", "order_id", newOrder.OrderID)

	return newOrder, nil
}

// CancelOrder は注文をキャンセルする
func (s *GoaTradeService) CancelOrder(ctx context.Context, orderID string) error {
	// TODO: orderClient を使って注文をキャンセルする
	s.logger.Info("GoaTradeService.CancelOrder called", "orderID", orderID)
	return fmt.Errorf("CancelOrder not implemented") // ダミー実装
}
