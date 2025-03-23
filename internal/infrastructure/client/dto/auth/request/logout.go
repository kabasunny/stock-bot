// infrastructure/client/dto/auth/request/login.go
package request

import "stock-bot/internal/infrastructure/client/dto"

type ReqLogout struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"` // 機能ID (固定値: "CLMAuthLogoutRequest")
}
