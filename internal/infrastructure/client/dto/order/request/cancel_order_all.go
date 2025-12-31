// request/cancel_order_all.go
package request

import "stock-bot/internal/infrastructure/client/dto"

// ReqCancelOrderAll は株式一括取消のリクエストを表すDTO
type ReqCancelOrderAll struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`          // 機能ID, CLMKabuCancelOrderAll
	SecondPassword  string `json:"sSecondPassword"` // 第二パスワード
}
