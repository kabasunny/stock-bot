// internal/infrastructure/client/tachibana_client.go
package client

import (
	"net/url"
	"sync"
	"time"
)

// TachibanaClient は、橘証券 e支店 API クライアントの構造体です。
type TachibanaClient struct {
	baseURL          *url.URL     // APIのベースURL (本番/デモ環境)
	sUserId          string       // e支店口座のログインＩＤ
	sPassword        string       // e支店口座のログインパスワード
	sSecondPassword  string       // 第二暗証番号（発注パスワード)
	loginInfo        *LoginInfo   // ログイン後に取得した情報 (各種URL、有効期限)
	loggined         bool         // ログイン状態 (true: ログイン中, false: 未ログイン)
	mu               sync.RWMutex // 排他制御用ミューテックス (読み取り/書き込みロック)
	p_no             int64        // リクエストに付与する一意な番号 (連番)
	p_NoMu           sync.Mutex   // p_no の排他制御用ミューテックス
	targetIssueCodes []string     // 利用する銘柄コード

	// 埋め込む構造体 (各機能別のクライアント実装)
	*authClientImpl
	*orderClientImpl
	*balanceClientImpl
	*masterDataClientImpl
	*priceInfoClientImpl
}

// NewTachibanaClient は TachibanaClient のコンストラクタです。
// 必要な情報を引数で受け取り、TachibanaClient インスタンスを生成して返します。
func NewTachibanaClient(baseURL *url.URL, sUserId, sPassword, sSecondPassword string, targetIssueCodes []string) *TachibanaClient {
	client := &TachibanaClient{
		baseURL:          baseURL,
		sUserId:          sUserId,
		sPassword:        sPassword,
		sSecondPassword:  sSecondPassword,
		targetIssueCodes: targetIssueCodes,
	}
	// 埋め込む構造体の初期化 (各機能別のクライアント実装を関連付け)
	client.authClientImpl = &authClientImpl{client: client}
	client.orderClientImpl = &orderClientImpl{client: client}
	client.balanceClientImpl = &balanceClientImpl{client: client}
	client.masterDataClientImpl = &masterDataClientImpl{client: client}
	client.priceInfoClientImpl = &priceInfoClientImpl{client: client}

	return client
}

// LoginInfo は、ログイン後に取得する情報を保持する構造体
type LoginInfo struct {
	RequestURL string    // リクエスト用URL (業務機能にアクセスするためのURL)
	MasterURL  string    // マスタ用URL (マスタ情報にアクセスするためのURL)
	PriceURL   string    // 時価情報用URL (時価情報にアクセスするためのURL)
	EventURL   string    // イベント用URL (注文約定通知などを受信するためのURL)
	Expiry     time.Time // 各URLの有効期限
}
