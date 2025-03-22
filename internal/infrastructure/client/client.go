// infrastructure/client/client.go
package client

import (
	"context"
	"stock-bot/internal/infrastructure/client/dto/request"
	"stock-bot/internal/infrastructure/client/dto/response"
)

type Client interface {
	Login(ctx context.Context, userID, password string) (*LoginInfo, error)
	Logout(ctx context.Context) error

	// 注文関連
	NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error)
	CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error)
	CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error)
	CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error)
	GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error)
	GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error)

	// 残高・余力関連
	GetGenbutuKabuList(ctx context.Context, req request.ReqGenbutuKabuList) (*response.ResGenbutuKabuList, error)
	GetShinyouTategyokuList(ctx context.Context, req request.ReqShinyouTategyokuList) (*response.ResShinyouTategyokuList, error)
	GetZanKaiKanougaku(ctx context.Context, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error)
	GetZanKaiKanougakuSuii(ctx context.Context, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error)
	GetZanKaiSummary(ctx context.Context, req request.ReqZanKaiSummary) (*response.ResZanKaiSummary, error)
	GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, req request.ReqZanKaiGenbutuKaitukeSyousai) (*response.ResZanKaiGenbutuKaitukeSyousai, error)
	GetZanKaiSinyouSinkidateSyousai(ctx context.Context, req request.ReqZanKaiSinyouSinkidateSyousai) (*response.ResZanKaiSinyouSinkidateSyousai, error)
	GetZanRealHosyoukinRitu(ctx context.Context, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error)
	GetZanShinkiKanoIjiritu(ctx context.Context, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error)
	GetZanUriKanousuu(ctx context.Context, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error)

	// マスタデータ (後で実装)
	// GetMasterData(ctx context.Context, req request.ReqMasterData) (*response.ResMasterData, error)

	// ... 他の必要なメソッド ...
}
