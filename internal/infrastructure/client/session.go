package client

import (

	"net/http"

	"net/http/cookiejar"

	"sync/atomic" // p_noをアトミックに扱うため



	"stock-bot/internal/infrastructure/client/dto/auth/response" // ResLoginのためにインポート

)



// Session はAPIセッション情報を保持します。

// 各ログインによって生成され、そのセッションに紐づくAPIリクエストで使用されます。

type Session struct {

	ResultCode string

	ResultText string



	// 認証情報 (Login時にクライアントからコピー)

	SecondPassword string



	// 各種URL (ResLoginから取得)

	RequestURL string

	MasterURL  string

	PriceURL   string

	EventURL   string



	// セッション管理情報

	CookieJar  http.CookieJar // セッションCookieを保持



	// P_no (リクエスト番号) の管理

	pNo atomic.Int32

}

// NewSession は新しいSessionインスタンスを生成します。

func NewSession() *Session {

	s := &Session{}

	s.pNo.Store(0) // 初期値は0

	jar, _ := cookiejar.New(nil)

	s.CookieJar = jar

	return s

}

// GetPNo は現在のp_noを取得し、次のリクエストのためにインクリメントします。
func (s *Session) GetPNo() int32 {
	return s.pNo.Add(1)
}

// SetLoginResponse は ResLogin の情報で Session を初期化します。
func (s *Session) SetLoginResponse(res *response.ResLogin) {
	s.ResultCode = res.ResultCode
	s.ResultText = res.ResultText
	// ResLoginのフィールド名に合わせて修正
	s.RequestURL = res.RequestURL // sUrlRequestにマッピングされている
	s.MasterURL = res.MasterURL   // sUrlMasterにマッピングされている
	s.PriceURL = res.PriceURL     // sUrlPriceにマッピングされている
	s.EventURL = res.SUrlEventWebSocket // ここを WebSocket用のURLに修正
}
