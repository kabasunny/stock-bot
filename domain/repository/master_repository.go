// domain/repository/master_repository.go
package repository

import (
	"context"
	"stock-bot/domain/model"
)

type MasterRepository interface {
	SaveStockMaster(ctx context.Context, master *model.Master) error
	SaveStockMasters(ctx context.Context, masters []*model.Master) error // 一括保存
	//...
	FindStockMasterByIssueCode(ctx context.Context, issueCode string) (*model.Master, error)
	// 他のエンティティ（市場など）についても、必要に応じてメソッドを定義
}
