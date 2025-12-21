package agent

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"strconv"
	// "stock-bot/internal/infrastructure/client/dto/balance/request"
)

// GoaTradeService は TradeService インターフェースのGoaクライアント実装
type GoaTradeService struct {
	balanceClient client.BalanceClient
	orderClient   client.OrderClient
	appSession    *client.Session
	logger        *slog.Logger
}

// NewGoaTradeService は GoaTradeService の新しいインスタンスを作成する
func NewGoaTradeService(
	balanceClient client.BalanceClient,
	orderClient client.OrderClient,
	appSession *client.Session,
	logger *slog.Logger,
) *GoaTradeService {
	return &GoaTradeService{
		balanceClient: balanceClient,
		orderClient:   orderClient,
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


// PlaceOrder は注文を発行する
func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error) {
	// TODO: orderClient を使って注文を発行する
	s.logger.Info("GoaTradeService.PlaceOrder called", "request", req)
	return nil, fmt.Errorf("PlaceOrder not implemented") // ダミー実装
}

// CancelOrder は注文をキャンセルする
func (s *GoaTradeService) CancelOrder(ctx context.Context, orderID string) error {
	// TODO: orderClient を使って注文をキャンセルする
	s.logger.Info("GoaTradeService.CancelOrder called", "orderID", orderID)
	return fmt.Errorf("CancelOrder not implemented") // ダミー実装
}
