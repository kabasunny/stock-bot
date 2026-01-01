package client

import (
	"os"
	"path/filepath"
)

// EnsureSessionDirectory はセッションディレクトリが存在することを確認し、必要に応じて作成する
func EnsureSessionDirectory(sessionDir string) error {
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return err
	}
	return nil
}

// GetDefaultSessionDirectory はデフォルトのセッションディレクトリパスを取得する
func GetDefaultSessionDirectory() string {
	return filepath.Join(".", "data", "sessions")
}

// CleanupOldSessions は古いセッションファイルをクリーンアップする（オプション機能）
func CleanupOldSessions(sessionDir string, keepDays int) error {
	// 実装は後で追加（現在は何もしない）
	return nil
}
