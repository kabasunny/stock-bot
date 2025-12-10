// internal/infrastructure/client/order_client_impl.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
)

type orderClientImpl struct {
	client *TachibanaClientImpl
}

func (o *orderClientImpl) NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuNewOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	var hensaiDataJSON string
	if len(req.CLMKabuHensaiData) > 0 {
		hensaiBytes, err := json.Marshal(req.CLMKabuHensaiData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal CLMKabuHensaiData to JSON")
		}
		hensaiDataJSON = string(hensaiBytes)
	}

	tempReq := req
	tempReq.CLMKabuHensaiData = nil

	params, err := structToMapString(tempReq)
	if err != nil {
		return nil, err
	}

	if hensaiDataJSON != "" {
		params["aCLMKabuHensaiData"] = hensaiDataJSON
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		if k == "aCLMKabuHensaiData" {
			buf.WriteString(fmt.Sprintf(`"%s":%s`, k, v))
		} else {
			buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
		}
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "new order failed")
	}

	res, err := ConvertResponse[response.ResNewOrder](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderClientImpl) CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCorrectOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "correct order failed")
	}

	res, err := ConvertResponse[response.ResCorrectOrder](respMap)
	if err != nil {
		return nil, err
	}
	if res.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s (errno: %s)", res.ResultText, res.ResultCode)
	}

	return res, nil
}

func (o *orderClientImpl) CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCancelOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "cancel order failed")
	}

	res, err := ConvertResponse[response.ResCancelOrder](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderClientImpl) CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCancelOrderAll"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "cancel all order failed")
	}

	res, err := ConvertResponse[response.ResCancelOrderAll](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderClientImpl) GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMOrderList"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get order list failed")
	}

	res, err := ConvertResponse[response.ResOrderList](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderClientImpl) GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMOrderListDetail"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(o.client.httpClient, httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get order list detail failed")
	}

	res, err := ConvertResponse[response.ResOrderListDetail](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
