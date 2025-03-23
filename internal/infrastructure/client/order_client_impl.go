// internal/infrastructure/client/order_client_impl.go

package client

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
)

type orderClientImpl struct {
	client *TachibanaClient
}

func (o *orderClientImpl) NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	fmt.Println("Dummy NewOrder")
	return nil, nil
}
func (o *orderClientImpl) CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error) {
	fmt.Println("Dummy CorrectOrder")
	return nil, nil
}
func (o *orderClientImpl) CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error) {
	fmt.Println("Dummy CancelOrder")
	return nil, nil
}
func (o *orderClientImpl) CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error) {
	fmt.Println("Dummy CancelOrderAll")
	return nil, nil
}
func (o *orderClientImpl) GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error) {
	fmt.Println("Dummy GetOrderList")
	return nil, nil
}
func (o *orderClientImpl) GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error) {
	fmt.Println("Dummy GetOrderListDetail")
	return nil, nil
}
