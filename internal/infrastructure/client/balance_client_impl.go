// internal/infrastructure/client/balance_client_impl.go
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
)

type balanceClientImpl struct {
	client *TachibanaClientImpl
}

func (b *balanceClientImpl) GetGenbutuKabuList(ctx context.Context, session *Session) (*response.ResGenbutuKabuList, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req := request.ReqGenbutuKabuList{}
	req.CLMID = "CLMGenbutuKabuList"
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
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v)) // すべての値を文字列として扱う
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
		return nil, errors.Wrap(err, "GetGenbutuKabuList failed")
	}

	res, err := ConvertResponse[response.ResGenbutuKabuList](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (b *balanceClientImpl) GetShinyouTategyokuList(ctx context.Context, session *Session) (*response.ResShinyouTategyokuList, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req := request.ReqShinyouTategyokuList{}
	req.CLMID = "CLMShinyouTategyokuList"
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
		return nil, errors.Wrap(err, "GetShinyouTategyokuList failed")
	}

	res, err := ConvertResponse[response.ResShinyouTategyokuList](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (b *balanceClientImpl) GetZanKaiKanougaku(ctx context.Context, session *Session, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMZanKaiKanougaku"
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
		return nil, errors.Wrap(err, "GetZanKaiKanougaku failed")
	}

	res, err := ConvertResponse[response.ResZanKaiKanougaku](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (b *balanceClientImpl) GetZanKaiKanougakuSuii(ctx context.Context, session *Session, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMZanKaiKanougakuSuii"
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
		return nil, errors.Wrap(err, "GetZanKaiKanougakuSuii failed")
	}

	res, err := ConvertResponse[response.ResZanKaiKanougakuSuii](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiSummary(ctx context.Context, session *Session) (*response.ResZanKaiSummary, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req := request.ReqZanKaiSummary{}
	req.CLMID = "CLMZanKaiSummary"
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
		return nil, errors.Wrap(err, "GetZanKaiSummary failed")
	}

	res, err := ConvertResponse[response.ResZanKaiSummary](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, session *Session, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req := request.ReqZanKaiGenbutuKaitukeSyousai{}
	req.CLMID = "CLMZanKaiGenbutuKaitukeSyousai"
	req.HitukeIndex = strconv.Itoa(tradingDay)
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
		return nil, errors.Wrap(err, "GetZanKaiGenbutuKaitukeSyousai failed")
	}

	res, err := ConvertResponse[response.ResZanKaiGenbutuKaitukeSyousai](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, session *Session, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req := request.ReqZanKaiSinyouSinkidateSyousai{}
	req.CLMID = "CLMZanKaiSinyouSinkidateSyousai"
	req.HitukeIndex = strconv.Itoa(tradingDay)
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
		return nil, errors.Wrap(err, "GetZanKaiSinyouSinkidateSyousai failed")
	}

	res, err := ConvertResponse[response.ResZanKaiSinyouSinkidateSyousai](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (b *balanceClientImpl) GetZanRealHosyoukinRitu(ctx context.Context, session *Session, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMZanRealHosyoukinRitu"
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
		return nil, errors.Wrap(err, "GetZanRealHosyoukinRitu failed")
	}

	res, err := ConvertResponse[response.ResZanRealHosyoukinRitu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanShinkiKanoIjiritu(ctx context.Context, session *Session, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMZanShinkiKanoIjiritu"
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
		return nil, errors.Wrap(err, "GetZanShinkiKanoIjiritu failed")
	}

	res, err := ConvertResponse[response.ResZanShinkiKanoIjiritu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanUriKanousuu(ctx context.Context, session *Session, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.RequestURL) // sessionからURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	req.CLMID = "CLMZanUriKanousuu"
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
		return nil, errors.Wrap(err, "GetZanUriKanousuu failed")
	}

	res, err := ConvertResponse[response.ResZanUriKanousuu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
