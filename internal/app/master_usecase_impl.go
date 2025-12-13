package app

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
)

// masterUseCaseImpl implements the MasterUseCase interface.
type masterUseCaseImpl struct {
	masterClient client.MasterDataClient
}

// NewMasterUseCaseImpl creates a new MasterUseCase.
func NewMasterUseCaseImpl(masterClient client.MasterDataClient) MasterUseCase {
	return &masterUseCaseImpl{masterClient: masterClient}
}

// GetStock retrieves basic master data for a single stock.
func (uc *masterUseCaseImpl) GetStock(ctx context.Context, symbol string) (*StockMasterResult, error) {
	// Request the entire Stock Master list for now (inefficient, to be improved with caching)
	req := request.ReqGetMasterData{
		TargetCLMID: "CLMIssueMstKabu", // Request Stock Master data
		// TargetColumn can be left empty to get all columns, or specified as needed.
	}

	res, err := uc.masterClient.GetMasterDataQuery(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get master data query from client: %w", err)
	}

	if res == nil || len(res.StockMaster) == 0 {
		return nil, ErrNotFound
	}

	// Find the requested symbol in the returned list
	var stockFound *response.ResStockMaster
	for i := range res.StockMaster {
		if res.StockMaster[i].IssueCode == symbol {
			stockFound = &res.StockMaster[i]
			break
		}
	}

	if stockFound == nil {
		return nil, ErrNotFound
	}

	result := &StockMasterResult{
		Symbol:       stockFound.IssueCode,
		Name:         stockFound.IssueName,
		NameKana:     stockFound.IssueNameKana,
		Market:       stockFound.PreferredMarket,
		IndustryCode: stockFound.IndustryCode,
		IndustryName: stockFound.IndustryName,
	}

	return result, nil
}
