// internal/infrastructure/client/tachibana_client.go
package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"stock-bot/internal/config"
	_ "stock-bot/internal/logger"
	"strconv"
	"sync"
	"time"
)

// TachibanaClientImpl は、橘証券 e支店 API クライアントの構造体です。
type TachibanaClientImpl struct {
	baseURL          *url.URL
	httpClient       *http.Client // Add httpClient to manage cookies
	sUserId          string
	sPassword        string
	sSecondPassword  string
	loginInfo        *LoginInfo
	loggined         bool
	mu               sync.RWMutex
	p_no             int64
	p_NoMu           sync.Mutex
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
		sSecondPassword:  cfg.TachibanaPassword,
		targetIssueCodes: []string{},
		loginInfo:        nil,
		loggined:         false,
		mu:               sync.RWMutex{},
		p_no:             0,
		p_NoMu:           sync.Mutex{},
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


// LoginInfo は、ログイン後に取得する情報を保持する構造体
type LoginInfo struct {
	RequestURL string    // リクエスト用URL (業務機能にアクセスするためのURL)
	MasterURL  string    // マスタ用URL (マスタ情報にアクセスするためのURL)
	PriceURL   string    // 時価情報用URL (時価情報にアクセスするためのURL)
	EventURL   string    // イベント用URL (注文約定通知などを受信するためのURL)
	Expiry     time.Time // 各URLの有効期限
}

// getPNo は p_no を取得し、インクリメントする (スレッドセーフ)
func (tc *TachibanaClientImpl) getPNo() string {
	tc.p_NoMu.Lock()
	defer tc.p_NoMu.Unlock()
	tc.p_no++
	return strconv.FormatInt(tc.p_no, 10)
}
