package client

import (
	"context"
	"fmt"
	balance_request "stock-bot/internal/infrastructure/client/dto/balance/request"
	balance_response "stock-bot/internal/infrastructure/client/dto/balance/response"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	order_response "stock-bot/internal/infrastructure/client/dto/order/response"
	price_request "stock-bot/internal/infrastructure/client/dto/price/request"
	price_response "stock-bot/internal/infrastructure/client/dto/price/response"
)

// TachibanaUnifiedClientAdapter は TachibanaUnifiedClient を既存のインターフェースに適合させるアダプター
type TachibanaUnifiedClientAdapter struct {
	unifiedClient *TachibanaUnifiedClient
}

// NewTachibanaUnifiedClientAdapter は新しいアダプターを作成します
func NewTachibanaUnifiedClientAdapter(unifiedClient *TachibanaUnifiedClient) *TachibanaUnifiedClientAdapter {
	return &TachibanaUnifiedClientAdapter{
		unifiedClient: unifiedClient,
	}
}

// --- BalanceClient interface implementation ---

func (a *TachibanaUnifiedClientAdapter) GetGenbutuKabuList(ctx context.Context, session *Session) (*balance_response.ResGenbutuKabuList, error) {
	return a.unifiedClient.GetGenbutuKabuList(ctx)
}

func (a *TachibanaUnifiedClientAdapter) GetShinyouTategyokuList(ctx context.Context, session *Session) (*balance_response.ResShinyouTategyokuList, error) {
	return a.unifiedClient.GetShinyouTategyokuList(ctx)
}

func (a *TachibanaUnifiedClientAdapter) GetZanKaiSummary(ctx context.Context, session *Session) (*balance_response.ResZanKaiSummary, error) {
	return a.unifiedClient.GetZanKaiSummary(ctx)
}

func (a *TachibanaUnifiedClientAdapter) GetZanKaiKanougaku(ctx context.Context, session *Session, req balance_request.ReqZanKaiKanougaku) (*balance_response.ResZanKaiKanougaku, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanKaiKanougaku not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanKaiKanougakuSuii(ctx context.Context, session *Session, req balance_request.ReqZanKaiKanougakuSuii) (*balance_response.ResZanKaiKanougakuSuii, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanKaiKanougakuSuii not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, session *Session, tradingDay int) (*balance_response.ResZanKaiGenbutuKaitukeSyousai, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanKaiGenbutuKaitukeSyousai not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, session *Session, tradingDay int) (*balance_response.ResZanKaiSinyouSinkidateSyousai, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanKaiSinyouSinkidateSyousai not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanRealHosyoukinRitu(ctx context.Context, session *Session, req balance_request.ReqZanRealHosyoukinRitu) (*balance_response.ResZanRealHosyoukinRitu, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanRealHosyoukinRitu not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanShinkiKanoIjiritu(ctx context.Context, session *Session, req balance_request.ReqZanShinkiKanoIjiritu) (*balance_response.ResZanShinkiKanoIjiritu, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanShinkiKanoIjiritu not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetZanUriKanousuu(ctx context.Context, session *Session, req balance_request.ReqZanUriKanousuu) (*balance_response.ResZanUriKanousuu, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetZanUriKanousuu not implemented in unified client")
}

// --- OrderClient interface implementation ---

func (a *TachibanaUnifiedClientAdapter) NewOrder(ctx context.Context, session *Session, params NewOrderParams) (*order_response.ResNewOrder, error) {
	return a.unifiedClient.NewOrder(ctx, params)
}

func (a *TachibanaUnifiedClientAdapter) CorrectOrder(ctx context.Context, session *Session, params CorrectOrderParams) (*order_response.ResCorrectOrder, error) {
	return a.unifiedClient.CorrectOrder(ctx, params)
}

func (a *TachibanaUnifiedClientAdapter) CancelOrder(ctx context.Context, session *Session, params CancelOrderParams) (*order_response.ResCancelOrder, error) {
	return a.unifiedClient.CancelOrder(ctx, params)
}

func (a *TachibanaUnifiedClientAdapter) GetOrderList(ctx context.Context, session *Session, req order_request.ReqOrderList) (*order_response.ResOrderList, error) {
	return a.unifiedClient.GetOrderList(ctx, req)
}

func (a *TachibanaUnifiedClientAdapter) CancelOrderAll(ctx context.Context, session *Session, params CancelOrderAllParams) (*order_response.ResCancelOrderAll, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("CancelOrderAll not implemented in unified client")
}

func (a *TachibanaUnifiedClientAdapter) GetOrderListDetail(ctx context.Context, session *Session, req order_request.ReqOrderListDetail) (*order_response.ResOrderListDetail, error) {
	// TODO: Implement in unified client
	return nil, fmt.Errorf("GetOrderListDetail not implemented in unified client")
}

// --- PriceInfoClient interface implementation ---

func (a *TachibanaUnifiedClientAdapter) GetPriceInfo(ctx context.Context, session *Session, req price_request.ReqGetPriceInfo) (*price_response.ResGetPriceInfo, error) {
	return a.unifiedClient.GetPriceInfo(ctx, req)
}

func (a *TachibanaUnifiedClientAdapter) GetPriceInfoHistory(ctx context.Context, session *Session, req price_request.ReqGetPriceInfoHistory) (*price_response.ResGetPriceInfoHistory, error) {
	return a.unifiedClient.GetPriceInfoHistory(ctx, req)
}
