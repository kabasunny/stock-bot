package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"stock-bot/internal/config"
	"stock-bot/internal/infrastructure/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("=== 保存された注文確認 ===")

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

	// データベース接続
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// OrderRepositoryの作成
	orderRepo := repository.NewOrderRepository(db)

	// 最近の注文を取得
	ctx := context.Background()
	orders, err := orderRepo.FindOrderHistory(ctx, nil, nil, 10)
	if err != nil {
		fmt.Printf("❌ 注文履歴取得エラー: %v\n", err)
		return
	}

	fmt.Printf("最近の注文 (%d件):\n", len(orders))
	fmt.Println("----------------------------------------")

	for i, order := range orders {
		fmt.Printf("%d. OrderID: %s\n", i+1, order.OrderID)
		fmt.Printf("   Symbol: %s\n", order.Symbol)
		fmt.Printf("   TradeType: %s\n", order.TradeType)
		fmt.Printf("   OrderType: %s\n", order.OrderType)
		fmt.Printf("   Quantity: %d\n", order.Quantity)
		fmt.Printf("   Price: %.2f\n", order.Price)
		fmt.Printf("   TriggerPrice: %.2f\n", order.TriggerPrice)
		fmt.Printf("   PositionAccountType: %s\n", order.PositionAccountType)
		fmt.Printf("   OrderStatus: %s\n", order.OrderStatus)
		fmt.Printf("   CreatedAt: %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("   ---")
	}

	fmt.Println("\n=== 確認完了 ===")
}
