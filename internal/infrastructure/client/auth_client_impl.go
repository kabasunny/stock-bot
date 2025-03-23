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
	u.Path += "auth/" // ★★★ auth/ を追加 (末尾のスラッシュは不要) ★★★

	// 2. リクエストパラメータの作成 (JSON + URLエンコード)
	payload := map[string]string{ // ★★★ map[string]string を使用 ★★★
		"sCLMID":    "CLMAuthLoginRequest",
		"sUserId":   req.UserId,
		"sPassword": req.Password,
		"sJsonOfmt": "4",
		"p_no":      "1", // Login時は初期値"1"
		"p_sd_date": formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload) // ★★★ JSON文字列に変換 ★★★
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON)) // ★★★ URLエンコード ★★★
	u.RawQuery = encodedPayload                            // クエリパラメータとして設定

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json") // Content-Type は GET では不要

	// 4. リクエストの送信 (sendRequest を使用)
	respMap, err := SendRequest(httpReq, 3, a.client, a.logger)
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	// 5. レスポンスの処理
	res, err := a.convertToResLogin(respMap)
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
	u, err := url.Parse(a.client.loginInfo.RequestURL) // RequestURLをそのまま使う
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}
	// u.Path = "/request"  ★★★ この行は不要 ★★★

	// 2. リクエストパラメータの作成
	payload := map[string]string{ // ★★★ map[string]string を使用 ★★★
		"sCLMID":    "CLMAuthLogoutRequest",
		"sJsonOfmt": "4",
		"p_no":      a.client.getPNo(), // getPNo() でインクリメント & 値取得
		"p_sd_date": formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload) // ★★★ JSON文字列に変換 ★★★
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON)) // ★★★ URLエンコード ★★★
	u.RawQuery = encodedPayload                            // クエリパラメータとして設定

	// 3. HTTPリクエストの作成
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	httpReq.Header.Set("Content-Type", "application/json") //念のため

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, a.client, a.logger) // リトライ回数3
	if err != nil {
		//ログアウト失敗時も、ログイン状態はfalseにする
		a.client.loggined = false
		a.client.loginInfo = nil
		return nil, errors.Wrap(err, "logout failed")
	}

	// 5. レスポンスの処理
	res, err := a.convertToResLogout(respMap)
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

// convertToResLogin は、map[string]interface{} を response.ResLogin に変換する
func (a *authClientImpl) convertToResLogin(respMap map[string]interface{}) (*response.ResLogin, error) {
	// レスポンスがエラーの場合
	if resultCode, ok := respMap["sResultCode"].(string); ok && resultCode != "0" {
		resultText := ""
		if rt, ok := respMap["sResultText"].(string); ok {
			resultText = rt
		}
		return nil, errors.Errorf("API error: ResultCode=%s, ResultText=%s", resultCode, resultText)
	}

	// レスポンスが正常な場合
	resBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response map to json")
	}
	var res response.ResLogin
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response json")
	}
	return &res, nil
}

// convertToResLogout は、map[string]interface{} を response.ResLogout に変換する
func (a *authClientImpl) convertToResLogout(respMap map[string]interface{}) (*response.ResLogout, error) {
	// レスポンスがエラーの場合の処理
	if resultCode, ok := respMap["sResultCode"].(string); ok && resultCode != "0" {
		resultText := ""
		if rt, ok := respMap["sResultText"].(string); ok {
			resultText = rt
		}
		return nil, errors.Errorf("API error: ResultCode=%s, ResultText=%s", resultCode, resultText)
	}

	// レスポンスが正常な場合
	resBytes, err := json.Marshal(respMap) //mapをjsonに変換
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response map to json")
	}
	var res response.ResLogout
	if err := json.Unmarshal(resBytes, &res); err != nil { //jsonを構造体に変換
		return nil, errors.Wrap(err, "failed to unmarshal response json")
	}
	return &res, nil
}
