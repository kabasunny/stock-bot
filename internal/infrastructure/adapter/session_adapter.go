package adapter

import (
	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// SessionAdapter はインフラ層のSessionをドメイン層のSessionに変換するアダプター
type SessionAdapter struct{}

// NewSessionAdapter は新しいSessionAdapterを作成する
func NewSessionAdapter() *SessionAdapter {
	return &SessionAdapter{}
}

// ToDomainSession はclient.Sessionをmodel.Sessionに変換する
func (a *SessionAdapter) ToDomainSession(clientSession *client.Session) *model.Session {
	if clientSession == nil {
		return nil
	}

	// セッション有効期限を8時間後に設定（立花証券の仕様）
	expiresAt := time.Now().Add(8 * time.Hour)

	return &model.Session{
		SessionID:    generateSessionID(clientSession),
		UserID:       extractUserID(clientSession),
		LoginTime:    time.Now(), // 実際のログイン時刻は別途管理が必要
		ExpiresAt:    expiresAt,
		IsActive:     clientSession.ResultCode == "0",
		ResultCode:   clientSession.ResultCode,
		ResultText:   clientSession.ResultText,
		LastActivity: time.Now(),
	}
}

// ToClientSession はmodel.Sessionをclient.Sessionに変換する（必要に応じて）
func (a *SessionAdapter) ToClientSession(domainSession *model.Session) *client.Session {
	if domainSession == nil {
		return nil
	}

	return &client.Session{
		ResultCode: domainSession.ResultCode,
		ResultText: domainSession.ResultText,
		// 他のフィールドは必要に応じて追加
	}
}

// generateSessionID はクライアントセッションからセッションIDを生成する
func generateSessionID(clientSession *client.Session) string {
	// 実際の実装では、クライアントセッションの情報を基にユニークなIDを生成
	// 現在は簡易実装
	if clientSession.ResultCode == "0" {
		return "session-" + time.Now().Format("20060102150405")
	}
	return ""
}

// extractUserID はクライアントセッションからユーザーIDを抽出する
func extractUserID(clientSession *client.Session) string {
	// 実際の実装では、クライアントセッションからユーザー情報を抽出
	// 現在は簡易実装
	if clientSession.ResultCode == "0" {
		return "user-placeholder"
	}
	return ""
}
