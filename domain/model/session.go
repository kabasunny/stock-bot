package model

import "time"

// Session はドメイン層のセッション情報を表す
type Session struct {
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	LoginTime    time.Time `json:"login_time"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsActive     bool      `json:"is_active"`
	ResultCode   string    `json:"result_code"`
	ResultText   string    `json:"result_text"`
	LastActivity time.Time `json:"last_activity"`
}

// IsValid はセッションが有効かどうかを判定する
func (s *Session) IsValid() bool {
	return s.IsActive && s.ResultCode == "0" && time.Now().Before(s.ExpiresAt)
}

// IsExpired はセッションが期限切れかどうかを判定する
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UpdateActivity は最終アクティビティ時刻を更新する
func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now()
}
