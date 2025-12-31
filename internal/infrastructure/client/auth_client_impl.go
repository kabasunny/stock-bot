// internal/infrastructure/client/auth_client_impl.go
package client

import (
	"context"
	"log/slog" // 追加
	"net/http"
	"net/http/cookiejar" // 追加
	"net/url"
	"strconv"
	"time"

	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/auth/response"

	"github.com/cockroachdb/errors"
)

type authClientImpl struct {
	client *TachibanaClientImpl
}

// LoginWithPost は、POSTメソッドを使用してログインを行い、Sessionを返す
func (a *authClientImpl) LoginWithPost(ctx context.Context, req request.ReqLogin) (*Session, error) {
	// 1. リクエストURLの作成
	authPath, _ := url.Parse("auth/")
	fullURL := a.client.baseURL.ResolveReference(authPath).String()

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMAuthLoginRequest"
	req.P_no = "1" // Login時は初期値"1"
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	// SendPostRequest を使用してリクエストを送信
	respMap, err := SendPostRequest(ctx, a.client.httpClient, fullURL, req, 3) // SendPostRequest は変更しない
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	// デバッグのため、レスポンスマップ全体をログに出力
	slog.InfoContext(ctx, "login response received", slog.Any("response_map", respMap))

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResLogin](respMap)
	if err != nil {
		// 失敗した場合、元のマップ情報を含めてエラーを返す
		return nil, errors.Wrapf(err, "failed to convert login response map: %+v", respMap)
	}

	// 6. ログイン成功/失敗の判定
	if res.ResultCode == "0" {

		// ログイン成功時の処理

		session := NewSession()

		session.SetLoginResponse(res) // ResLoginからSessionに情報をコピー

		session.SecondPassword = a.client.sSecondPassword // TachibanaClientImplからSecondPasswordをコピー

		// p_noの初期値をAPIレスポンスから設定
		if pNoStr, ok := respMap["p_no"].(string); ok {
			if pNo, err := strconv.ParseInt(pNoStr, 10, 32); err == nil {
				session.pNo.Store(int32(pNo))
			}
		}

		// *** CookieJar のコピー処理 ***
		// クライアントのhttpClient.JarからCookieをコピーし、新しいCookieJarを作成してSessionに設定
		newCookieJar, _ := cookiejar.New(nil)
		if clientJar, ok := a.client.httpClient.Jar.(*cookiejar.Jar); ok {
			// 現在のクライアントのCookieJarから全てのURLのCookieを取得
			allCookies := clientJar.Cookies(a.client.baseURL) // 例: ログインURLのCookieを取得
			newCookieJar.SetCookies(a.client.baseURL, allCookies) // 新しいJarに設定
		}
		session.CookieJar = newCookieJar // 独立したCookieJarをSessionに設定

		return session, nil
	}

	// ログイン失敗時の処理
	return nil, errors.Errorf("login failed with result code '%s': %s. raw response: %+v", res.ResultCode, res.ResultText, respMap)
}

func (a *authClientImpl) LogoutWithPost(ctx context.Context, session *Session, req request.ReqLogout) (*response.ResLogout, error) {
	// 1. リクエストURLの作成
	if session == nil {
		return nil, errors.New("session is nil")
	}
	u, err := url.Parse(session.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL from session")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMAuthLogoutRequest"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // Sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	// 3. HTTPリクエストの作成 (クライアントのhttpClientを使用し、セッションのCookieJarを設定)
	// クライアントのhttpClientは通常、自身のJarを持たないため、
	// セッション固有のCookieJarを持つ新しいhttpClientを作成するか、
	// リクエストごとにSessionのCookieをヘッダーに設定する必要がある。
	// 今回のケースでは、a.client.httpClientを使用し、リクエストにSessionのCookieJarを設定する方法を取る。

	// 一時的な http.Client を作成 (セッション固有のCookieJarを使用)
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	// 3. SendPostRequest を使用してリクエストを送信
	respMap, err := SendPostRequest(ctx, tempClient, u.String(), req, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "logout failed")
	}

	// デバッグログを追加
	slog.Info("Logout API response map", slog.Any("response", respMap))

	// 4. レスポンスの処理
	res, err := ConvertResponse[response.ResLogout](respMap)
	if err != nil {
		return nil, err
	}

	// 5. クライアントのログイン状態をリセットするコードは削除 (クライアントは状態を持たないため)

	return res, nil
}
