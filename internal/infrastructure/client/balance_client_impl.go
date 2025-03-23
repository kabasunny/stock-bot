// internal/infrastructure/client/balance_client_impl.go
package client

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
)

type balanceClientImpl struct {
	client *TachibanaClient
}

func (b *balanceClientImpl) GetGenbutuKabuList(ctx context.Context, req request.ReqGenbutuKabuList) (*response.ResGenbutuKabuList, error) {
	fmt.Println("Dummy GetGenbutuKabuList")
	return nil, nil
}

func (b *balanceClientImpl) GetShinyouTategyokuList(ctx context.Context, req request.ReqShinyouTategyokuList) (*response.ResShinyouTategyokuList, error) {
	fmt.Println("Dummy GetShinyouTategyokuList")
	return nil, nil
}

func (b *balanceClientImpl) GetZanKaiKanougaku(ctx context.Context, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	fmt.Println("Dummy GetZanKaiKanougaku")
	return nil, nil
}
func (b *balanceClientImpl) GetZanKaiKanougakuSuii(ctx context.Context, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	fmt.Println("Dummy GetZanKaiKanougakuSuii")
	return nil, nil
}
func (b *balanceClientImpl) GetZanKaiSummary(ctx context.Context, req request.ReqZanKaiSummary) (*response.ResZanKaiSummary, error) {
	fmt.Println("Dummy GetZanKaiSummary")
	return nil, nil
}
func (b *balanceClientImpl) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, req request.ReqZanKaiGenbutuKaitukeSyousai) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	fmt.Println("Dummy GetZanKaiGenbutuKaitukeSyousai")
	return nil, nil
}
func (b *balanceClientImpl) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, req request.ReqZanKaiSinyouSinkidateSyousai) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	fmt.Println("Dummy GetZanKaiSinyouSinkidateSyousai")
	return nil, nil
}
func (b *balanceClientImpl) GetZanRealHosyoukinRitu(ctx context.Context, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	fmt.Println("Dummy GetZanRealHosyoukinRitu")
	return nil, nil
}
func (b *balanceClientImpl) GetZanShinkiKanoIjiritu(ctx context.Context, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	fmt.Println("Dummy GetZanShinkiKanoIjiritu")
	return nil, nil
}
func (b *balanceClientImpl) GetZanUriKanousuu(ctx context.Context, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	fmt.Println("Dummy GetZanUriKanousuu")
	return nil, nil
}
