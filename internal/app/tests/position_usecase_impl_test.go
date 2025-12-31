package tests

import (
	"context"
	"errors"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// NOTE: BalanceClientMock is redefined here for test isolation.
// In a real project, this might be shared.
type PositionUseCaseBalanceClientMock struct {
	mock.Mock
}

func (m *PositionUseCaseBalanceClientMock) GetGenbutuKabuList(ctx context.Context, session *client.Session) (*response.ResGenbutuKabuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGenbutuKabuList), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetShinyouTategyokuList(ctx context.Context, session *client.Session) (*response.ResShinyouTategyokuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResShinyouTategyokuList), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanKaiKanougaku(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougaku), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanKaiKanougakuSuii(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougakuSuii), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanKaiSummary(ctx context.Context, session *client.Session) (*response.ResZanKaiSummary, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSummary), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiGenbutuKaitukeSyousai), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSinyouSinkidateSyousai), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanRealHosyoukinRitu(ctx context.Context, session *client.Session, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanRealHosyoukinRitu), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanShinkiKanoIjiritu(ctx context.Context, session *client.Session, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanShinkiKanoIjiritu), args.Error(1)
}

func (m *PositionUseCaseBalanceClientMock) GetZanUriKanousuu(ctx context.Context, session *client.Session, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanUriKanousuu), args.Error(1)
}

func TestListPositions_All(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}
	balanceClientMock := new(PositionUseCaseBalanceClientMock)

	// --- Mock Data ---
	// 1. Cash Positions
	cashPositionsResponse := &response.ResGenbutuKabuList{
		ResultCode: "0",
		GenbutuKabuList: []response.ResGenbutuKabu{
			{
				UriOrderIssueCode:              "7203", // Toyota
				UriOrderZanKabuSuryou:          "100",
				UriOrderGaisanBokaTanka:        "7000.5",
				UriOrderHyoukaTanka:            "7500.0",
				UriOrderGaisanHyoukaSoneki:     "50000",
				UriOrderGaisanHyoukaSonekiRitu: "7.14",
			},
		},
	}

	// 2. Margin Positions
	marginPositionsResponse := &response.ResShinyouTategyokuList{
		ResultCode: "0",
		SinyouTategyokuList: []response.ResShinyouTategyoku{
			{
				OrderIssueCode:              "6758", // Sony
				OrderBaibaiKubun:            "2",    // '2' is SELL (Short)
				OrderTategyokuSuryou:        "200",
				OrderTategyokuTanka:         "13000",
				OrderHyoukaTanka:            "12500",
				OrderGaisanHyoukaSoneki:     "100000",
				OrderGaisanHyoukaSonekiRitu: "3.84",
				OrderTategyokuDay:           "20251201",
			},
			{
				OrderIssueCode:              "9984", // SoftBank
				OrderBaibaiKubun:            "1",    // '1' is BUY (Long)
				OrderTategyokuSuryou:        "300",
				OrderTategyokuTanka:         "6000",
				OrderHyoukaTanka:            "6200",
				OrderGaisanHyoukaSoneki:     "60000",
				OrderGaisanHyoukaSonekiRitu: "3.33",
				OrderTategyokuDay:           "20251202",
			},
		},
	}

	balanceClientMock.On("GetGenbutuKabuList", ctx, session).Return(cashPositionsResponse, nil).Once()
	balanceClientMock.On("GetShinyouTategyokuList", ctx, session).Return(marginPositionsResponse, nil).Once()

	uc := app.NewPositionUseCaseImpl(balanceClientMock)

	// --- Execute ---
	result, err := uc.ListPositions(ctx, session, "all")

	// --- Assert ---
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)

	// Assert Cash Position (Toyota)
	assert.Equal(t, "7203", result[0].Symbol)
	assert.Equal(t, app.PositionTypeCash, result[0].PositionType)
	assert.Equal(t, 100.0, result[0].Quantity)
	assert.Equal(t, 7000.5, result[0].AverageCost)

	// Assert Margin Short Position (Sony)
	assert.Equal(t, "6758", result[1].Symbol)
	assert.Equal(t, app.PositionTypeMarginShort, result[1].PositionType)
	assert.Equal(t, 200.0, result[1].Quantity)
	assert.Equal(t, 13000.0, result[1].AverageCost)
	assert.Equal(t, "20251201", result[1].OpenedDate)

	// Assert Margin Long Position (SoftBank)
	assert.Equal(t, "9984", result[2].Symbol)
	assert.Equal(t, app.PositionTypeMarginLong, result[2].PositionType)
	assert.Equal(t, 300.0, result[2].Quantity)
	assert.Equal(t, 6000.0, result[2].AverageCost)
	assert.Equal(t, "20251202", result[2].OpenedDate)

	balanceClientMock.AssertExpectations(t)
}

func TestListPositions_ClientError(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}
	balanceClientMock := new(PositionUseCaseBalanceClientMock)
	expectedErr := errors.New("API error")

	// Mock only one of the calls to fail
	balanceClientMock.On("GetGenbutuKabuList", ctx, session).Return(nil, expectedErr).Once()

	uc := app.NewPositionUseCaseImpl(balanceClientMock)

	// --- Execute ---
	result, err := uc.ListPositions(ctx, session, "all")

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())

	// Ensure the second client call was not made
	balanceClientMock.AssertNotCalled(t, "GetShinyouTategyokuList", ctx, session)
	balanceClientMock.AssertExpectations(t)
}

// go test -v ./internal/app/tests/position_usecase_impl_test.go
