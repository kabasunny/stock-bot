// internal/infrastructure/repository/test_helper.go
package repository

import (
	"fmt"
	"log"
	"stock-bot/domain/model"
	"stock-bot/internal/config"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDatabase はテスト用のデータベース環境を構築し、GORM DB インスタンスとクリーンアップ関数を返します。
func SetupTestDatabase(t *testing.T) (*gorm.DB, func(), error) {
	t.Helper()

	// Docker pool の初期化
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	// PostgreSQL コンテナのオプション設定
	runOpts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15", // PostgreSQL のバージョンを指定
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpassword",
			"POSTGRES_DB=testdb",
		},
	}

	// Docker コンテナの起動
	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	// データベースへの接続設定
	var db *gorm.DB
	dsn := fmt.Sprintf("host=localhost port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", resource.GetPort("5432/tcp"))

	// 接続のリトライ設定
	pool.MaxWait = 120 * time.Second // 最大待ち時間

	// データベースへの接続をリトライ
	if err := pool.Retry(func() error {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Database connection error:", err)
			return err
		}

		// 接続確認
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	// テストに必要なテーブルのマイグレーションを実行
	err = db.AutoMigrate(&model.Order{}, &model.Position{}, &model.Signal{}, &model.StockMaster{}, &model.StockMarketMaster{}, &model.TickRule{}, &model.TickLevel{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// クリーンアップ関数の定義
	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	return db, cleanup, nil
}

// CreateTestConfig はテスト用のConfigインスタンスを作成します。
func CreateTestConfig() *config.Config {
	// 必要に応じて、Config構造体のフィールドを設定します。
	// ここでは、最小限の設定のみ行っています。
	return &config.Config{
		DBHost:     "localhost",    // Dockerコンテナのホスト
		DBPort:     5432,           // Dockerコンテナのポート
		DBUser:     "testuser",     // ユーザー名
		DBPassword: "testpassword", // パスワード
		DBName:     "testdb",       // データベース名
	}
}

// setupLogger はテスト用のロガーをセットアップします。
func SetupLogger(t *testing.T) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	return logger
}
