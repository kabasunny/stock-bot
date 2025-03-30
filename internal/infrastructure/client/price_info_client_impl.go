// internal/infrastructure/client/price_info_client_impl.go

package client

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"

	"go.uber.org/zap"
)

type priceInfoClientImpl struct {
	client *TachibanaClientImpl
	logger *zap.Logger // 追加
}

func (p *priceInfoClientImpl) GetPriceInfo(ctx context.Context, req request.ReqGetPriceInfo) (*response.ResGetPriceInfo, error) {
	fmt.Println("Dummy GetPriceInfo")
	return nil, nil
}
func (p *priceInfoClientImpl) GetPriceInfoHistory(ctx context.Context, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	fmt.Println("Dummy GetPriceInfoHistory")
	return nil, nil
}
