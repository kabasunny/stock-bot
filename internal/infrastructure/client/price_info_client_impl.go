// internal/infrastructure/client/price_info_client_impl.go

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
)

type priceInfoClientImpl struct {
	client *TachibanaClientImpl
}

func (p *priceInfoClientImpl) GetPriceInfo(ctx context.Context, req request.ReqGetPriceInfo) (*response.ResGetPriceInfo, error) {
	if !p.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(p.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMMfdsGetMarketPrice" // 修正
	req.P_no = p.client.getPNo()
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

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get price info with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *priceInfoClientImpl) GetPriceInfoHistory(ctx context.Context, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	if !p.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(p.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMMfdsGetMarketPriceHistory" // 修正
	req.P_no = p.client.getPNo()
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

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get price info history with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfoHistory](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *priceInfoClientImpl) GetPriceInfoHistoryWithPost(ctx context.Context, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	if !p.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(p.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMMfdsGetMarketPriceHistory" // 修正
	req.P_no = p.client.getPNo()
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

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get price info history with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfoHistory](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
