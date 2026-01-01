package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"stock-bot/internal/config"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("=== データベーススキーマ確認 ===")

	// 設定ファイルの読み込み
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller information")
	}
	envPath := filepath.Join(filepath.Dir(filename), ".env")

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// データベース接続文字列を構築
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// ordersテーブルの構造を確認
	fmt.Println("\n1. ordersテーブルのカラム一覧:")
	query := `
		SELECT column_name, data_type, is_nullable, column_default 
		FROM information_schema.columns 
		WHERE table_name = 'orders' 
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to query table schema: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-25s %-20s %-10s %s\n", "Column Name", "Data Type", "Nullable", "Default")
	fmt.Println(strings.Repeat("-", 80))

	for rows.Next() {
		var columnName, dataType, isNullable string
		var columnDefault sql.NullString

		err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		defaultValue := "NULL"
		if columnDefault.Valid {
			defaultValue = columnDefault.String
		}

		fmt.Printf("%-25s %-20s %-10s %s\n", columnName, dataType, isNullable, defaultValue)
	}

	// position_account_typeカラムの存在確認
	fmt.Println("\n2. position_account_typeカラムの存在確認:")
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'orders' AND column_name = 'position_account_type'
		)
	`

	err = db.QueryRow(checkQuery).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check column existence: %v", err)
	}

	if exists {
		fmt.Println("✅ position_account_typeカラムは存在します")
	} else {
		fmt.Println("❌ position_account_typeカラムが存在しません")
	}

	// マイグレーション履歴の確認
	fmt.Println("\n3. マイグレーション履歴:")
	migrationQuery := `
		SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 5
	`

	migrationRows, err := db.Query(migrationQuery)
	if err != nil {
		log.Printf("Failed to query migration history: %v", err)
		return
	}
	defer migrationRows.Close()

	fmt.Printf("%-20s %s\n", "Version", "Dirty")
	fmt.Println(strings.Repeat("-", 30))

	for migrationRows.Next() {
		var version string
		var dirty bool

		err := migrationRows.Scan(&version, &dirty)
		if err != nil {
			log.Printf("Error scanning migration row: %v", err)
			continue
		}

		fmt.Printf("%-20s %t\n", version, dirty)
	}

	fmt.Println("\n=== 確認完了 ===")
}
