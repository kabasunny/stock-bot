package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockOrderRepository はテスト用の OrderRepository のモック実装です。
type MockOrderRepository struct {
	// 必要に応じてメソッドをモックするためにフィールドを追加
}

// Save は OrderRepository インターフェースの Save メソッドのダミー実装です。
func (m *MockOrderRepository) Save(ctx context.Context, order *model.Order) error {
	return nil
}

// FindByID は OrderRepository インターフェースの FindByID メソッドのダミー実装です。
func (m *MockOrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	// テストによっては、特定の orderID に対してダミーの Order を返すように設定することもできます
	// 今回の position_repository のテストでは、FindByID が呼ばれることはないので、nil を返します
	return nil, nil
}

// FindByStatus は OrderRepository インターフェースの FindByStatus メソッドのダミー実装です。
func (m *MockOrderRepository) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return nil, nil
}

// UpdateOrderStatusByExecution は OrderRepository インターフェースの UpdateOrderStatusByExecution メソッドのダミー実装です。
func (m *MockOrderRepository) UpdateOrderStatusByExecution(ctx context.Context, execution *model.Execution) error {
	return nil
}

func TestPositionRepositoryImpl_Save(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	mockOrderRepo := &MockOrderRepository{}
	repo := repository.NewPositionRepository(db, mockOrderRepo)

	t.Run("正常系: Position を保存できること", func(t *testing.T) {
		ctx := context.Background()
		position := &model.Position{
			Symbol:              "1234",
			PositionType:        model.PositionTypeLong,
			PositionAccountType: model.PositionAccountTypeCash, // AccountType を追加
			AveragePrice:        1000.0,
			Quantity:            100,
		}

		err := repo.Save(ctx, position)
		assert.NoError(t, err)

		retrievedPosition, err := repo.FindBySymbol(ctx, "1234")
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPosition)
		assert.Equal(t, model.PositionTypeLong, retrievedPosition.PositionType)
		assert.Equal(t, model.PositionAccountTypeCash, retrievedPosition.PositionAccountType) // AccountType の比較を追加
	})
}

func TestPositionRepositoryImpl_FindBySymbol(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	mockOrderRepo := &MockOrderRepository{}
	repo := repository.NewPositionRepository(db, mockOrderRepo)

	t.Run("正常系: Symbol で Position を取得できること", func(t *testing.T) {
		ctx := context.Background()
		symbol := "5678"
				position := &model.Position{
					Symbol:              symbol,
					PositionType:        model.PositionTypeShort,
					PositionAccountType: model.PositionAccountTypeMarginNew, // AccountType を更新
					AveragePrice:        1200.0,
					Quantity:            50,
				}
		err := repo.Save(ctx, position)
		assert.NoError(t, err)

		retrievedPosition, err := repo.FindBySymbol(ctx, symbol)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPosition)
		assert.Equal(t, model.PositionTypeShort, retrievedPosition.PositionType)
		assert.Equal(t, model.PositionAccountTypeMarginNew, retrievedPosition.PositionAccountType) // AccountType の比較を更新
	})

	t.Run("正常系: 存在しない Symbol を指定した場合 nil が返ること", func(t *testing.T) {
		ctx := context.Background()
		retrievedPosition, err := repo.FindBySymbol(ctx, "non-existent-symbol")
		assert.NoError(t, err)
		assert.Nil(t, retrievedPosition)
	})
}

func TestPositionRepositoryImpl_FindAll(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	mockOrderRepo := &MockOrderRepository{}
	repo := repository.NewPositionRepository(db, mockOrderRepo)

	t.Run("正常系: すべての Position が取得できること", func(t *testing.T) {
		ctx := context.Background()
		position1 := &model.Position{
			Symbol:              "9012",
			PositionType:        model.PositionTypeLong,
			PositionAccountType: model.PositionAccountTypeCash, // AccountType を追加
			AveragePrice:        1500.0,
			Quantity:            200,
		}
				position2 := &model.Position{
					Symbol:              "3456",
					PositionType:        model.PositionTypeShort,
					PositionAccountType: model.PositionAccountTypeMarginNew, // AccountType を更新
					AveragePrice:        800.0,
					Quantity:            100,
				}
		err = repo.Save(ctx, position1)
		assert.NoError(t, err)
		err = repo.Save(ctx, position2)
		assert.NoError(t, err)

		retrievedPositions, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, retrievedPositions, 2)
		// AccountType も確認
		assert.Contains(t, retrievedPositions, position1) // GORMの比較はポインタなので注意
		assert.Contains(t, retrievedPositions, position2)
	})
}
