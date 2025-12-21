package agent_test

import (
	"stock-bot/domain/model"
	"stock-bot/internal/agent"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_Positions(t *testing.T) {
	state := agent.NewState()

	// 初期状態の確認
	pos, ok := state.GetPosition("7203")
	assert.False(t, ok)
	assert.Nil(t, pos)

	// ポジションの更新
	positions := []*model.Position{
		{Symbol: "7203", Quantity: 100, AveragePrice: 3000},
		{Symbol: "9984", Quantity: 200, AveragePrice: 5000},
	}
	state.UpdatePositions(positions)

	// 更新後の取得確認
	pos7203, ok7203 := state.GetPosition("7203")
	assert.True(t, ok7203)
	assert.NotNil(t, pos7203)
	assert.Equal(t, 100, pos7203.Quantity)

	pos9984, ok9984 := state.GetPosition("9984")
	assert.True(t, ok9984)
	assert.NotNil(t, pos9984)
	assert.Equal(t, 200, pos9984.Quantity)

	// 存在しない銘柄の確認
	pos_none, ok_none := state.GetPosition("XXXX")
	assert.False(t, ok_none)
	assert.Nil(t, pos_none)

	// 再度更新（上書き）
	newPositions := []*model.Position{
		{Symbol: "7203", Quantity: 50, AveragePrice: 3100}, // 数量変更
	}
	state.UpdatePositions(newPositions)

	pos7203_updated, ok7203_updated := state.GetPosition("7203")
	assert.True(t, ok7203_updated)
	assert.Equal(t, 50, pos7203_updated.Quantity)

	// 以前のポジションが消えていることを確認
	pos9984_deleted, ok9984_deleted := state.GetPosition("9984")
	assert.False(t, ok9984_deleted)
	assert.Nil(t, pos9984_deleted)
}

func TestState_Orders(t *testing.T) {
	state := agent.NewState()

	// 初期状態の確認
	ord, ok := state.GetOrder("order-001")
	assert.False(t, ok)
	assert.Nil(t, ord)

	// 注文の更新
	orders := []*model.Order{
		{OrderID: "order-001", Symbol: "7203", OrderStatus: model.OrderStatusNew, Quantity: 100},
		{OrderID: "order-002", Symbol: "9984", OrderStatus: model.OrderStatusPartiallyFilled, Quantity: 200},
	}
	state.UpdateOrders(orders)

	// 更新後の取得確認
	ord001, ok001 := state.GetOrder("order-001")
	assert.True(t, ok001)
	assert.NotNil(t, ord001)
	assert.Equal(t, model.OrderStatusNew, ord001.OrderStatus)

	ord002, ok002 := state.GetOrder("order-002")
	assert.True(t, ok002)
	assert.NotNil(t, ord002)
	assert.Equal(t, model.OrderStatusPartiallyFilled, ord002.OrderStatus)
}

func TestState_Balance(t *testing.T) {
	state := agent.NewState()

	// 初期状態の確認 (0値の構造体)
	balance := state.GetBalance()
	assert.NotNil(t, balance)
	assert.Equal(t, 0.0, balance.Cash)
	assert.Equal(t, 0.0, balance.BuyingPower)

	// 残高の更新
	newBalance := &agent.Balance{Cash: 1000000, BuyingPower: 500000}
	state.UpdateBalance(newBalance)

	// 更新後の取得確認
	updatedBalance := state.GetBalance()
	assert.NotNil(t, updatedBalance)
	assert.Equal(t, 1000000.0, updatedBalance.Cash)
	assert.Equal(t, 500000.0, updatedBalance.BuyingPower)

	// GetBalanceがコピーを返すかの確認
	updatedBalance.Cash = 999 // 取得したものを変更してみる
	originalBalance := state.GetBalance()
	assert.Equal(t, 1000000.0, originalBalance.Cash) // 元の値は変わらないはず
}

func TestState_ThreadSafety(t *testing.T) {
	state := agent.NewState()
	var wg sync.WaitGroup
	numGoroutines := 100

	// 複数のゴルーチンから同時に書き込みと読み込みを行う
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			
			// 書き込み
			state.UpdateBalance(&agent.Balance{Cash: float64(i), BuyingPower: float64(i)})
			
			positions := []*model.Position{
				{Symbol: "7203", Quantity: i},
			}
			state.UpdatePositions(positions)

			// 読み込み
			state.GetBalance()
			state.GetPosition("7203")

		}(i)
	}

	wg.Wait()
	// このテストは、-raceフラグ付きで実行した際にデータ競合が検出されないことで成功とみなす
	// ここでは単純にパニックが起きずに終了することを確認
	t.Log("Thread safety test completed without panic.")
}
