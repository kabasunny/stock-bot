package state

import (
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"sync"
)

// State はエージェントの内部状態を管理する
// 全てのフィールドへのアクセスはスレッドセーフである必要がある
type State struct {
	mutex     sync.RWMutex
	positions map[string]*model.Position // キーは銘柄コード(Symbol)
	orders    map[string]*model.Order    // キーは証券会社の注文ID(OrderID)
	prices    map[string]float64         // キーは銘柄コード(Symbol), 値は現在の価格
	balance   *service.Balance
}

// NewState は新しいStateを初期化して返す
func NewState() *State {
	return &State{
		positions: make(map[string]*model.Position),
		orders:    make(map[string]*model.Order),
		prices:    make(map[string]float64),
		balance:   &service.Balance{},
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
	for _, pos := range positions {
		newPositions[pos.Symbol] = pos
	}
	s.positions = newPositions
}

// GetPositions は現在の保有ポジション一覧を取得する
func (s *State) GetPositions() []*model.Position {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	positions := make([]*model.Position, 0, len(s.positions))
	for _, pos := range s.positions {
		positions = append(positions, pos)
	}
	return positions
}

// GetPosition は指定した銘柄のポジションを取得する
func (s *State) GetPosition(symbol string) (*model.Position, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	pos, ok := s.positions[symbol]
	return pos, ok
}

// UpdateOrders は注文情報を更新する
func (s *State) UpdateOrders(orders []*model.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newOrders := make(map[string]*model.Order)
	for _, order := range orders {
		newOrders[order.OrderID] = order
	}
	s.orders = newOrders
}

// GetOrders は現在の注文一覧を取得する
func (s *State) GetOrders() []*model.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	orders := make([]*model.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

// AddOrder は新しい注文を追加する
func (s *State) AddOrder(order *model.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.orders[order.OrderID] = order
}

// UpdateBalance は残高情報を更新する
func (s *State) UpdateBalance(balance *service.Balance) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.balance = balance
}

// GetBalance は現在の残高情報を取得する
func (s *State) GetBalance() *service.Balance {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.balance
}

// UpdatePositionHighestPrice は指定したポジションの最高価格を更新する
func (s *State) UpdatePositionHighestPrice(symbol string, highestPrice float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if pos, ok := s.positions[symbol]; ok {
		pos.HighestPrice = highestPrice
	}
}

// UpdatePositionTrailingStopPrice は指定したポジションのトレーリングストップ価格を更新する
func (s *State) UpdatePositionTrailingStopPrice(symbol string, trailingStopPrice float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if pos, ok := s.positions[symbol]; ok {
		pos.TrailingStopPrice = trailingStopPrice
	}
}
