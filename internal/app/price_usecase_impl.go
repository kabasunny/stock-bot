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

// GetHistory retrieves historical price data for a specified stock symbol.
func (uc *PriceUseCaseImpl) GetHistory(ctx context.Context, symbol string, days uint) (*HistoricalPriceResult, error) {
	req := request.ReqGetPriceInfoHistory{
		IssueCode: symbol,
	}

	res, err := uc.priceInfoClient.GetPriceInfoHistory(ctx, uc.session, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history from client: %w", err)
	}

	if res == nil || len(res.CLMMfdsGetMarketPriceHistory) == 0 {
		return &HistoricalPriceResult{
			Symbol: symbol,
			History: []*HistoricalPriceItem{},
		}, nil
	}

	historyItems := make([]*HistoricalPriceItem, 0, len(res.CLMMfdsGetMarketPriceHistory))
	for _, item := range res.CLMMfdsGetMarketPriceHistory {
		open, _ := strconv.ParseFloat(item.PDOPxK, 64)
		high, _ := strconv.ParseFloat(item.PDHPxK, 64)
		low, _ := strconv.ParseFloat(item.PDLPxK, 64)
		close, _ := strconv.ParseFloat(item.PDPPxK, 64)
		volume, _ := strconv.ParseUint(item.PDVxK, 10, 64)

		historyItems = append(historyItems, &HistoricalPriceItem{
			Date:   item.SDate, // Assuming YYYYMMDD and Goa will format
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	// Filter by days if needed (the underlying client might return all history)
	if days > 0 && len(historyItems) > int(days) {
		historyItems = historyItems[:days]
	}

	return &HistoricalPriceResult{
		Symbol: symbol,
		History: historyItems,
	}, nil
}