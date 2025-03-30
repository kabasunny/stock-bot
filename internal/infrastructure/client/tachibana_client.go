// internal/infrastructure/client/tachibana_client.go
package client

import (
	"net/url"
	"stock-bot/internal/config"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TachibanaClientImpl, NewTachibanaClient, LoginInfo, getPNo は変更なし (省略)
// TachibanaClientImpl は、橘証券 e支店 API クライアントの構造体です。
type TachibanaClientImpl struct {
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
	logger           *zap.Logger  // ロガー

	// 埋め込みフィールド
	// 型名がフィールド名となり、同じ型の複数の埋め込みは不可の模様
	*authClientImpl
	*orderClientImpl
	*balanceClientImpl
	*masterDataClientImpl
	*priceInfoClientImpl
}

// NewTachibanaClient は TachibanaClient のコンストラクタです。
// 必要な情報を引数で受け取り、TachibanaClient インスタンスを生成して返します。
func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) *TachibanaClientImpl {
	baseURL, _ := url.Parse(cfg.TachibanaBaseURL) // 文字列から *url.URL に変換
	client := &TachibanaClientImpl{
		baseURL:          baseURL, // *url.URL型
		sUserId:          cfg.TachibanaUserID,
		sPassword:        cfg.TachibanaPassword,
		sSecondPassword:  cfg.TachibanaPassword,
		targetIssueCodes: []string{},     // 必要に応じて設定
		loginInfo:        nil,            // 初期値はnil
		loggined:         false,          // 初期値はfalse
		mu:               sync.RWMutex{}, // 初期化
		p_no:             0,              // 初期値は0
		p_NoMu:           sync.Mutex{},   // 初期化
		logger:           logger,         // ロガー
	}
	// 埋め込む構造体の初期化 (各機能別のクライアント実装を関連付け)
	client.authClientImpl = &authClientImpl{client: client, logger: logger}
	client.orderClientImpl = &orderClientImpl{client: client, logger: logger}
	client.balanceClientImpl = &balanceClientImpl{client: client, logger: logger}
	client.masterDataClientImpl = &masterDataClientImpl{client: client, logger: logger}
	client.priceInfoClientImpl = &priceInfoClientImpl{client: client, logger: logger}

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

// getPNo は p_no を取得し、インクリメントする (スレッドセーフ)
func (tc *TachibanaClientImpl) getPNo() string {
	tc.p_NoMu.Lock()
	defer tc.p_NoMu.Unlock()
	tc.p_no++
	return strconv.FormatInt(tc.p_no, 10)
}
