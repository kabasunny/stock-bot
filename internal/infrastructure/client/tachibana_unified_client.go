package client

import (
	"context"
	"log/slog"
	balance_response "stock-bot/internal/infrastructure/client/dto/balance/response"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	order_response "stock-bot/internal/infrastructure/client/dto/order/response"
	price_request "stock-bot/internal/infrastructure/client/dto/price/request"
	price_response "stock-bot/internal/infrastructure/client/dto/price/response"
)

// TachibanaUnifiedClient は立花証券の3つのI/F（認証、REQUEST、EVENT）を統合したクライアント
type TachibanaUnifiedClient struct {
	authClient    AuthClient
	balanceClient BalanceClient
	orderClient   OrderClient
	priceClient   PriceInfoClient
	masterClient  MasterDataClient
	eventClient   EventClient

	sessionManager SessionManager
	logger         *slog.Logger
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
	// セッションディレクトリのパス
	sessionDir := "./data/sessions"

	// 日付ベースセッション管理を初期化
	sessionManager := NewDateBasedSessionManager(
		authClient,
		userID,
		password,
		secondPassword,
		sessionDir,
		logger,
	)

	return &TachibanaUnifiedClient{
		authClient:     authClient,
		balanceClient:  balanceClient,
		orderClient:    orderClient,
		priceClient:    priceClient,
		masterClient:   masterClient,
		eventClient:    eventClient,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// EnsureAuthenticated はセッションが有効であることを確認し、必要に応じて再認証を行います
func (c *TachibanaUnifiedClient) EnsureAuthenticated(ctx context.Context) error {
	return c.sessionManager.EnsureAuthenticated(ctx)
}

// GetSession は現在のセッションを取得します（認証が必要な場合は自動で実行）
func (c *TachibanaUnifiedClient) GetSession(ctx context.Context) (*Session, error) {
	return c.sessionManager.GetSession(ctx)
}

// Logout はセッションを終了します
func (c *TachibanaUnifiedClient) Logout(ctx context.Context) error {
	return c.sessionManager.Logout(ctx)
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
	return c.sessionManager.IsAuthenticated()
}
