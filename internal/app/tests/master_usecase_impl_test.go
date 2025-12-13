package tests

import (
	"context"
	"errors"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MasterDataClientMock implements the client.MasterDataClient interface for testing.
type MasterDataClientMock struct {
	mock.Mock
}

// Ensure all methods of client.MasterDataClient are implemented.
func (m *MasterDataClientMock) DownloadMasterData(ctx context.Context, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetMasterDataQuery(ctx context.Context, req request.ReqGetMasterData) (*response.ResGetMasterData, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGetMasterData), args.Error(1)
}
func (m *MasterDataClientMock) GetNewsHeader(ctx context.Context, req request.ReqGetNewsHead) (*response.ResGetNewsHeader, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetNewsBody(ctx context.Context, req request.ReqGetNewsBody) (*response.ResGetNewsBody, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetIssueDetail(ctx context.Context, req request.ReqGetIssueDetail) (*response.ResGetIssueDetail, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGetIssueDetail), args.Error(1)
}
func (m *MasterDataClientMock) GetMarginInfo(ctx context.Context, req request.ReqGetMarginInfo) (*response.ResGetMarginInfo, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetCreditInfo(ctx context.Context, req request.ReqGetCreditInfo) (*response.ResGetCreditInfo, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetMarginPremiumInfo(ctx context.Context, req request.ReqGetMarginPremiumInfo) (*response.ResGetMarginPremiumInfo, error) {
	args := m.Called(ctx, req)
	return nil, args.Error(1)
}

func TestGetStock_Success(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	symbol := "7203"

	// --- Mock Data ---
	apiResponse := &response.ResGetMasterData{
		StockMaster: []response.ResStockMaster{
			{
				IssueCode:       symbol,
				IssueName:       "トヨタ自動車",
				IssueNameKana:   "トヨタジドウシャ",
				PreferredMarket: "東証プライム",
				IndustryCode:    "3600",
				IndustryName:    "輸送用機器",
			},
			{
				IssueCode:       "6758",
				IssueName:       "ソニーグループ",
				IssueNameKana:   "ソニーグループ",
				PreferredMarket: "東証プライム",
				IndustryCode:    "3650",
				IndustryName:    "電気機器",
			},
		},
	}
	expectedReq := request.ReqGetMasterData{
		TargetCLMID: "CLMIssueMstKabu",
	}

	masterClientMock.On("GetMasterDataQuery", ctx, expectedReq).Return(apiResponse, nil).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock)

	// --- Execute ---
	result, err := uc.GetStock(ctx, symbol)

	// --- Assert ---
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, "トヨタ自動車", result.Name)
	assert.Equal(t, "トヨタジドウシャ", result.NameKana)
	assert.Equal(t, "東証プライム", result.Market)
	assert.Equal(t, "3600", result.IndustryCode)
	assert.Equal(t, "輸送用機器", result.IndustryName)

	masterClientMock.AssertExpectations(t)
}

func TestGetStock_NotFound(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	symbol := "9999" // Symbol not in mock data

	// --- Mock Data ---
	apiResponse := &response.ResGetMasterData{
		StockMaster: []response.ResStockMaster{
			{
				IssueCode:       "7203",
				IssueName:       "トヨタ自動車",
				IssueNameKana:   "トヨタジドウシャ",
				PreferredMarket: "東証プライム",
				IndustryCode:    "3600",
				IndustryName:    "輸送用機器",
			},
		},
	}
	expectedReq := request.ReqGetMasterData{
		TargetCLMID: "CLMIssueMstKabu",
	}

	masterClientMock.On("GetMasterDataQuery", ctx, expectedReq).Return(apiResponse, nil).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock)

	// --- Execute ---
	result, err := uc.GetStock(ctx, symbol)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, app.ErrNotFound))

	masterClientMock.AssertExpectations(t)
}

func TestGetStock_ClientError(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	symbol := "7203"
	expectedErr := errors.New("API communication failed")
	expectedReq := request.ReqGetMasterData{
		TargetCLMID: "CLMIssueMstKabu",
	}

	masterClientMock.On("GetMasterDataQuery", ctx, expectedReq).Return(nil, expectedErr).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock)

	// --- Execute ---
	result, err := uc.GetStock(ctx, symbol)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())

	masterClientMock.AssertExpectations(t)
}