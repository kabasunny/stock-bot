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

type BalanceClientMock struct {
	mock.Mock
}

func (m *BalanceClientMock) GetZanKaiSummary(ctx context.Context, session *client.Session) (*response.ResZanKaiSummary, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSummary), args.Error(1)
}

func (m *BalanceClientMock) GetGenbutuKabuList(ctx context.Context, session *client.Session) (*response.ResGenbutuKabuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGenbutuKabuList), args.Error(1)
}
func (m *BalanceClientMock) GetShinyouTategyokuList(ctx context.Context, session *client.Session) (*response.ResShinyouTategyokuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResShinyouTategyokuList), args.Error(1)
}
func (m *BalanceClientMock) GetZanKaiKanougaku(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougaku), args.Error(1)
}
func (m *BalanceClientMock) GetZanKaiKanougakuSuii(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougakuSuii), args.Error(1)
}
func (m *BalanceClientMock) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiGenbutuKaitukeSyousai), args.Error(1)
}
func (m *BalanceClientMock) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSinyouSinkidateSyousai), args.Error(1)
}
func (m *BalanceClientMock) GetZanRealHosyoukinRitu(ctx context.Context, session *client.Session, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanRealHosyoukinRitu), args.Error(1)
}
func (m *BalanceClientMock) GetZanShinkiKanoIjiritu(ctx context.Context, session *client.Session, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanShinkiKanoIjiritu), args.Error(1)
}
func (m *BalanceClientMock) GetZanUriKanousuu(ctx context.Context, session *client.Session, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanUriKanousuu), args.Error(1)
}

func TestGetBalance_Success(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}

	// Setup mock
	balanceClientMock := new(BalanceClientMock)

	// Expected response from the client
	apiResponse := &response.ResZanKaiSummary{
		ResultCode:         "0",
		GenbutuKabuKaituke: "1234567",
		SinyouSinkidate:    "2345678",
		HosyouKinritu:      "50.25",
		Syukkin:            "100000",
		OisyouHasseiFlg:    "1", // 1: 発生
	}

	balanceClientMock.On("GetZanKaiSummary", ctx, session).Return(apiResponse, nil).Once()

	// Initialize Usecase
	uc := app.NewBalanceUseCaseImpl(balanceClientMock)

	// Execute
	result, err := uc.GetBalance(ctx, session)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1234567.0, result.AvailableCashForStock)
	assert.Equal(t, 2345678.0, result.AvailableMarginForNewPosition)
	assert.Equal(t, 50.25, result.MarginMaintenanceRate)
	assert.Equal(t, 100000.0, result.WithdrawableCash)
	assert.True(t, result.HasMarginCall)

	balanceClientMock.AssertExpectations(t)
}

func TestGetBalance_ClientError(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}

	// Setup mock
	balanceClientMock := new(BalanceClientMock)
	expectedErr := errors.New("API error")

	balanceClientMock.On("GetZanKaiSummary", ctx, session).Return(nil, expectedErr).Once()

	// Initialize Usecase
	uc := app.NewBalanceUseCaseImpl(balanceClientMock)

	// Execute
	result, err := uc.GetBalance(ctx, session)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())

	balanceClientMock.AssertExpectations(t)
}

func TestGetBalance_ParseError(t *testing.T) {
	ctx := context.Background()
	session := &client.Session{}

	// Setup mock
	balanceClientMock := new(BalanceClientMock)

	// Invalid number format in response
	apiResponse := &response.ResZanKaiSummary{
		ResultCode:         "0",
		GenbutuKabuKaituke: "not-a-number",
		SinyouSinkidate:    "2345678",
		HosyouKinritu:      "50.25",
		Syukkin:            "100000",
		OisyouHasseiFlg:    "1",
	}

	balanceClientMock.On("GetZanKaiSummary", ctx, session).Return(apiResponse, nil).Once()

	// Initialize Usecase
	uc := app.NewBalanceUseCaseImpl(balanceClientMock)

	// Execute
	result, err := uc.GetBalance(ctx, session)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to parse") // Expecting a parsing error

	balanceClientMock.AssertExpectations(t)
}

// go test -v ./internal/app/tests/balance_usecase_impl_test.go
