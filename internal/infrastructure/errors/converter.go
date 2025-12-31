package errors

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"
)

// ConvertError は標準的なGoエラーをAppErrorに変換
func ConvertError(err error) *AppError {
	if err == nil {
		return nil
	}

	// 既にAppErrorの場合はそのまま返す
	if appErr, ok := GetAppError(err); ok {
		return appErr
	}

	// コンテキストエラー
	if err == context.DeadlineExceeded {
		return NewTimeoutError(err)
	}
	if err == context.Canceled {
		return NewAppError(ErrCodeInternal, "Operation was canceled", err)
	}

	// データベースエラー
	if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
		return NewNotFoundError("record")
	}
	if strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "UNIQUE constraint") {
		return NewAlreadyExistsError("record")
	}

	// ネットワークエラー
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return NewTimeoutError(err)
		}
		return NewNetworkError(err)
	}

	// URL解析エラー
	if _, ok := err.(*url.Error); ok {
		return NewNetworkError(err)
	}

	// システムコールエラー
	if syscallErr, ok := err.(*syscall.Errno); ok {
		switch *syscallErr {
		case syscall.ECONNREFUSED:
			return NewNetworkError(err)
		case syscall.ETIMEDOUT:
			return NewTimeoutError(err)
		default:
			return NewInternalError(err)
		}
	}

	// その他のエラーは内部エラーとして扱う
	return NewInternalError(err)
}

// WrapError は既存のエラーを指定されたAppErrorでラップ
func WrapError(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	return NewAppError(code, message, err)
}

// WrapErrorWithDetails は詳細付きでエラーをラップ
func WrapErrorWithDetails(err error, code ErrorCode, message, details string) *AppError {
	if err == nil {
		return nil
	}
	return NewAppErrorWithDetails(code, message, details, err)
}

// ChainError は複数のエラーをチェーン
func ChainError(baseErr *AppError, cause error) *AppError {
	if baseErr == nil {
		return ConvertError(cause)
	}
	if cause == nil {
		return baseErr
	}

	// 既存のAppErrorに新しい原因を追加
	return &AppError{
		Code:       baseErr.Code,
		Message:    baseErr.Message,
		Details:    baseErr.Details,
		Cause:      cause,
		HTTPStatus: baseErr.HTTPStatus,
	}
}

// RetryableError は再試行可能なエラーかどうかを判定
func IsRetryable(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		switch appErr.Code {
		case ErrCodeNetwork, ErrCodeTimeout, ErrCodeRateLimit:
			return true
		case ErrCodeAPIError:
			// API エラーの場合は詳細を確認
			return strings.Contains(appErr.Message, "temporary") ||
				strings.Contains(appErr.Message, "timeout") ||
				strings.Contains(appErr.Message, "rate limit")
		}
	}

	// ネットワークエラーやタイムアウトは再試行可能
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	return false
}

// IsPermanent は永続的なエラーかどうかを判定
func IsPermanent(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		switch appErr.Code {
		case ErrCodeInvalidInput, ErrCodeNotFound, ErrCodeUnauthorized,
			ErrCodeAlreadyExists, ErrCodeInvalidOrder:
			return true
		}
	}
	return false
}

// ShouldLog はエラーをログに出力すべきかどうかを判定
func ShouldLog(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		switch appErr.Code {
		case ErrCodeInvalidInput, ErrCodeNotFound:
			// クライアントエラーは通常ログ出力不要
			return false
		}
	}
	return true
}

// GetRetryDelay は再試行までの待機時間を取得
func GetRetryDelay(err error, attempt int) time.Duration {
	baseDelay := time.Second

	if appErr, ok := GetAppError(err); ok {
		switch appErr.Code {
		case ErrCodeRateLimit:
			// レート制限の場合は長めの待機
			baseDelay = 5 * time.Second
		case ErrCodeTimeout:
			// タイムアウトの場合は短めの待機
			baseDelay = 500 * time.Millisecond
		}
	}

	// 指数バックオフ
	delay := baseDelay * time.Duration(1<<uint(attempt))

	// 最大60秒まで
	if delay > 60*time.Second {
		delay = 60 * time.Second
	}

	return delay
}
