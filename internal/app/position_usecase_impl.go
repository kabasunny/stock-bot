package app

import (
	"context"

	"stock-bot/domain/model"
	"stock-bot/domain/repository"
)

type positionUseCaseImpl struct {
	positionRepo repository.PositionRepository
}

// NewPositionUseCase : PositionUseCaseのコンストラクタ
func NewPositionUseCase(
	positionRepo repository.PositionRepository,
) PositionUseCase {
	return &positionUseCaseImpl{
		positionRepo: positionRepo,
	}
}

// List : 建玉一覧を取得
func (p *positionUseCaseImpl) List(ctx context.Context) ([]*model.Position, error) {
	return p.positionRepo.FindAll(ctx)
}
