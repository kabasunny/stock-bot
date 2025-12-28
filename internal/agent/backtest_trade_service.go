package agent

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/internal/data"
	"stock-bot/internal/infrastructure/client"
	"time"

	"gorm.io/gorm"
)

// BacktestTradeService はバックテスト用のTradeServiceインターフェースの実装です。
// 実際の証券会社APIとの通信の代わりに、履歴データと内部状態を基に動作します。
type BacktestTradeService struct {
	historyReader  *data.PriceHistoryReader
	AllHistory     map[string][]*data.HistoricalPrice // 銘柄ごとの全履歴データ
	CurrentTick    time.Time                          // シミュレーション現在の時刻
	Balance        *Balance                           // 擬似的な口座残高
	Positions      map[string]*model.Position         // 擬似的な保有ポジション
	Orders         map[string]*model.Order            // 擬似的な注文情報
	OrderIDCounter int
}

// NewBacktestTradeService は新しいBacktestTradeServiceのインスタンスを作成します。
func NewBacktestTradeService(dataDir string, initialCash float64) *BacktestTradeService {
	return &BacktestTradeService{
		historyReader: data.NewPriceHistoryReader(dataDir),
		AllHistory:    make(map[string][]*data.HistoricalPrice),
		Balance: &Balance{
			Cash:        initialCash,
			BuyingPower: initialCash,
		},
		Positions:      make(map[string]*model.Position),
		Orders:         make(map[string]*model.Order),
		OrderIDCounter: 0,
	}
}

// LoadHistory はバックテストに必要なすべての履歴データを事前に読み込みます。
func (s *BacktestTradeService) LoadHistory(symbols []string) error {
	for _, symbol := range symbols {
		hist, err := s.historyReader.ReadHistory(symbol)
		if err != nil {
			return fmt.Errorf("failed to load history for %s: %w", symbol, err)
		}
		s.AllHistory[symbol] = hist
	}
	return nil
}

// SetCurrentTick はシミュレーションの現在時刻を設定します。
func (s *BacktestTradeService) SetCurrentTick(t time.Time) {
	s.CurrentTick = t
}

// GetSession はバックテストでは使用しないためnilを返します。
func (s *BacktestTradeService) GetSession() *client.Session {
	return nil
}

// GetPositions は現在の擬似的な保有ポジションを返します。
func (s *BacktestTradeService) GetPositions(ctx context.Context) ([]*model.Position, error) {
	// deep copy to prevent external modification
	posList := make([]*model.Position, 0, len(s.Positions))
	for _, p := range s.Positions {
		posList = append(posList, &model.Position{
			Model: gorm.Model{
				ID: p.ID,
			},
			Symbol:       p.Symbol,
			PositionType: p.PositionType,
			AveragePrice: p.AveragePrice,
			Quantity:     p.Quantity,
			HighestPrice: p.HighestPrice, // HighestPriceも返す
		})
	}
	return posList, nil
}

// GetOrders は現在の擬似的な注文情報を返します。
func (s *BacktestTradeService) GetOrders(ctx context.Context) ([]*model.Order, error) {
	// deep copy
	orderList := make([]*model.Order, 0, len(s.Orders))
	for _, o := range s.Orders {
		orderList = append(orderList, &model.Order{
			OrderID:   o.OrderID,
			Symbol:    o.Symbol,
			TradeType: o.TradeType,
			Quantity:  o.Quantity,
			Price:     o.Price,
			// ... 他のフィールドも必要に応じてコピー
		})
	}
	return orderList, nil
}

// GetBalance は現在の擬似的な口座残高を返します。
func (s *BacktestTradeService) GetBalance(ctx context.Context) (*Balance, error) {
	return s.Balance, nil
}

// GetPrice は現在のシミュレーション時刻における価格を返します。
// 終値（Close）を使用します。指定された日付のデータがない場合は、その直近の過去データを返します。
func (s *BacktestTradeService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	history, ok := s.AllHistory[symbol]
	if !ok {
		return 0, fmt.Errorf("history not loaded for symbol: %s", symbol)
	}

	var latestPrice float64 = -1
	for _, h := range history {
		if h.Date.After(s.CurrentTick) {
			break // 現在のtickを超えたらループを抜ける
		}
		latestPrice = h.Close
	}

	if latestPrice == -1 {
		return 0, fmt.Errorf("no price data found for symbol %s on or before %s", symbol, s.CurrentTick.Format("2006-01-02"))
	}

	return latestPrice, nil
}

// GetPriceHistory は現在のシミュレーション時刻より過去の履歴データを返します。
// シミュレーション時刻のデータがない場合でも、その直近の過去データを基準に履歴を返します。
func (s *BacktestTradeService) GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error) {
	allHist, ok := s.AllHistory[symbol]
	if !ok {
		return nil, fmt.Errorf("history not loaded for symbol: %s", symbol)
	}

	// s.CurrentTick 以前で最も新しい履歴データの日付を探す
	effectiveIndex := -1
	// allHist は日付昇順でソートされていると仮定
	for i := 0; i < len(allHist); i++ {
		if !allHist[i].Date.After(s.CurrentTick) { // CurrentTick 以前の日付であれば
			effectiveIndex = i // その日を有効なインデックスとする
		} else {
			break // CurrentTick より後の日付になったらループを抜ける
		}
	}

	if effectiveIndex == -1 {
		return nil, fmt.Errorf("no historical data found for symbol %s on or before %s", symbol, s.CurrentTick.Format("2006-01-02"))
	}

	// 過去 days 分のデータを取得
	// effectiveIndex は 0-indexed なので effectiveIndex+1 がデータ数
	actualDays := days
	if effectiveIndex+1 < actualDays {
		actualDays = effectiveIndex + 1 // 取得可能な最大日数に制限
	}
	
	result := make([]*HistoricalPrice, actualDays)
	for i := 0; i < actualDays; i++ {
		src := allHist[effectiveIndex-i] // effectiveIndex を基準に過去に遡る
		result[actualDays-1-i] = &HistoricalPrice{ // 古い方から新しい方へソートされるように逆順にコピー
			Date:   src.Date,
			Open:   src.Open,
			High:   src.High,
			Low:    src.Low,
			Close:  src.Close,
			Volume: src.Volume,
		}
	}
	return result, nil
}

// PlaceOrder は注文を処理し、内部状態を更新します。
func (s *BacktestTradeService) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error) {
	currentPrice, err := s.GetPrice(ctx, req.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get current price for order placement: %w", err)
	}

	// 成行注文として即時約定を仮定
	executedPrice := currentPrice // スリッページなし
	cost := executedPrice * float64(req.Quantity)

	// 注文IDを生成
	s.OrderIDCounter++
	orderID := fmt.Sprintf("BT-ORDER-%d", s.OrderIDCounter)

	newOrder := &model.Order{
		OrderID:     orderID,
		Symbol:      req.Symbol,
		TradeType:   req.TradeType,
		Quantity:    req.Quantity,
		Price:       executedPrice,           // 約定価格
		OrderStatus: model.OrderStatusFilled, // 即時約定
		// ... 他のフィールドを設定
	}
	s.Orders[orderID] = newOrder

	if req.TradeType == model.TradeTypeBuy {
		if s.Balance.BuyingPower < cost {
			return nil, fmt.Errorf("insufficient funds for buy order for %s", req.Symbol)
		}
		s.Balance.BuyingPower -= cost
		s.Balance.Cash -= cost // 現金も減る

		if pos, ok := s.Positions[req.Symbol]; ok {
			// 既存ポジションに追加
			totalCost := (pos.AveragePrice * float64(pos.Quantity)) + cost
			totalQuantity := pos.Quantity + req.Quantity
			pos.AveragePrice = totalCost / float64(totalQuantity)
			pos.Quantity = totalQuantity
			pos.HighestPrice = executedPrice // 購入時の価格をHighestPriceの初期値とする
		} else {
			// 新規ポジション
			s.Positions[req.Symbol] = &model.Position{
				Symbol:       req.Symbol,
				PositionType: model.PositionTypeLong, // バックテストでは全てLONG
				AveragePrice: executedPrice,
				Quantity:     req.Quantity,
				HighestPrice: executedPrice, // 購入時の価格をHighestPriceの初期値とする
			}
		}
	} else if req.TradeType == model.TradeTypeSell {
		pos, ok := s.Positions[req.Symbol]
		if !ok || pos.Quantity < req.Quantity {
			return nil, fmt.Errorf("insufficient position for sell order for %s", req.Symbol)
		}
		s.Balance.BuyingPower += cost
		s.Balance.Cash += cost

		pos.Quantity -= req.Quantity
		if pos.Quantity == 0 {
			delete(s.Positions, req.Symbol) // ポジションがゼロになったら削除
		}
	}

	return newOrder, nil
}

// CancelOrder はバックテストでは注文は即時約定するため、実装しません。
func (s *BacktestTradeService) CancelOrder(ctx context.Context, orderID string) error {
	return fmt.Errorf("CancelOrder is not implemented for backtest")
}

// UpdateHighestPriceForPosition は、agent.Stateが呼び出すHighestPriceの更新をモックするためのヘルパーメソッドです。
// BacktestTradeServiceでは直接DB更新は行わないため、内部のポジション情報を更新します。
func (s *BacktestTradeService) UpdateHighestPriceForPosition(symbol string, price float64) {
	if pos, ok := s.Positions[symbol]; ok {
		pos.HighestPrice = price
	}
}
