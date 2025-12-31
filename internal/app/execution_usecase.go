package app

import (
	"context"
	"stock-bot/domain/model"
)

// ExecutionUseCase は約定情報を受け取り、注文とポジションを更新するユースケースのインターフェースです。
type ExecutionUseCase interface {
	Execute(ctx context.Context, execution *model.Execution) error
}
