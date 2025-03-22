package repository

import (
	"context"
	"myapp/domain/model"
)

type SignalRepository interface {
	Save(ctx context.Context, signal *model.Signal) error
	FindByID(ctx context.Context, id uint) (*model.Signal, error)
	FindBySymbol(ctx context.Context, symbol string) ([]*model.Signal, error) // 例: 銘柄コードでシグナルを検索

	// 他の必要なメソッドを定義
}
