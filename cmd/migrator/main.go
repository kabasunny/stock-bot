// cmd/migrator/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. .env ファイルをロード
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// 2. 環境変数からDB接続情報を取得し、DSNを組み立てる
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	log.Println("Connecting to database...")

	// 3. データベースに接続
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection successful.")

	// 4. マイグレーションドライバを初期化
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	// 5. マイグレーションインスタンスを生成
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// 6. マイグレーションを適用 (Up)
	log.Println("Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// バージョンを確認
	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("Failed to get migration version: %v", err)
	}

	log.Printf("Migrations applied successfully. Current version: %d, Dirty: %v\n", version, dirty)
}
