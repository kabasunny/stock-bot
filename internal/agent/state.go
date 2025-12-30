package agent

import (
	"stock-bot/domain/model"
	"sync"
)

// Balance は口座残高の情報を保持する
type Balance struct {
	Cash        float64 // 現金残高
	BuyingPower float64 // 買付余力
}

// State はエージェントの内部状態を管理する
// 全てのフィールドへのアクセスはスレッドセーフである必要がある
type State struct {
	mutex     sync.RWMutex
	positions map[string]*model.Position // キーは銘柄コード(Symbol)
	orders    map[string]*model.Order    // キーは証券会社の注文ID(OrderID)
	prices    map[string]float64         // キーは銘柄コード(Symbol), 値は現在の価格
	balance   *Balance
}

// NewState は新しいStateを初期化して返す
func NewState() *State {
	return &State{
		positions: make(map[string]*model.Position),
		orders:    make(map[string]*model.Order),
		prices:    make(map[string]float64), // pricesマップを初期化
		balance:   &Balance{},
	}
}

// UpdatePrice は銘柄の現在価格を更新する
func (s *State) UpdatePrice(symbol string, price float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.prices[symbol] = price
}

// GetPrice は銘柄の現在価格を取得する
// 存在しない場合は(0, false)を返す
func (s *State) GetPrice(symbol string) (float64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	price, ok := s.prices[symbol]
	return price, ok
}

// UpdatePositions は保有ポジション情報を更新する
func (s *State) UpdatePositions(positions []*model.Position) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newPositions := make(map[string]*model.Position)
	for _, p := range positions {
		// 既存のポジション情報があれば、HighestPriceとTrailingStopPriceを引き継ぐ
		if existingPos, ok := s.positions[p.Symbol]; ok {
			p.HighestPrice = existingPos.HighestPrice
			p.TrailingStopPrice = existingPos.TrailingStopPrice
		}
		newPositions[p.Symbol] = p
	}
	s.positions = newPositions
}

// UpdatePositionHighestPrice は指定した銘柄のHighestPriceを更新する
func (s *State) UpdatePositionHighestPrice(symbol string, price float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if pos, ok := s.positions[symbol]; ok {
		pos.HighestPrice = price
	}
}

// UpdatePositionTrailingStopPrice は指定した銘柄のTrailingStopPriceを更新する
func (s *State) UpdatePositionTrailingStopPrice(symbol string, price float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if pos, ok := s.positions[symbol]; ok {
		pos.TrailingStopPrice = price
	}
}

// GetPosition は指定した銘柄のポジションを取得する
// 存在しない場合は(nil, false)を返す
func (s *State) GetPosition(symbol string) (*model.Position, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	pos, ok := s.positions[symbol]
	return pos, ok
}


// UpdateOrders は発注中注文の情報を更新する
func (s *State) UpdateOrders(orders []*model.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newOrders := make(map[string]*model.Order)
	for _, o := range orders {
		newOrders[o.OrderID] = o
	}
	s.orders = newOrders
}

// GetOrder は指定した注文IDの注文を取得する
// 存在しない場合は(nil, false)を返す
func (s *State) GetOrder(orderID string) (*model.Order, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	ord, ok := s.orders[orderID]
	return ord, ok
}

// AddOrder は新しい注文を一件追加する
func (s *State) AddOrder(order *model.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.orders[order.OrderID] = order
}

// UpdateBalance は口座残高の情報を更新する
func (s *State) UpdateBalance(balance *Balance) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.balance = balance
}

// GetBalance は現在の口座残高の情報を取得する
func (s *State) GetBalance() *Balance {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	// 読み取り専用で返すためにコピーを返す
	b := *s.balance
	return &b
}

// GetPositions は現在の保有ポジションのリストをコピーして取得する
func (s *State) GetPositions() []*model.Position {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	positions := make([]*model.Position, 0, len(s.positions))
	for _, p := range s.positions {
		positions = append(positions, p)
	}
	return positions
}

// GetOrders は現在の発注中注文のリストをコピーして取得する
func (s *State) GetOrders() []*model.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	orders := make([]*model.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}
	return orders
}
