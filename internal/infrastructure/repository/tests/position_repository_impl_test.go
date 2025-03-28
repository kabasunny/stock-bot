// internal/infrastructure/repository/tests/position_repository_impl_test.go

package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionRepositoryImpl_Save(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewPositionRepository(db)

	t.Run("正常系: Position を保存できること", func(t *testing.T) {
		ctx := context.Background()
		position := &model.Position{
			Symbol:       "1234",
			PositionType: model.PositionTypeLong,
			AveragePrice: 1000.0,
			Quantity:     100,
		}

		err := repo.Save(ctx, position)
		assert.NoError(t, err)

		retrievedPosition, err := repo.FindBySymbol(ctx, "1234")
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPosition)
		assert.Equal(t, model.PositionTypeLong, retrievedPosition.PositionType)
	})
}

func TestPositionRepositoryImpl_FindBySymbol(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewPositionRepository(db)

	t.Run("正常系: Symbol で Position を取得できること", func(t *testing.T) {
		ctx := context.Background()
		symbol := "5678"
		position := &model.Position{
			Symbol:       symbol,
			PositionType: model.PositionTypeShort,
			AveragePrice: 1200.0,
			Quantity:     50,
		}
		err := repo.Save(ctx, position)
		assert.NoError(t, err)

		retrievedPosition, err := repo.FindBySymbol(ctx, symbol)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPosition)
		assert.Equal(t, model.PositionTypeShort, retrievedPosition.PositionType)
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

	repo := repository.NewPositionRepository(db)

	t.Run("正常系: すべての Position が取得できること", func(t *testing.T) {
		ctx := context.Background()
		position1 := &model.Position{
			Symbol:       "9012",
			PositionType: model.PositionTypeLong,
			AveragePrice: 1500.0,
			Quantity:     200,
		}
		position2 := &model.Position{
			Symbol:       "3456",
			PositionType: model.PositionTypeShort,
			AveragePrice: 800.0,
			Quantity:     100,
		}
		err = repo.Save(ctx, position1)
		assert.NoError(t, err)
		err = repo.Save(ctx, position2)
		assert.NoError(t, err)

		retrievedPositions, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, retrievedPositions, 2)
	})
}

// go test -v ./internal/infrastructure/repository/tests/position_repository_impl_test.go
