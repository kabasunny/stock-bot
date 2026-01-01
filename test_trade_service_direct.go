package main

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/tradeservice"
)

// MockOrderRepository はテスト用のモックリポジトリ
type MockOrderRepository struct{}

func (m *MockOrderRepository) Save(ctx context.Context, order *model.Order) error {
	fmt.Printf("Mock: 注文を保存しました - ID: %s, Symbol: %s\n", order.OrderID, order.Symbol)
	return nil
}

func (m *MockOrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	return nil, fmt.Errorf("order not found")
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

func (m *MockOrderRepository) FindBySymbol(ctx context.Context, symbol string) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

// MockMasterRepository はテスト用のモックリポジトリ
type MockMasterRepository struct{}

func (m *MockMasterRepository) SaveSymbol(ctx context.Context, symbol *model.Symbol) error {
	return nil
}

func (m *MockMasterRepository) FindSymbolByCode(ctx context.Context, code string) (*model.Symbol, error) {
	return &model.Symbol{Code: code, Name: "テスト銘柄"}, nil
}

func main() {
	fmt.Println("=== TradeService 直接テスト ===")

	// 1. 立花クライアントを作成
	tachibanaClient := client.CreateTestClient(nil)

	// 2. ログイン
	ctx := context.Background()
	loginReq := client.ReqLogin{
		UserId:   tachibanaClient.GetUserIDForTest(),
		Password: tachibanaClient.GetPasswordForTest(),
	}

	session, err := tachibanaClient.LoginWithPost(ctx, loginReq)
	if err != nil {
		fmt.Printf("ログインエラー: %v\n", err)
		return
	}
	fmt.Printf("ログイン成功: %s\n", session.RequestURL)

	// 3. TradeServiceを作成
	mockOrderRepo := &MockOrderRepository{}
	mockMasterRepo := &MockMasterRepository{}
	logger := slog.Default()

	tradeService := tradeservice.NewGoaTradeService(
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		mockOrderRepo,
		mockMasterRepo,
		session,
		logger,
	)

	// 4. 注文リクエストを作成（立花クライアントで成功したパラメータ）
	orderReq := &service.PlaceOrderRequest{
		Symbol:              "6658", // シスメックス
		TradeType:           model.TradeTypeBuy,
		OrderType:           model.OrderTypeMarket,
		Quantity:            100,
		Price:               0,   // 成行なので0
		TriggerPrice:        nil, // 成行なのでnil
		PositionAccountType: model.PositionAccountTypeCash,
	}

	fmt.Printf("注文リクエスト: %+v\n", orderReq)

	// 5. 注文を実行
	order, err := tradeService.PlaceOrder(ctx, orderReq)
	if err != nil {
		fmt.Printf("❌ 注文エラー: %v\n", err)
		return
	}

	fmt.Printf("✅ 注文成功!\n")
	fmt.Printf("注文ID: %s\n", order.OrderID)
	fmt.Printf("銘柄: %s\n", order.Symbol)
	fmt.Printf("売買: %s\n", order.TradeType)
	fmt.Printf("注文種別: %s\n", order.OrderType)
	fmt.Printf("数量: %d\n", order.Quantity)
	fmt.Printf("ステータス: %s\n", order.OrderStatus)
	fmt.Printf("口座区分: %s\n", order.PositionAccountType)

	fmt.Println("\n=== テスト完了 ===")
}
