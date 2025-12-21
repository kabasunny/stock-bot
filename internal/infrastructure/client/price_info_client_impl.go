// internal/infrastructure/client/price_info_client_impl.go

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv" // 追加
	"time"

	"stock-bot/internal/infrastructure/client/dto/price/request"
	"stock-bot/internal/infrastructure/client/dto/price/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
)
type priceInfoClientImpl struct {
	client *TachibanaClientImpl
}

func (p *priceInfoClientImpl) GetPriceInfo(ctx context.Context, session *Session, req request.ReqGetPriceInfo) (*response.ResGetPriceInfo, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMMfdsGetMarketPrice"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
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

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get price info with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *priceInfoClientImpl) GetPriceInfoHistory(ctx context.Context, session *Session, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMMfdsGetMarketPriceHistory"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
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

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get price info history with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfoHistory](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (p *priceInfoClientImpl) GetPriceInfoHistoryWithPost(ctx context.Context, session *Session, req request.ReqGetPriceInfoHistory) (*response.ResGetPriceInfoHistory, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMMfdsGetMarketPriceHistory"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
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

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get price info history with post failed")
	}

	res, err := ConvertResponse[response.ResGetPriceInfoHistory](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
