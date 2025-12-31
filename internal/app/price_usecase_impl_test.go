package app_test

import (
	"context"
	"errors"
	"stock-bot/internal/app"
	"stock-bot/internal/app/mocks"
	client_request "stock-bot/internal/infrastructure/client/dto/price/request"
	client_response "stock-bot/internal/infrastructure/client/dto/price/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPriceUseCase_Get(t *testing.T) {
	mockPriceInfoClient := new(mocks.PriceInfoClient)
	// PriceUseCaseImpl はまだ実装していないので、ここでは仮のコンストラクタを想定
	// 実装時に app.NewPriceUseCaseImpl に置き換える
	useCase := app.NewPriceUseCaseImpl(mockPriceInfoClient, nil)

	ctx := context.Background()
	symbol := "6758"
	expectedPrice := 1000.5
	expectedTimestamp := "2025-12-22T09:00:00Z"

	// 成功ケース
	t.Run("成功ケース: 指定された銘柄の価格を取得できること", func(t *testing.T) {
		// client.GetPriceInfo が返すレスポンスをモック
		mockResponse := &client_response.ResGetPriceInfo{
			CLMID: "CLMMfdsGetMarketPrice",
			CLMMfdsMarketPrice: []client_response.ResMarketPriceInfoItem{
				{
					IssueCode: symbol,
					Values: map[string]string{
						"CurrentPrice": "1000.5", // 文字列として返す
						"Timestamp":    expectedTimestamp,
					},
				},
			},
		}

		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(mockResponse, nil).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, symbol, res.Symbol)
		assert.Equal(t, expectedPrice, res.Price)
		assert.Equal(t, expectedTimestamp, res.Timestamp)
		mockPriceInfoClient.AssertExpectations(t)
	})

	// エラーケース
	t.Run("エラーケース: 価格取得に失敗した場合、エラーを返すこと", func(t *testing.T) {
		expectedErr := errors.New("APIエラー")
		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(nil, expectedErr).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockPriceInfoClient.AssertExpectations(t)
	})

	// エラーケース: レスポンスのCLMMfdsMarketPriceが空の場合
	t.Run("エラーケース: CLMMfdsMarketPriceが空の場合、エラーを返すこと", func(t *testing.T) {
		mockResponse := &client_response.ResGetPriceInfo{
			CLMID:              "CLMMfdsGetMarketPrice",
			CLMMfdsMarketPrice: []client_response.ResMarketPriceInfoItem{},
		}
		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(mockResponse, nil).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "no price info found for symbol")
		mockPriceInfoClient.AssertExpectations(t)
	})

	// エラーケース: ValuesマップにCurrentPriceキーがない場合
	t.Run("エラーケース: ValuesマップにCurrentPriceキーがない場合、エラーを返すこと", func(t *testing.T) {
		mockResponse := &client_response.ResGetPriceInfo{
			CLMID: "CLMMfdsGetMarketPrice",
			CLMMfdsMarketPrice: []client_response.ResMarketPriceInfoItem{
				{
					IssueCode: symbol,
					Values: map[string]string{
						"Timestamp": expectedTimestamp, // CurrentPrice missing
					},
				},
			},
		}
		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(mockResponse, nil).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "CurrentPrice not found in response")
		mockPriceInfoClient.AssertExpectations(t)
	})

	// エラーケース: ValuesマップにTimestampキーがない場合
	t.Run("エラーケース: ValuesマップにTimestampキーがない場合、エラーを返すこと", func(t *testing.T) {
		mockResponse := &client_response.ResGetPriceInfo{
			CLMID: "CLMMfdsGetMarketPrice",
			CLMMfdsMarketPrice: []client_response.ResMarketPriceInfoItem{
				{
					IssueCode: symbol,
					Values: map[string]string{
						"CurrentPrice": "1000.5", // Timestamp missing
					},
				},
			},
		}
		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(mockResponse, nil).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Timestamp not found in response")
		mockPriceInfoClient.AssertExpectations(t)
	})

	// エラーケース: 価格の文字列変換失敗
	t.Run("エラーケース: 価格の文字列変換失敗", func(t *testing.T) {
		mockResponse := &client_response.ResGetPriceInfo{
			CLMID: "CLMMfdsGetMarketPrice",
			CLMMfdsMarketPrice: []client_response.ResMarketPriceInfoItem{
				{
					IssueCode: symbol,
					Values: map[string]string{
						"CurrentPrice": "invalid_price",
						"Timestamp":    expectedTimestamp,
					},
				},
			},
		}
		mockPriceInfoClient.On("GetPriceInfo",
			mock.Anything,                            // context.Background()
			mock.Anything,                            // session *Session
			client_request.ReqGetPriceInfo{CLMID: "CLMMfdsGetMarketPrice", TargetIssueCode: symbol, TargetColumn: "CurrentPrice,Timestamp"}, // req request.ReqGetPriceInfo
		).Return(mockResponse, nil).Once()

		res, err := useCase.Get(ctx, symbol)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "failed to parse price")
		mockPriceInfoClient.AssertExpectations(t)
	})
}
