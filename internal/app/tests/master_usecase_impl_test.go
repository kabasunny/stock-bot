package tests

import (
	"context"
	"errors"
	"stock-bot/domain/model"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
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
func (m *MasterDataClientMock) DownloadMasterData(ctx context.Context, session *client.Session, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResDownloadMaster), args.Error(1)
}
func (m *MasterDataClientMock) GetMasterDataQuery(ctx context.Context, session *client.Session, req request.ReqGetMasterData) (*response.ResGetMasterData, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGetMasterData), args.Error(1)
}
func (m *MasterDataClientMock) GetNewsHeader(ctx context.Context, session *client.Session, req request.ReqGetNewsHead) (*response.ResGetNewsHeader, error) {
	args := m.Called(ctx, session, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetNewsBody(ctx context.Context, session *client.Session, req request.ReqGetNewsBody) (*response.ResGetNewsBody, error) {
	args := m.Called(ctx, session, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetIssueDetail(ctx context.Context, session *client.Session, req request.ReqGetIssueDetail) (*response.ResGetIssueDetail, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGetIssueDetail), args.Error(1)
}
func (m *MasterDataClientMock) GetMarginInfo(ctx context.Context, session *client.Session, req request.ReqGetMarginInfo) (*response.ResGetMarginInfo, error) {
	args := m.Called(ctx, session, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetCreditInfo(ctx context.Context, session *client.Session, req request.ReqGetCreditInfo) (*response.ResGetCreditInfo, error) {
	args := m.Called(ctx, session, req)
	return nil, args.Error(1)
}
func (m *MasterDataClientMock) GetMarginPremiumInfo(ctx context.Context, session *client.Session, req request.ReqGetMarginPremiumInfo) (*response.ResGetMarginPremiumInfo, error) {
	args := m.Called(ctx, session, req)
	return nil, args.Error(1)
}

// MasterRepositoryMock implements the repository.MasterRepository interface for testing.
type MasterRepositoryMock struct {
	mock.Mock
}

func (m *MasterRepositoryMock) Save(ctx context.Context, entity interface{}) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}
func (m *MasterRepositoryMock) SaveAll(ctx context.Context, entities []interface{}) error {
	args := m.Called(ctx, entities)
	return args.Error(0)
}
func (m *MasterRepositoryMock) FindByIssueCode(ctx context.Context, issueCode string, entityType string) (interface{}, error) {
	args := m.Called(ctx, issueCode, entityType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}
func (m *MasterRepositoryMock) UpsertStockMasters(ctx context.Context, stocks []*model.StockMaster) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}
func (m *MasterRepositoryMock) UpsertTickRules(ctx context.Context, tickRules []*model.TickRule) error {
	args := m.Called(ctx, tickRules)
	return args.Error(0)
}

func TestGetStock_Success(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	masterRepoMock := new(MasterRepositoryMock)
	symbol := "7203"

	// --- Mock Data ---
	mockStockMaster := &model.StockMaster{
		IssueCode:    symbol,
		IssueName:    "トヨタ自動車",
		IssueNameKana: "トヨタジドウシャ",
		MarketCode:   "東証プライム",
		IndustryCode: "3600",
		IndustryName: "輸送用機器",
	}

	masterRepoMock.On("FindByIssueCode", ctx, symbol, "StockMaster").Return(mockStockMaster, nil).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock, masterRepoMock)

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

	masterRepoMock.AssertExpectations(t)
}

func TestGetStock_NotFound(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	masterRepoMock := new(MasterRepositoryMock)
	symbol := "9999" // Symbol not in mock data

	masterRepoMock.On("FindByIssueCode", ctx, symbol, "StockMaster").Return(nil, nil).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock, masterRepoMock)

	// --- Execute ---
	result, err := uc.GetStock(ctx, symbol)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, app.ErrNotFound))

	masterRepoMock.AssertExpectations(t)
}

func TestGetStock_RepoError(t *testing.T) {
	ctx := context.Background()
	masterClientMock := new(MasterDataClientMock)
	masterRepoMock := new(MasterRepositoryMock)
	symbol := "7203"
	expectedErr := errors.New("repository error")

	masterRepoMock.On("FindByIssueCode", ctx, symbol, "StockMaster").Return(nil, expectedErr).Once()

	uc := app.NewMasterUseCaseImpl(masterClientMock, masterRepoMock)

	// --- Execute ---
	result, err := uc.GetStock(ctx, symbol)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())

	masterRepoMock.AssertExpectations(t)
}

func TestDownloadAndStoreMasterData_Success(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}
	masterClientMock := new(MasterDataClientMock)
	masterRepoMock := new(MasterRepositoryMock)

	// --- Mock Data ---
	dummyMasterData := &response.ResDownloadMaster{
		SystemStatus: response.ResSystemStatus{SystemStatus: "1"},
		StockMaster: []response.ResStockMaster{
			{
				IssueCode:               "7203",
				IssueName:               "トヨタ自動車",
				IssueNameShort:          "トヨタ",
				IssueNameKana:           "トヨタジドウシャ",
				IssueNameEnglish:        "TOYOTA MOTOR",
				PreferredMarket:         "00",
				IndustryCode:            "3700",
				IndustryName:            "輸送用機器",
				TradingUnit:             "100",
				ListedSharesOutstanding: "3262997492",
			},
			{
				IssueCode:               "9984",
				IssueName:               "ソフトバンクグループ",
				IssueNameShort:          "SBG",
				IssueNameKana:           "ソフトバンクグループ",
				IssueNameEnglish:        "SOFTBANK GROUP",
				PreferredMarket:         "00",
				IndustryCode:            "5250",
				IndustryName:            "情報・通信業",
				TradingUnit:             "100",
				ListedSharesOutstanding: "1452632314",
			},
		},
		StockMarketMaster: []response.ResStockMarketMaster{
			{
				IssueCode:  "7203",
				UpperLimit: "10000.0",
				LowerLimit: "1000.0",
			},
			{
				IssueCode:  "9984",
				UpperLimit: "8000.0",
				LowerLimit: "500.0",
			},
		},
		TickRule: []response.ResTickRule{
			{
				TickUnitNumber: "101",
				ApplicableDate: "20140101",
				BasePrice1:     "3000.0",
				TickValue1:     "1.0",
				BasePrice2:     "5000.0",
				TickValue2:     "5.0",
			},
		},
	}

	// --- Mock Setup ---
	masterClientMock.On("DownloadMasterData", ctx, session, request.ReqDownloadMaster{}).Return(dummyMasterData, nil).Once()

	var capturedStocks []*model.StockMaster
	masterRepoMock.On("UpsertStockMasters", ctx, mock.AnythingOfType("[]*model.StockMaster")).Run(func(args mock.Arguments) {
		capturedStocks = args.Get(1).([]*model.StockMaster)
	}).Return(nil).Once()

	masterRepoMock.On("UpsertTickRules", ctx, mock.AnythingOfType("[]*model.TickRule")).Return(nil).Once()

	// --- Execute ---
	uc := app.NewMasterUseCaseImpl(masterClientMock, masterRepoMock)
	err := uc.DownloadAndStoreMasterData(ctx, session)

	// --- Assert ---
	assert.NoError(t, err)
	masterClientMock.AssertExpectations(t)
	masterRepoMock.AssertExpectations(t)

	// Assert captured data
	assert.Len(t, capturedStocks, 2)
	// Check first stock
	assert.Equal(t, "7203", capturedStocks[0].IssueCode)
	assert.Equal(t, "トヨタ自動車", capturedStocks[0].IssueName)
	assert.Equal(t, "トヨタ", capturedStocks[0].IssueNameShort)
	assert.Equal(t, "輸送用機器", capturedStocks[0].IndustryName)
	assert.Equal(t, 100, capturedStocks[0].TradingUnit)
	assert.Equal(t, int64(3262997492), capturedStocks[0].ListedSharesOutstanding)
	assert.Equal(t, 10000.0, capturedStocks[0].UpperLimit)
	// Check second stock
	assert.Equal(t, "9984", capturedStocks[1].IssueCode)
	assert.Equal(t, "ソフトバンクグループ", capturedStocks[1].IssueName)
	assert.Equal(t, 100, capturedStocks[1].TradingUnit)
	assert.Equal(t, 8000.0, capturedStocks[1].UpperLimit)
}
