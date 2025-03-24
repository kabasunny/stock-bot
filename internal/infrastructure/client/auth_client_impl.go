// internal/infrastructure/client/auth_client_impl.go
package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/auth/response"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type authClientImpl struct {
	client *TachibanaClient
	logger *zap.Logger
}

func (a *authClientImpl) Login(ctx context.Context, req request.ReqLogin) (*response.ResLogin, error) {
	// 1. リクエストURLの作成
	u, err := url.Parse(a.client.baseURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse base URL")
	}
	u.Path += "auth/"

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMAuthLoginRequest"
	req.P_no = "1" // Login時は初期値"1"
	req.P_sd_date = formatSDDate(time.Now())
	req.SJsonOfmt = "4"

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
	respMap, err := SendRequest(httpReq, 3, a.logger)
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResLogin](respMap) //utliの関数
	if err != nil {
		return nil, err
	}

	// 6. ログイン成功時の処理
	if res.ResultCode == "0" {
		a.client.loggined = true

		// p_noの更新 (Login APIのレスポンスで返ってくる値で更新)
		if pNoStr, ok := respMap["p_no"].(string); ok {
			if pNo, err := strconv.ParseInt(pNoStr, 10, 64); err == nil {
				a.client.p_NoMu.Lock()
				a.client.p_no = pNo
				a.client.p_NoMu.Unlock()
			}
		}

		// LoginInfo を更新
		a.client.loginInfo = &LoginInfo{
			RequestURL: res.RequestURL,
			MasterURL:  res.MasterURL,
			PriceURL:   res.PriceURL,
			EventURL:   res.EventURL,
			Expiry:     time.Now().Add(24 * time.Hour),
		}
	}

	return res, nil
}

func (a *authClientImpl) Logout(ctx context.Context, req request.ReqLogout) (*response.ResLogout, error) {
	// 1. リクエストURLの作成
	if a.client.loginInfo == nil {
		return nil, errors.New("not logged in")
	}
	u, err := url.Parse(a.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMAuthLogoutRequest"
	req.P_no = a.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.SJsonOfmt = "4"

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

	// 3. HTTPリクエストの作成
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, a.logger)
	if err != nil {
		//ログアウト失敗時も、ログイン状態はfalseにする
		a.client.loggined = false
		a.client.loginInfo = nil
		return nil, errors.Wrap(err, "logout failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResLogout](respMap) //utliの関数
	if err != nil {
		//ログアウト失敗時も、ログイン状態はfalseにする
		a.client.loggined = false
		a.client.loginInfo = nil
		return nil, err
	}

	// 6. ログアウト成功/失敗に関わらず、状態をリセット
	a.client.loggined = false
	a.client.loginInfo = nil

	return res, nil
}
