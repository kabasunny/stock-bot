package app

import (
	"context"
	"stock-bot/domain/model"
)

// PositionUseCase : 建玉に関するユースケース
type PositionUseCase interface {
	// List : 建玉一覧を取得
	List(ctx context.Context) ([]*model.Position, error)
}
