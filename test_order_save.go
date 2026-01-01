package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"stock-bot/domain/model"
	"stock-bot/internal/config"
	"stock-bot/internal/infrastructure/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("=== 注文保存テスト ===")

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

	// テスト用の注文データを作成
	testOrder := &model.Order{
		OrderID:             "TEST_ORDER_001",
		Symbol:              "3632",
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeStopLimit,
		Quantity:            100,
		Price:               455.0,
		TriggerPrice:        460.0,
		TimeInForce:         model.TimeInForceDay,
		OrderStatus:         model.OrderStatusNew,
		IsMargin:            false,
		PositionAccountType: model.PositionAccountTypeCash,
	}

	fmt.Printf("保存する注文データ:\n")
	fmt.Printf("  OrderID: %s\n", testOrder.OrderID)
	fmt.Printf("  Symbol: %s\n", testOrder.Symbol)
	fmt.Printf("  TradeType: %s\n", testOrder.TradeType)
	fmt.Printf("  OrderType: %s\n", testOrder.OrderType)
	fmt.Printf("  Quantity: %d\n", testOrder.Quantity)
	fmt.Printf("  Price: %.2f\n", testOrder.Price)
	fmt.Printf("  TriggerPrice: %.2f\n", testOrder.TriggerPrice)
	fmt.Printf("  PositionAccountType: %s\n", testOrder.PositionAccountType)

	// 注文を保存
	ctx := context.Background()
	err = orderRepo.Save(ctx, testOrder)
	if err != nil {
		fmt.Printf("❌ 注文保存エラー: %v\n", err)
		return
	}

	fmt.Printf("✅ 注文保存成功!\n")

	// 保存された注文を取得して確認
	savedOrder, err := orderRepo.FindByID(ctx, testOrder.OrderID)
	if err != nil {
		fmt.Printf("❌ 注文取得エラー: %v\n", err)
		return
	}

	if savedOrder == nil {
		fmt.Printf("❌ 注文が見つかりません\n")
		return
	}

	fmt.Printf("\n保存された注文データ:\n")
	fmt.Printf("  OrderID: %s\n", savedOrder.OrderID)
	fmt.Printf("  Symbol: %s\n", savedOrder.Symbol)
	fmt.Printf("  TradeType: %s\n", savedOrder.TradeType)
	fmt.Printf("  OrderType: %s\n", savedOrder.OrderType)
	fmt.Printf("  Quantity: %d\n", savedOrder.Quantity)
	fmt.Printf("  Price: %.2f\n", savedOrder.Price)
	fmt.Printf("  TriggerPrice: %.2f\n", savedOrder.TriggerPrice)
	fmt.Printf("  PositionAccountType: %s\n", savedOrder.PositionAccountType)

	fmt.Println("\n=== テスト完了 ===")
}
