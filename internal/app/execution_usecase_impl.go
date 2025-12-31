package app

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
)

// ExecutionUseCaseImpl は ExecutionUseCase の実装です。
type ExecutionUseCaseImpl struct {
	orderRepo    repository.OrderRepository
	positionRepo repository.PositionRepository
}

// NewExecutionUseCaseImpl は ExecutionUseCaseImpl の新しいインスタンスを作成します。
func NewExecutionUseCaseImpl(
	orderRepo repository.OrderRepository,
	positionRepo repository.PositionRepository,
) *ExecutionUseCaseImpl {
	return &ExecutionUseCaseImpl{
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
	}
}

// Execute は約定情報に基づいて注文とポジションを更新します。
func (uc *ExecutionUseCaseImpl) Execute(ctx context.Context, execution *model.Execution) error {
	fmt.Println("start execution")
	// 注文の状態を更新
	err := uc.orderRepo.UpdateOrderStatusByExecution(ctx, execution)
	if err != nil {
		fmt.Println("error from orderRepo.UpdateOrderStatusByExecution")
		return fmt.Errorf("failed to update order status by execution: %w", err)
	}

	fmt.Println("after orderRepo.UpdateOrderStatusByExecution")
	// ポジションを更新
	err = uc.positionRepo.UpsertPositionByExecution(ctx, execution)
	if err != nil {
		fmt.Println("error from positionRepo.UpsertPositionByExecution")
		return fmt.Errorf("failed to upsert position by execution: %w", err)
	}

	fmt.Println("end execution")
	return nil
}
