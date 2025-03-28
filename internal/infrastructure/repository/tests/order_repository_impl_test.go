// internal/infrastructure/repository/tests/order_repository_impl_test.go

package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepositoryImpl_Save(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewOrderRepository(db)

	t.Run("正常系: Order を保存できること", func(t *testing.T) {
		ctx := context.Background()
		order := &model.Order{
			OrderID:     "test-order-1",
			Symbol:      "1234",
			TradeType:   model.TradeTypeBuy,
			OrderType:   model.OrderTypeLimit,
			Quantity:    100,
			Price:       1000.0,
			OrderStatus: model.OrderStatusNew,
			IsMargin:    false,
		}

		err := repo.Save(ctx, order)
		assert.NoError(t, err)

		retrievedOrder, err := repo.FindByID(ctx, "test-order-1")
		assert.NoError(t, err)
		assert.NotNil(t, retrievedOrder)
		assert.Equal(t, "1234", retrievedOrder.Symbol)
	})
}

func TestOrderRepositoryImpl_FindByID(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewOrderRepository(db)

	t.Run("正常系: OrderID で Order を取得できること", func(t *testing.T) {
		ctx := context.Background()
		orderID := "test-order-2"
		order := &model.Order{
			OrderID:     orderID,
			Symbol:      "5678",
			TradeType:   model.TradeTypeSell,
			OrderType:   model.OrderTypeMarket,
			Quantity:    50,
			Price:       0, // 成行注文なので 0
			OrderStatus: model.OrderStatusFilled,
			IsMargin:    true,
		}
		err := repo.Save(ctx, order)
		assert.NoError(t, err)

		retrievedOrder, err := repo.FindByID(ctx, orderID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedOrder)
		assert.Equal(t, "5678", retrievedOrder.Symbol)
	})

	t.Run("正常系: 存在しない OrderID を指定した場合 nil が返ること", func(t *testing.T) {
		ctx := context.Background()
		retrievedOrder, err := repo.FindByID(ctx, "non-existent-order")
		assert.NoError(t, err)
		assert.Nil(t, retrievedOrder)
	})
}

func TestOrderRepositoryImpl_FindByStatus(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewOrderRepository(db)

	t.Run("正常系: 指定した OrderStatus の Order が取得できること", func(t *testing.T) {
		ctx := context.Background()
		status := model.OrderStatusNew
		order1 := &model.Order{
			OrderID:     "test-order-3",
			Symbol:      "9012",
			TradeType:   model.TradeTypeBuy,
			OrderType:   model.OrderTypeLimit,
			Quantity:    200,
			Price:       1200.0,
			OrderStatus: status,
			IsMargin:    false,
		}
		order2 := &model.Order{
			OrderID:     "test-order-4",
			Symbol:      "3456",
			TradeType:   model.TradeTypeSell,
			OrderType:   model.OrderTypeMarket,
			Quantity:    100,
			Price:       0,
			OrderStatus: model.OrderStatusFilled, // 別のステータス
			IsMargin:    true,
		}
		err = repo.Save(ctx, order1)
		assert.NoError(t, err)
		err = repo.Save(ctx, order2)
		assert.NoError(t, err)

		retrievedOrders, err := repo.FindByStatus(ctx, status)
		assert.NoError(t, err)
		assert.Len(t, retrievedOrders, 1)
		assert.Equal(t, "9012", retrievedOrders[0].Symbol)
	})
}

// go test -v ./internal/infrastructure/repository/tests/order_repository_impl_test.go
