// internal/infrastructure/client/client.go
package client

// Client は、すべてのクライアント機能をまとめたインターフェース
type Client interface {
	AuthClient       // 認証関連の API (ログイン、ログアウト) を扱うインターフェース
	OrderClient      // 注文関連の API を扱うインターフェース
	BalanceClient    // 残高・余力関連の API を扱うインターフェース
	MasterDataClient // マスタデータ関連の API を扱うインターフェース
	PriceInfoClient  // 時価情報関連の API を扱うインターフェース
}
