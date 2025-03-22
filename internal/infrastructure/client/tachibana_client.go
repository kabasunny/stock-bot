// infrastructure/client/tachibana_client.go
package client

import (
	"context"
	"net/url"
	"sync"
	"time"

	"stock-bot/internal/infrastructure/client/dto/request"
	"stock-bot/internal/infrastructure/client/dto/response"
)

// TachibanaClient 構造体の定義
type TachibanaClient struct {
	baseURL         *url.URL // 本番、デモ
	sUserId         string   // e支店口座のログインＩＤ
	sPassword       string   // e支店口座のログインパスワード
	sSecondPassword string   // 第二暗証番号（発注パスワード)

	loginInfo *LoginInfo   // ログイン情報
	loggined  bool         // ログインフラグ
	mu        sync.RWMutex // 排他制御用
	p_no      int64        // p_no の連番管理用
	p_NoMu    sync.Mutex   // pNo の排他制御用

	targetIssueCodes []string
	// masterData       *domain.MasterData // リポジトリから取得するマスターデータの要否確認

}

// LoginInfo はログイン後に取得する情報を保持します
type LoginInfo struct {
	RequestURL string    // キャッシュする仮想URL（REQUEST)
	MasterURL  string    // キャッシュする仮想URL（Master)
	PriceURL   string    // キャッシュする仮想URL（Price)
	EventURL   string    // キャッシュする仮想URL（EVENT)
	Expiry     time.Time // 仮想URLの有効期限
}

// Client インターフェースのメソッドを実装
func (c *TachibanaClient) Login(ctx context.Context, userID, password string) (*LoginInfo, error) {
	// ...
	// 	return &loginInfo, nil
	return nil, nil
}

func (c *TachibanaClient) Logout(ctx context.Context) error {
	// ...
	return nil
}
func (c *TachibanaClient) NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	// ...
	return nil, nil
}

// 他のメソッドも同様に実装
