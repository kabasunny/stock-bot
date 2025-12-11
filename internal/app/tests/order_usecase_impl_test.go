package tests

import (
	"context"
	"errors"
	"stock-bot/domain/model"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client" // Added for client.NewOrderParams
	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for client.OrderClient
type OrderClientMock struct {
	mock.Mock
}

func (m *OrderClientMock) NewOrder(ctx context.Context, params client.NewOrderParams) (*response.ResNewOrder, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResNewOrder), args.Error(1)
}

func (m *OrderClientMock) CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResCorrectOrder), args.Error(1)
}

func (m *OrderClientMock) CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResCancelOrder), args.Error(1)
}

func (m *OrderClientMock) CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResCancelOrderAll), args.Error(1)
}

func (m *OrderClientMock) GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResOrderList), args.Error(1)
}

func (m *OrderClientMock) GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResOrderListDetail), args.Error(1)
}

// Mock for repository.OrderRepository
type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Save(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

// OrderUsecaseの実装をテスト
func TestExecuteOrder_Success(t *testing.T) {
	ctx := context.Background()

	// Mockのセットアップ
	orderClientMock := new(OrderClientMock)
	orderRepositoryMock := new(OrderRepositoryMock)

	// Tachibana OrderClient が成功レスポンスを返すように設定
	expectedResNewOrder := &response.ResNewOrder{
		ResultCode:  "0",
		OrderNumber: "test-order-id-123",
	}
	orderClientMock.On("NewOrder", ctx, mock.AnythingOfType("client.NewOrderParams")).Return(expectedResNewOrder, nil).Once() // Changed to client.NewOrderParams

	// OrderRepository がエラーなく保存するように設定
	orderRepositoryMock.On("Save", ctx, mock.AnythingOfType("*model.Order")).Return(nil).Once()

	// Usecaseの初期化
	uc := app.NewOrderUseCaseImpl(orderClientMock, orderRepositoryMock) // Removed secondPassword

	// 実行
	orderParams := app.OrderParams{
		Symbol:    "7203",
		TradeType: model.TradeTypeBuy,
		OrderType: model.OrderTypeMarket,
		Quantity:  100,
		Price:     0,
		IsMargin:  false,
	}
	result, err := uc.ExecuteOrder(ctx, orderParams)

	// アサート
	if assert.NoError(t, err) {
		assert.NotNil(t, result)
		assert.Equal(t, expectedResNewOrder.OrderNumber, result.OrderID)
		assert.Equal(t, orderParams.Symbol, result.Symbol)
		assert.Equal(t, orderParams.TradeType, result.TradeType)
		assert.Equal(t, orderParams.OrderType, result.OrderType)
		assert.Equal(t, int(orderParams.Quantity), result.Quantity) // uint64からintへの変換
		assert.Equal(t, orderParams.Price, result.Price)
		assert.Equal(t, orderParams.IsMargin, result.IsMargin)
		assert.Equal(t, model.OrderStatusNew, result.OrderStatus) // 新規注文が成功したら"NEW"ステータス
	}
	orderClientMock.AssertExpectations(t)
	orderRepositoryMock.AssertExpectations(t)
}

func TestExecuteOrder_ClientError(t *testing.T) {
	ctx := context.Background()

	// Mockのセットアップ
	orderClientMock := new(OrderClientMock)
	orderRepositoryMock := new(OrderRepositoryMock)

	// Tachibana OrderClient がエラーを返すように設定
	expectedErr := errors.New("failed to call Tachibana API")
	orderClientMock.On("NewOrder", ctx, mock.AnythingOfType("client.NewOrderParams")).Return(nil, expectedErr).Once() // Changed to client.NewOrderParams

	// Usecaseの初期化
	uc := app.NewOrderUseCaseImpl(orderClientMock, orderRepositoryMock) // Removed secondPassword

	// 実行
	orderParams := app.OrderParams{
		Symbol:    "7203",
		TradeType: model.TradeTypeBuy,
		OrderType: model.OrderTypeMarket,
		Quantity:  100,
		Price:     0,
		IsMargin:  false,
	}
	result, err := uc.ExecuteOrder(ctx, orderParams)

	// アサート
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	assert.Nil(t, result)
	orderClientMock.AssertExpectations(t)
	orderRepositoryMock.AssertNotCalled(t, "Save", ctx, mock.Anything) // Clientエラーの場合はRepositoryは呼ばれない
}

func TestExecuteOrder_RepositoryError(t *testing.T) {
	ctx := context.Background()

	// Mockのセットアップ
	orderClientMock := new(OrderClientMock)
	orderRepositoryMock := new(OrderRepositoryMock)

	// Tachibana OrderClient が成功レスポンスを返すように設定
	expectedResNewOrder := &response.ResNewOrder{
		ResultCode:  "0",
		OrderNumber: "test-order-id-456",
	}
	orderClientMock.On("NewOrder", ctx, mock.AnythingOfType("client.NewOrderParams")).Return(expectedResNewOrder, nil).Once() // Changed to client.NewOrderParams

	// OrderRepository がエラーを返すように設定
	expectedErr := errors.New("failed to save order to DB")
	orderRepositoryMock.On("Save", ctx, mock.AnythingOfType("*model.Order")).Return(expectedErr).Once()

	// Usecaseの初期化
	uc := app.NewOrderUseCaseImpl(orderClientMock, orderRepositoryMock) // Removed secondPassword

	// 実行
	orderParams := app.OrderParams{
		Symbol:    "7203",
		TradeType: model.TradeTypeBuy,
		OrderType: model.OrderTypeMarket,
		Quantity:  100,
		Price:     0,
		IsMargin:  false,
	}
	result, err := uc.ExecuteOrder(ctx, orderParams)

	// アサート
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	assert.Nil(t, result) // Repositoryエラーの場合は結果は返さない
	orderClientMock.AssertExpectations(t)
	orderRepositoryMock.AssertExpectations(t)
}

// go test -v ./internal/app/tests/order_usecase_impl_test.go
