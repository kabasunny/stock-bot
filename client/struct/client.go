// tachibana/strct/client.go
package tachibana

import (
	"net/url"
	"sync"
	"time"
)

// TachibanaClientImple 構造体の定義
type TachibanaClientImple struct {
	baseURL         *url.URL // 本番、デモ
	sUserId         string   // e支店口座のログインＩＤ
	sPassword       string   // e支店口座のログインパスワード
	sSecondPassword string   // 第二暗証番号（発注パスワード)

	loggined   bool         // ログインフラグ
	requestURL string       // キャッシュする仮想URL（REQUEST)
	masterURL  string       // キャッシュする仮想URL（Master)
	priceURL   string       // キャッシュする仮想URL（Price)
	eventURL   string       // キャッシュする仮想URL（EVENT)
	expiry     time.Time    // 仮想URLの有効期限
	mu         sync.RWMutex // 排他制御用
	p_no       int64        // p_no の連番管理用
	p_NoMu     sync.Mutex   // pNo の排他制御用

	targetIssueCodes []string
	// masterData       *domain.MasterData // リポジトリから取得するマスターデータの要否確認

}
