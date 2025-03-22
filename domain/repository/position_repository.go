package repository

import (
	"context"
	"myapp/domain/model"
)

type PositionRepository interface {
	Save(ctx context.Context, position *model.Position) error
	FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) // 例: 銘柄コードでポジションを検索
	FindAll(ctx context.Context) ([]*model.Position, error)                   // 例: すべてのポジションを取得
	// 他の必要なメソッドを定義
}
