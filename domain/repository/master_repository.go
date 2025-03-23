package repository

import (
	"context"
)

type MasterRepository interface {
	Save(ctx context.Context, entity interface{}) error
	SaveAll(ctx context.Context, entities []interface{}) error
	FindByIssueCode(ctx context.Context, issueCode string, entityType string) (interface{}, error)
	// Find(ctx context.Context, conditions map[string]interface{}, entityType string) ([]interface{}, error) // より汎用的な検索
	// Delete(ctx context.Context, entity interface{}) error // 削除が必要な場合
}
