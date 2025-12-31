package repository

import (
	"context"
	"stock-bot/domain/model"
)

type PositionRepository interface {
	Save(ctx context.Context, position *model.Position) error
	FindBySymbol(ctx context.Context, symbol string) (*model.Position, error) // 例: 銘柄コードでポジションを検索
	FindAll(ctx context.Context) ([]*model.Position, error)                   // 例: すべてのポジションを取得
	UpdateHighestPrice(ctx context.Context, symbol string, price float64) error
	// UpsertPositionByExecution は約定情報に基づいてポジションを新規作成または更新します。
	UpsertPositionByExecution(ctx context.Context, execution *model.Execution) error
	// DeletePosition はポジションを削除します（例：全株売却時）。
	DeletePosition(ctx context.Context, symbol string, accountType model.PositionAccountType) error
	// 他の必要なメソッドを定義
}
