// internal/infrastructure/repository/tests/signal_repository_impl_test.go

package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignalRepositoryImpl_Save(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewSignalRepository(db)

	t.Run("正常系: Signal を保存できること", func(t *testing.T) {
		ctx := context.Background()
		signal := &model.Signal{
			Symbol:      "1234",
			SignalType:  model.SignalTypeBuy,
			GeneratedAt: time.Now(),
			Rationale:   "テストシグナル",
			Price:       1000.0,
		}

		err := repo.Save(ctx, signal)
		assert.NoError(t, err)

		retrievedSignal, err := repo.FindByID(ctx, signal.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedSignal)
		assert.Equal(t, model.SignalTypeBuy, retrievedSignal.SignalType)
	})
}

func TestSignalRepositoryImpl_FindByID(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewSignalRepository(db)

	t.Run("正常系: ID で Signal を取得できること", func(t *testing.T) {
		ctx := context.Background()
		signal := &model.Signal{
			Symbol:      "5678",
			SignalType:  model.SignalTypeSell,
			GeneratedAt: time.Now(),
			Rationale:   "テストシグナル",
			Price:       1200.0,
		}
		err := repo.Save(ctx, signal)
		assert.NoError(t, err)

		retrievedSignal, err := repo.FindByID(ctx, signal.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedSignal)
		assert.Equal(t, model.SignalTypeSell, retrievedSignal.SignalType)
	})

	t.Run("正常系: 存在しない ID を指定した場合 nil が返ること", func(t *testing.T) {
		ctx := context.Background()
		retrievedSignal, err := repo.FindByID(ctx, 9999)
		assert.NoError(t, err)
		assert.Nil(t, retrievedSignal)
	})
}

func TestSignalRepositoryImpl_FindBySymbol(t *testing.T) {
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := repository.NewSignalRepository(db)

	t.Run("正常系: Symbol で Signal が取得できること", func(t *testing.T) {
		ctx := context.Background()
		symbol := "9012"
		signal1 := &model.Signal{
			Symbol:      symbol,
			SignalType:  model.SignalTypeBuy,
			GeneratedAt: time.Now(),
			Rationale:   "テストシグナル",
			Price:       1500.0,
		}
		signal2 := &model.Signal{
			Symbol:      "3456",
			SignalType:  model.SignalTypeSell,
			GeneratedAt: time.Now(),
			Rationale:   "テストシグナル",
			Price:       800.0,
		}
		err = repo.Save(ctx, signal1)
		assert.NoError(t, err)
		err = repo.Save(ctx, signal2)
		assert.NoError(t, err)

		retrievedSignals, err := repo.FindBySymbol(ctx, symbol)
		assert.NoError(t, err)
		assert.Len(t, retrievedSignals, 1)
		assert.Equal(t, model.SignalTypeBuy, retrievedSignals[0].SignalType)
	})
}

// go test -v ./internal/infrastructure/repository/tests/signal_repository_impl_test.go
