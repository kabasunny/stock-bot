package app

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
	"strconv"
)

// masterUseCaseImpl implements the MasterUseCase interface.
type masterUseCaseImpl struct {
	masterClient client.MasterDataClient
	masterRepo   repository.MasterRepository
}

// NewMasterUseCaseImpl creates a new MasterUseCase.
func NewMasterUseCaseImpl(masterClient client.MasterDataClient, masterRepo repository.MasterRepository) MasterUseCase {
	return &masterUseCaseImpl{
		masterClient: masterClient,
		masterRepo:   masterRepo,
	}
}

// DownloadAndStoreMasterData downloads all master data and stores the watched stocks into the database.
func (uc *masterUseCaseImpl) DownloadAndStoreMasterData(ctx context.Context) error {
	slog.Info("Starting master data download...")
	// 1. Call DownloadMasterData
	res, err := uc.masterClient.DownloadMasterData(ctx, request.ReqDownloadMaster{})
	if err != nil {
		return fmt.Errorf("failed to download master data: %w", err)
	}
	slog.Info("Master data download completed.", "system_status", res.SystemStatus.SystemStatus)

	// MarketMasterを銘柄コードで検索できるようにマップに変換
	marketMasterMap := make(map[string]response.ResStockMarketMaster)
	for _, mm := range res.StockMarketMaster {
		marketMasterMap[mm.IssueCode] = mm
	}

	var modelsToUpsert []*model.StockMaster
	slog.Info("Converting stock master data to domain models...", "count", len(res.StockMaster))
	for _, sm := range res.StockMaster {
		// string to int/float conversions with error handling
		tradingUnit, err := strconv.Atoi(sm.TradingUnit)
		if err != nil {
			slog.Warn("Failed to parse TradingUnit, setting 0", "value", sm.TradingUnit, "error", err)
			tradingUnit = 0
		}
		listedShares, err := strconv.ParseInt(sm.ListedSharesOutstanding, 10, 64)
		if err != nil {
			slog.Warn("Failed to parse ListedSharesOutstanding, setting 0", "value", sm.ListedSharesOutstanding, "error", err)
			listedShares = 0
		}

		var upperLimit, lowerLimit float64
		if mm, ok := marketMasterMap[sm.IssueCode]; ok {
			upperLimit, err = strconv.ParseFloat(mm.UpperLimit, 64)
			if err != nil {
				slog.Warn("Failed to parse UpperLimit", "value", mm.UpperLimit, "error", err)
			}
			lowerLimit, err = strconv.ParseFloat(mm.LowerLimit, 64)
			if err != nil {
				slog.Warn("Failed to parse LowerLimit", "value", mm.LowerLimit, "error", err)
			}
		}

		m := &model.StockMaster{
			IssueCode:               sm.IssueCode,
			IssueName:               sm.IssueName,
			IssueNameShort:          sm.IssueNameShort,
			IssueNameKana:           sm.IssueNameKana,
			IssueNameEnglish:        sm.IssueNameEnglish,
			MarketCode:              sm.PreferredMarket,
			IndustryCode:            sm.IndustryCode,
			IndustryName:            sm.IndustryName,
			TradingUnit:             tradingUnit,
			ListedSharesOutstanding: listedShares,
			UpperLimit:              upperLimit,
			LowerLimit:              lowerLimit,
		}
		modelsToUpsert = append(modelsToUpsert, m)
	}
	slog.Info("Finished converting stocks.", "upsert_count", len(modelsToUpsert))

	// 4. Call uc.masterRepo.UpsertStockMasters
	if len(modelsToUpsert) > 0 {
		slog.Info("Upserting stock masters to database...")
		if err := uc.masterRepo.UpsertStockMasters(ctx, modelsToUpsert); err != nil {
			return fmt.Errorf("failed to upsert stock masters: %w", err)
		}
		slog.Info("Successfully upserted stock masters.")
	}

	// 5. Process and Upsert TickRules
	var tickRulesToUpsert []*model.TickRule
	for _, rt := range res.TickRule {
		tickRule := &model.TickRule{
			TickUnitNumber: rt.TickUnitNumber,
			ApplicableDate: rt.ApplicableDate,
			TickLevels:     []model.TickLevel{},
		}

		levels := []struct {
			BasePrice string
			TickValue string
		}{
			{rt.BasePrice1, rt.TickValue1}, {rt.BasePrice2, rt.TickValue2},
			{rt.BasePrice3, rt.TickValue3}, {rt.BasePrice4, rt.TickValue4},
			{rt.BasePrice5, rt.TickValue5}, {rt.BasePrice6, rt.TickValue6},
			{rt.BasePrice7, rt.TickValue7}, {rt.BasePrice8, rt.TickValue8},
			{rt.BasePrice9, rt.TickValue9}, {rt.BasePrice10, rt.TickValue10},
			{rt.BasePrice11, rt.TickValue11}, {rt.BasePrice12, rt.TickValue12},
			{rt.BasePrice13, rt.TickValue13}, {rt.BasePrice14, rt.TickValue14},
			{rt.BasePrice15, rt.TickValue15}, {rt.BasePrice16, rt.TickValue16},
			{rt.BasePrice17, rt.TickValue17}, {rt.BasePrice18, rt.TickValue18},
			{rt.BasePrice19, rt.TickValue19}, {rt.BasePrice20, rt.TickValue20},
		}

		var upperPrice float64 = 0.0
		for _, levelData := range levels {
			if levelData.BasePrice == "" || levelData.TickValue == "" {
				continue
			}

			basePrice, errP := strconv.ParseFloat(levelData.BasePrice, 64)
			tickValue, errT := strconv.ParseFloat(levelData.TickValue, 64)

			if errP != nil || errT != nil {
				slog.Warn("Failed to parse tick level data", "rule", rt.TickUnitNumber, "basePrice", levelData.BasePrice, "tickValue", levelData.TickValue)
				continue
			}

			tickLevel := model.TickLevel{
				LowerPrice: upperPrice,
				UpperPrice: basePrice,
				TickValue:  tickValue,
			}
			tickRule.TickLevels = append(tickRule.TickLevels, tickLevel)
			upperPrice = basePrice // 次のレベルの下限値は、現在のレベルの上限値
		}
		tickRulesToUpsert = append(tickRulesToUpsert, tickRule)
	}

	if len(tickRulesToUpsert) > 0 {
		slog.Info("Upserting tick rules to database...")
		if err := uc.masterRepo.UpsertTickRules(ctx, tickRulesToUpsert); err != nil {
			return fmt.Errorf("failed to upsert tick rules: %w", err)
		}
		slog.Info("Successfully upserted tick rules.")
	}

	return nil
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
