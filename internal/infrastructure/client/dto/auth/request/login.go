// infrastructure/client/dto/auth/request/login.go
package request

import "stock-bot/internal/infrastructure/client/dto"

type ReqLogin struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`    // 機能ID (固定値: "CLMAuthLoginRequest")
	UserId          string `json:"sUserId"`   // e支店口座のログインＩＤ
	Password        string `json:"sPassword"` // e支店口座のログインパスワード
}
