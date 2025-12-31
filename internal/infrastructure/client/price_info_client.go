// internal/infrastructure/client/price_info_client.go
package client

import (
	"context"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"
)

// PriceInfoClient は、時価情報関連の API を扱うインターフェース
type PriceInfoClient interface {
	// GetPriceInfo は、指定した銘柄の時価情報を取得
	GetPriceInfo(ctx context.Context, session *Session, req request.ReqGetPriceInfo) (*response.ResGetPriceInfo, error)
	// GetPriceInfoHistory は、指定した銘柄の過去の時価情報（四本値、出来高など）を取得
	GetPriceInfoHistory(ctx context.Context, session *Session, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error)
}
