// internal/infrastructure/client/price_info_client_impl.go

package client

import (
	"context"
	"log/slog"

	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"
	_ "stock-bot/internal/logger"
)

type priceInfoClientImpl struct {
	client *TachibanaClientImpl
}

func (p *priceInfoClientImpl) GetPriceInfo(ctx context.Context, req request.ReqGetPriceInfo) (*response.ResGetPriceInfo, error) {
	slog.Info("Dummy GetPriceInfo")
	return nil, nil
}
func (p *priceInfoClientImpl) GetPriceInfoHistory(ctx context.Context, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	slog.Info("Dummy GetPriceInfoHistory")
	return nil, nil
}
