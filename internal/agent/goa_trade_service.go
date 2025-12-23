package agent

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"strconv"
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

// GetPositions は現在の保有ポジションを取得する
func (s *GoaTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	s.logger.Info("GoaTradeService.GetPositions called")

	// balanceClient を使って現物保有銘柄リストを取得
	genbutuList, err := s.balanceClient.GetGenbutuKabuList(ctx, s.appSession)
	if err != nil {
		return nil, fmt.Errorf("failed to get genbutu kabu list: %w", err)
	}

	// APIのレスポンスDTOからドメインモデルに変換
	// パースできないレコードはスキップするため、可変長の positions スライスを準備
	positions := make([]*model.Position, 0, len(genbutuList.GenbutuKabuList))
	for _, kabu := range genbutuList.GenbutuKabuList {
		quantity, err := strconv.Atoi(kabu.UriOrderZanKabuSuryou)
		if err != nil {
			s.logger.Warn("could not parse quantity, skipping position record", "raw", kabu.UriOrderZanKabuSuryou, "error", err)
			continue
		}
		if quantity == 0 {
			continue // 残高0のポジションは無視
		}

		avgPrice, err := strconv.ParseFloat(kabu.UriOrderGaisanBokaTanka, 64)
		if err != nil {
			s.logger.Warn("could not parse average price, skipping position record", "raw", kabu.UriOrderGaisanBokaTanka, "error", err)
			continue
		}

		positions = append(positions, &model.Position{
			Symbol:       kabu.UriOrderIssueCode,
			PositionType: model.PositionTypeLong, // 現物はLONG
			AveragePrice: avgPrice,
			Quantity:     quantity,
		})
	}

	// TODO: 信用建玉も取得してマージする必要がある

	return positions, nil
}


// GetOrders は発注中の注文を取得する
func (s *GoaTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	s.logger.Info("GoaTradeService.GetOrders called")
    // TODO: orderClient を使って発注中注文を取得し、model.Order に変換する
	return []*model.Order{}, nil // ダミー実装
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

	// OrderPrice のマッピング
	var orderPrice string
	switch req.OrderType {
	case model.OrderTypeMarket:
		orderPrice = "0" // 成行
	case model.OrderTypeLimit:
		orderPrice = strconv.FormatFloat(req.Price, 'f', -1, 64)
	default:
		return nil, fmt.Errorf("unknown order type: %s", req.OrderType)
	}

	// APIクライアントに渡すパラメータを作成
	params := client.NewOrderParams{
		IssueCode:          req.Symbol,
		SizyouC:            "00", // 東証
		BaibaiKubun:        baibaiKubun,
		Condition:          "0", // 指定なし
		OrderPrice:         orderPrice,
		OrderSuryou:        strconv.Itoa(req.Quantity),
		GenkinShinyouKubun: "0", // 現物
		OrderExpireDay:     "0", // 当日
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
