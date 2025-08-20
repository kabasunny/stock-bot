package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	// 開発中はテキスト形式の方が見やすいため TextHandler を使用します。
	// 本番環境では、機械処理しやすい JSONHandler の利用を推奨します。
	// opts := &slog.HandlerOptions{
	// 	Level: slog.LevelDebug, // ログレベルをDEBUGに設定
	// }
	// handler := slog.NewJSONHandler(os.Stdout, opts)
	logLevel := slog.LevelInfo // デフォルトはInfoレベル
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	log = slog.New(handler)

	slog.SetDefault(log)
}

// L returns the application logger.
// Deprecated: Use slog.Default() instead.
func L() *slog.Logger {
	return log
}
