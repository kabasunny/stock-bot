package client

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	balance_response "stock-bot/internal/infrastructure/client/dto/balance/response"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	order_response "stock-bot/internal/infrastructure/client/dto/order/response"
	price_request "stock-bot/internal/infrastructure/client/dto/price/request"
	price_response "stock-bot/internal/infrastructure/client/dto/price/response"
	"sync"
	"time"
)

// TachibanaUnifiedClient は立花証券の3つのI/F（認証、REQUEST、EVENT）を統合したクライアント
type TachibanaUnifiedClient struct {
	authClient    AuthClient
	balanceClient BalanceClient
	orderClient   OrderClient
	priceClient   PriceInfoClient
	masterClient  MasterDataClient
	eventClient   EventClient

	session      *Session
	sessionMutex sync.RWMutex
	logger       *slog.Logger

	// 認証情報
	userID         string
	password       string
	secondPassword string

	// セッション管理
	lastLoginTime time.Time
	loginMutex    sync.Mutex
}

// NewTachibanaUnifiedClient は新しいTachibanaUnifiedClientを作成します
func NewTachibanaUnifiedClient(
	authClient AuthClient,
	balanceClient BalanceClient,
	orderClient OrderClient,
	priceClient PriceInfoClient,
	masterClient MasterDataClient,
	eventClient EventClient,
	userID, password, secondPassword string,
	logger *slog.Logger,
) *TachibanaUnifiedClient {
	return &TachibanaUnifiedClient{
		authClient:     authClient,
		balanceClient:  balanceClient,
		orderClient:    orderClient,
		priceClient:    priceClient,
		masterClient:   masterClient,
		eventClient:    eventClient,
		userID:         userID,
		password:       password,
		secondPassword: secondPassword,
		logger:         logger,
	}
}

// EnsureAuthenticated はセッションが有効であることを確認し、必要に応じて再認証を行います
func (c *TachibanaUnifiedClient) EnsureAuthenticated(ctx context.Context) error {
	c.loginMutex.Lock()
	defer c.loginMutex.Unlock()

	// セッションが存在し、まだ有効な場合はそのまま使用
	if c.session != nil && time.Since(c.lastLoginTime) < 8*time.Hour {
		return nil
	}

	c.logger.Info("performing authentication to Tachibana API")

	// ログインリクエストを作成
	loginReq := request.ReqLogin{
		UserId:   c.userID,
		Password: c.password,
	}

	// 認証実行
	session, err := c.authClient.LoginWithPost(ctx, loginReq)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// セッション情報を設定
	session.SecondPassword = c.secondPassword

	c.sessionMutex.Lock()
	c.session = session
	c.lastLoginTime = time.Now()
	c.sessionMutex.Unlock()

	c.logger.Info("authentication successful")
	return nil
}

// GetSession は現在のセッションを取得します（認証が必要な場合は自動で実行）
func (c *TachibanaUnifiedClient) GetSession(ctx context.Context) (*Session, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.sessionMutex.RLock()
	defer c.sessionMutex.RUnlock()
	return c.session, nil
}

// Logout はセッションを終了します
func (c *TachibanaUnifiedClient) Logout(ctx context.Context) error {
	c.sessionMutex.RLock()
	session := c.session
	c.sessionMutex.RUnlock()

	if session == nil {
		return nil // 既にログアウト済み
	}

	// ログアウトリクエストを作成
	logoutReq := request.ReqLogout{}

	// ログアウト実行
	_, err := c.authClient.LogoutWithPost(ctx, session, logoutReq)
	if err != nil {
		c.logger.Warn("logout request failed", "error", err)
		// ログアウトエラーでもセッションはクリアする
	}

	c.sessionMutex.Lock()
	c.session = nil
	c.sessionMutex.Unlock()

	c.logger.Info("logout completed")
	return err
}

// --- BalanceClient methods ---

func (c *TachibanaUnifiedClient) GetZanKaiSummary(ctx context.Context) (*balance_response.ResZanKaiSummary, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.balanceClient.GetZanKaiSummary(ctx, session)
}

func (c *TachibanaUnifiedClient) GetGenbutuKabuList(ctx context.Context) (*balance_response.ResGenbutuKabuList, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.balanceClient.GetGenbutuKabuList(ctx, session)
}

func (c *TachibanaUnifiedClient) GetShinyouTategyokuList(ctx context.Context) (*balance_response.ResShinyouTategyokuList, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.balanceClient.GetShinyouTategyokuList(ctx, session)
}

// --- OrderClient methods ---

func (c *TachibanaUnifiedClient) NewOrder(ctx context.Context, params NewOrderParams) (*order_response.ResNewOrder, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.orderClient.NewOrder(ctx, session, params)
}

func (c *TachibanaUnifiedClient) CorrectOrder(ctx context.Context, params CorrectOrderParams) (*order_response.ResCorrectOrder, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.orderClient.CorrectOrder(ctx, session, params)
}

func (c *TachibanaUnifiedClient) CancelOrder(ctx context.Context, params CancelOrderParams) (*order_response.ResCancelOrder, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.orderClient.CancelOrder(ctx, session, params)
}

func (c *TachibanaUnifiedClient) GetOrderList(ctx context.Context, req order_request.ReqOrderList) (*order_response.ResOrderList, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.orderClient.GetOrderList(ctx, session, req)
}

// --- PriceInfoClient methods ---

func (c *TachibanaUnifiedClient) GetPriceInfo(ctx context.Context, req price_request.ReqGetPriceInfo) (*price_response.ResGetPriceInfo, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.priceClient.GetPriceInfo(ctx, session, req)
}

func (c *TachibanaUnifiedClient) GetPriceInfoHistory(ctx context.Context, req price_request.ReqGetPriceInfoHistory) (*price_response.ResGetPriceInfoHistory, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return c.priceClient.GetPriceInfoHistory(ctx, session, req)
}

// --- EventClient methods ---

func (c *TachibanaUnifiedClient) ConnectEvents(ctx context.Context, symbols []string) (<-chan []byte, <-chan error, error) {
	session, err := c.GetSession(ctx)
	if err != nil {
		return nil, nil, err
	}
	return c.eventClient.Connect(ctx, session, symbols)
}

func (c *TachibanaUnifiedClient) CloseEvents() {
	c.eventClient.Close()
}

// IsAuthenticated はセッションが有効かどうかを確認します
func (c *TachibanaUnifiedClient) IsAuthenticated() bool {
	c.sessionMutex.RLock()
	defer c.sessionMutex.RUnlock()
	return c.session != nil && time.Since(c.lastLoginTime) < 8*time.Hour
}
