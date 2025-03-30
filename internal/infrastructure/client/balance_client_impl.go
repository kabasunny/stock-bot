// internal/infrastructure/client/balance_client_impl.go
package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type balanceClientImpl struct {
	client *TachibanaClientImpl
	logger *zap.Logger // Loggerを追加
}

func NewBalanceClientImpl(client *TachibanaClientImpl, logger *zap.Logger) *balanceClientImpl {
	return &balanceClientImpl{
		client: client,
		logger: logger,
	}
}

func (b *balanceClientImpl) GetGenbutuKabuList(ctx context.Context) (*response.ResGenbutuKabuList, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req := request.ReqGenbutuKabuList{}
	req.CLMID = "CLMGenbutuKabuList"
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req) //utilの関数
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger) // b.logger を使用
	if err != nil {
		return nil, errors.Wrap(err, "GetGenbutuKabuList failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGenbutuKabuList](respMap) //utliの関数
	if err != nil {
		return nil, err
	}
	return res, nil
}

// internal/infrastructure/client/balance_client_impl.go

func (b *balanceClientImpl) GetShinyouTategyokuList(ctx context.Context) (*response.ResShinyouTategyokuList, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req := request.ReqShinyouTategyokuList{}
	req.CLMID = "CLMShinyouTategyokuList" // 修正: CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetShinyouTategyokuList failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResShinyouTategyokuList](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// internal/infrastructure/client/balance_client_impl.go

func (b *balanceClientImpl) GetZanKaiKanougaku(ctx context.Context, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMZanKaiKanougaku" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanKaiKanougaku failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanKaiKanougaku](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiKanougakuSuii(ctx context.Context, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMZanKaiKanougakuSuii" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanKaiKanougakuSuii failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanKaiKanougakuSuii](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiSummary(ctx context.Context) (*response.ResZanKaiSummary, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req := request.ReqZanKaiSummary{}
	req.CLMID = "CLMZanKaiSummary" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanKaiSummary failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanKaiSummary](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req := request.ReqZanKaiGenbutuKaitukeSyousai{}
	req.CLMID = "CLMZanKaiGenbutuKaitukeSyousai" // CLMID を設定
	req.HitukeIndex = strconv.Itoa(tradingDay)
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanKaiGenbutuKaitukeSyousai failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanKaiGenbutuKaitukeSyousai](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// internal/infrastructure/client/balance_client_impl.go

func (b *balanceClientImpl) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req := request.ReqZanKaiSinyouSinkidateSyousai{}
	req.CLMID = "CLMZanKaiSinyouSinkidateSyousai" // CLMID を設定
	req.HitukeIndex = strconv.Itoa(tradingDay)
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanKaiSinyouSinkidateSyousai failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanKaiSinyouSinkidateSyousai](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// internal/infrastructure/client/balance_client_impl.go

func (b *balanceClientImpl) GetZanRealHosyoukinRitu(ctx context.Context, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMZanRealHosyoukinRitu" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanRealHosyoukinRitu failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResZanRealHosyoukinRitu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanShinkiKanoIjiritu(ctx context.Context, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}
	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMZanShinkiKanoIjiritu" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanShinkiKanoIjiritu failed")
	}

	res, err := ConvertResponse[response.ResZanShinkiKanoIjiritu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *balanceClientImpl) GetZanUriKanousuu(ctx context.Context, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	if !b.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(b.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMZanUriKanousuu" // CLMID を設定
	req.P_no = b.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	respMap, err := SendRequest(httpReq, 3, b.logger)
	if err != nil {
		return nil, errors.Wrap(err, "GetZanUriKanousuu failed")
	}

	res, err := ConvertResponse[response.ResZanUriKanousuu](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}
