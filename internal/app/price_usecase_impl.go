package app

import (
	"context"
	"fmt"
	"strconv"
	"stock-bot/gen/price"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"
)

// PriceUseCaseImpl implements the PriceUseCase interface.
type PriceUseCaseImpl struct {
	priceInfoClient client.PriceInfoClient
	session         *client.Session // Sessionを追加
}

// NewPriceUseCaseImpl creates a new PriceUseCaseImpl.
func NewPriceUseCaseImpl(priceInfoClient client.PriceInfoClient, session *client.Session) *PriceUseCaseImpl {
	return &PriceUseCaseImpl{
		priceInfoClient: priceInfoClient,
		session:         session,
	}
}

// Get retrieves the current price for a specified stock symbol.
func (uc *PriceUseCaseImpl) Get(ctx context.Context, symbol string) (*price.StockbotPrice, error) {
	req := request.ReqGetPriceInfo{
		CLMID:           "CLMMfdsGetMarketPrice",
		TargetIssueCode: symbol,
		TargetColumn:    "CurrentPrice,Timestamp", // 必要な情報を指定
	}

	resGetPriceInfo, err := uc.priceInfoClient.GetPriceInfo(ctx, uc.session, req)
	_ = response.ResGetPriceInfo{} // Workaround for "imported and not used" error
	if err != nil {
		return nil, fmt.Errorf("failed to get price info from client: %w", err)
	}

	if resGetPriceInfo == nil || len(resGetPriceInfo.CLMMfdsMarketPrice) == 0 {
		return nil, fmt.Errorf("no price info found for symbol %s", symbol)
	}

	item := resGetPriceInfo.CLMMfdsMarketPrice[0] // 最初のアイテムを使用すると仮定

	currentPriceStr, ok := item.Values["CurrentPrice"]
	if !ok {
		return nil, fmt.Errorf("CurrentPrice not found in response for symbol %s", symbol)
	}

	timestampStr, ok := item.Values["Timestamp"]
	if !ok {
		return nil, fmt.Errorf("Timestamp not found in response for symbol %s", symbol)
	}

	currentPrice, err := strconv.ParseFloat(currentPriceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price '%s' for symbol %s: %w", currentPriceStr, symbol, err)
	}

	return &price.StockbotPrice{
		Symbol:    symbol,
		Price:     currentPrice,
		Timestamp: timestampStr,
	}, nil
}