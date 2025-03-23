// internal/infrastructure/client/auth_client.go
package client

import (
	"context"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/auth/response"
)

// AuthClient は、認証関連の API (ログイン、ログアウト) を扱うインターフェース
type AuthClient interface {
	// Login は、ユーザーIDとパスワードで認証を行い、API利用に必要な情報を取得
	Login(ctx context.Context, req request.ReqLogin) (*response.ResLogin, error)
	// Logout は、ログインセッションを終了
	Logout(ctx context.Context, req request.ReqLogout) (*response.ResLogout, error)
}
