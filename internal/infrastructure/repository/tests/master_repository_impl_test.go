// internal/infrastructure/repository/tests/master_repository_impl_test.go
package tests

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterRepositoryImpl_Save(t *testing.T) {
	// テスト用データベースのセットアップ
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	// リポジトリの実装を作成
	repo := repository.NewMasterRepository(db)

	t.Run("正常系: StockMaster を保存できること", func(t *testing.T) {
		ctx := context.Background()
		stockMaster := &model.StockMaster{
			IssueCode:   "1234",
			IssueName:   "テスト銘柄",
			TradingUnit: 100,
			MarketCode:  "01",
		}

		// 保存処理を実行
		err := repo.Save(ctx, stockMaster)
		assert.NoError(t, err)

		// 保存されたデータを検証 (例: データベースから取得)
		var retrievedStockMaster model.StockMaster
		result := db.WithContext(ctx).Where("issue_code = ?", "1234").First(&retrievedStockMaster)
		assert.NoError(t, result.Error)
		assert.Equal(t, "テスト銘柄", retrievedStockMaster.IssueName)
	})

	// 他のテストケース (異常系など) を追加
}

func TestMasterRepositoryImpl_FindByIssueCode(t *testing.T) {
	// テスト用データベースのセットアップ
	db, cleanup, err := repository.SetupTestDatabase(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	// リポジトリの実装を作成
	repo := repository.NewMasterRepository(db)

	t.Run("正常系: IssueCodeでStockMasterを取得できること", func(t *testing.T) {
		ctx := context.Background()
		issueCode := "7777"

		// 事前にテストデータを登録
		stockMaster := &model.StockMaster{
			IssueCode:   issueCode,
			IssueName:   "テスト銘柄",
			TradingUnit: 100,
			MarketCode:  "01",
		}
		err := repo.Save(ctx, stockMaster)
		assert.NoError(t, err)

		// FindByIssueCodeを実行
		retrievedEntity, err := repo.FindByIssueCode(ctx, issueCode, "StockMaster")
		assert.NoError(t, err)

		// 取得できたエンティティがStockMaster型であることを確認
		retrievedStockMaster, ok := retrievedEntity.(*model.StockMaster)
		assert.True(t, ok)
		assert.NotNil(t, retrievedStockMaster)

		// 取得したデータの検証
		assert.Equal(t, issueCode, retrievedStockMaster.IssueCode)
		assert.Equal(t, "テスト銘柄", retrievedStockMaster.IssueName)

	})

	t.Run("正常系: 存在しないIssueCodeを指定した場合nilが返ること", func(t *testing.T) {
		ctx := context.Background()
		issueCode := "9999" // 存在しないIssueCode

		// FindByIssueCodeを実行
		retrievedEntity, err := repo.FindByIssueCode(ctx, issueCode, "StockMaster")
		assert.NoError(t, err)
		assert.Nil(t, retrievedEntity) // nilが返ることを確認
	})

	t.Run("異常系: 存在しないエンティティタイプを指定した場合", func(t *testing.T) {
		ctx := context.Background()
		issueCode := "1234"

		// FindByIssueCodeを実行
		_, err := repo.FindByIssueCode(ctx, issueCode, "InvalidType")
		assert.Error(t, err) // エラーが返ることを確認
	})
}

// go test -v ./internal/infrastructure/repository/tests/master_repository_impl_test.go
