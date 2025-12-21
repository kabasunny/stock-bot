// internal/infrastructure/client/tachibana_client.go
package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"stock-bot/internal/config"
	_ "stock-bot/internal/logger"
	"sync"
)

// TachibanaClientImpl は、橘証券 e支店 API クライアントの構造体です。
type TachibanaClientImpl struct {
	baseURL          *url.URL
	httpClient       *http.Client // Add httpClient to manage cookies
	sUserId          string
	sPassword        string
	sSecondPassword  string // 追加
	mu               sync.RWMutex
	targetIssueCodes []string

	*authClientImpl
	*orderClientImpl
	*balanceClientImpl
	*masterDataClientImpl
	*priceInfoClientImpl
}

// NewTachibanaClient は TachibanaClient のコンストラクタです。
func NewTachibanaClient(cfg *config.Config) *TachibanaClientImpl {
	baseURL, _ := url.Parse(cfg.TachibanaBaseURL)
	jar, _ := cookiejar.New(nil) // Create a new cookie jar
	client := &TachibanaClientImpl{
		baseURL:   baseURL,
		httpClient: &http.Client{ // Initialize httpClient with the jar
			Jar: jar,
		},
		sUserId:          cfg.TachibanaUserID,
		sPassword:        cfg.TachibanaPassword,
		sSecondPassword:  cfg.TachibanaPassword, // 追加
		mu:               sync.RWMutex{},
		targetIssueCodes: []string{},
	}
	client.authClientImpl = &authClientImpl{client: client}
	client.orderClientImpl = &orderClientImpl{client: client}
	client.balanceClientImpl = &balanceClientImpl{client: client}
	client.masterDataClientImpl = &masterDataClientImpl{client: client}
	client.priceInfoClientImpl = &priceInfoClientImpl{client: client}

	return client
}

// CookieJar returns the http.CookieJar used by the client.
func (tc *TachibanaClientImpl) CookieJar() http.CookieJar {
	return tc.httpClient.Jar
}
