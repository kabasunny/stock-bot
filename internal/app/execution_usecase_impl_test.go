package app

import (
	"context"
	"errors"
	"stock-bot/domain/model"
	"stock-bot/internal/app/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecutionUseCase_Execute(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockPositionRepo := new(mocks.PositionRepository)
	useCase := NewExecutionUseCaseImpl(mockOrderRepo, mockPositionRepo)

	ctx := context.Background()

	t.Run("successful execution processing for buy order", func(t *testing.T) {
		execution := &model.Execution{
			ExecutionID: "EXEC001",
			OrderID:     "ORDER001",
			Symbol:      "7203",
			TradeType:   model.TradeTypeBuy,
			Quantity:    100,
			Price:       1500.0,
			ExecutedAt:  time.Now(),
		}

		mockOrderRepo.On("UpdateOrderStatusByExecution", ctx, execution).Return(nil).Once()
		mockPositionRepo.On("UpsertPositionByExecution", ctx, execution).Return(nil).Once()

		err := useCase.Execute(ctx, execution)
		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockPositionRepo.AssertExpectations(t)
	})

	t.Run("successful execution processing for sell order", func(t *testing.T) {
		execution := &model.Execution{
			ExecutionID: "EXEC002",
			OrderID:     "ORDER002",
			Symbol:      "6758",
			TradeType:   model.TradeTypeSell,
			Quantity:    200,
			Price:       20000.0,
			ExecutedAt:  time.Now(),
		}

		mockOrderRepo.On("UpdateOrderStatusByExecution", ctx, execution).Return(nil).Once()
		mockPositionRepo.On("UpsertPositionByExecution", ctx, execution).Return(nil).Once()

		err := useCase.Execute(ctx, execution)
		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockPositionRepo.AssertExpectations(t)
	})

	t.Run("order repository update fails", func(t *testing.T) {
		mockOrderRepo := new(mocks.OrderRepository)
		mockPositionRepo := new(mocks.PositionRepository)
		useCase := NewExecutionUseCaseImpl(mockOrderRepo, mockPositionRepo)
		execution := &model.Execution{
			ExecutionID: "EXEC003",
			OrderID:     "ORDER003",
			Symbol:      "7203",
			TradeType:   model.TradeTypeBuy,
			Quantity:    100,
			Price:       1500.0,
			ExecutedAt:  time.Now(),
		}
		expectedErr := errors.New("failed to update order status")

		mockOrderRepo.On("UpdateOrderStatusByExecution", ctx, execution).Return(expectedErr).Once()
		err := useCase.Execute(ctx, execution)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockOrderRepo.AssertExpectations(t)
		mockPositionRepo.AssertNotCalled(t, "UpsertPositionByExecution", mock.Anything, mock.Anything)
	})
	t.Run("position repository upsert fails", func(t *testing.T) {
		execution := &model.Execution{
			ExecutionID: "EXEC004",
			OrderID:     "ORDER004",
			Symbol:      "7203",
			TradeType:   model.TradeTypeBuy,
			Quantity:    100,
			Price:       1500.0,
			ExecutedAt:  time.Now(),
		}
		expectedErr := errors.New("failed to upsert position")

		mockOrderRepo.On("UpdateOrderStatusByExecution", ctx, execution).Return(nil).Once()
		mockPositionRepo.On("UpsertPositionByExecution", ctx, execution).Return(expectedErr).Once()

		err := useCase.Execute(ctx, execution)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockOrderRepo.AssertExpectations(t)
		mockPositionRepo.AssertExpectations(t)
	})
}

// 既存の OrderRepository と PositionRepository のモックインターフェースを生成
// mockery --name OrderRepository --inpkg --output internal/app/mocks
// mockery --name PositionRepository --inpkg --output internal/app/mocks
