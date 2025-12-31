package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode はエラーコードの型
type ErrorCode string

const (
	// システムエラー
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase     ErrorCode = "DATABASE_ERROR"
	ErrCodeNetwork      ErrorCode = "NETWORK_ERROR"
	ErrCodeTimeout      ErrorCode = "TIMEOUT_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// ビジネスロジックエラー
	ErrCodeInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists     ErrorCode = "ALREADY_EXISTS"
	ErrCodeInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodeInvalidOrder      ErrorCode = "INVALID_ORDER"
	ErrCodeOrderNotFound     ErrorCode = "ORDER_NOT_FOUND"
	ErrCodePositionNotFound  ErrorCode = "POSITION_NOT_FOUND"

	// 外部API関連エラー
	ErrCodeAPIError       ErrorCode = "API_ERROR"
	ErrCodeSessionExpired ErrorCode = "SESSION_EXPIRED"
	ErrCodeRateLimit      ErrorCode = "RATE_LIMIT"
	ErrCodeMarketClosed   ErrorCode = "MARKET_CLOSED"
)

// AppError はアプリケーション固有のエラー
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	Cause      error     `json:"-"`
	HTTPStatus int       `json:"-"`
}

// Error は error インターフェースを実装
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap に対応
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError は新しいAppErrorを作成
func NewAppError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Cause:      cause,
		HTTPStatus: getHTTPStatusFromCode(code),
	}
}

// NewAppErrorWithDetails は詳細付きのAppErrorを作成
func NewAppErrorWithDetails(code ErrorCode, message, details string, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		Cause:      cause,
		HTTPStatus: getHTTPStatusFromCode(code),
	}
}

// getHTTPStatusFromCode はエラーコードからHTTPステータスコードを取得
func getHTTPStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidInput, ErrCodeInvalidOrder:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeSessionExpired:
		return http.StatusUnauthorized
	case ErrCodeNotFound, ErrCodeOrderNotFound, ErrCodePositionNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeInsufficientFunds, ErrCodeMarketClosed:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// 便利な関数群

// IsAppError は error が AppError かどうかを判定
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError は error から AppError を取得
func GetAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// HasCode は error が指定されたコードを持つかどうかを判定
func HasCode(err error, code ErrorCode) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Code == code
	}
	return false
}

// 事前定義されたエラー

// システムエラー
func NewInternalError(cause error) *AppError {
	return NewAppError(ErrCodeInternal, "Internal server error", cause)
}

func NewDatabaseError(cause error) *AppError {
	return NewAppError(ErrCodeDatabase, "Database operation failed", cause)
}

func NewNetworkError(cause error) *AppError {
	return NewAppError(ErrCodeNetwork, "Network operation failed", cause)
}

func NewTimeoutError(cause error) *AppError {
	return NewAppError(ErrCodeTimeout, "Operation timed out", cause)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, nil)
}

// ビジネスロジックエラー
func NewInvalidInputError(details string) *AppError {
	return NewAppErrorWithDetails(ErrCodeInvalidInput, "Invalid input", details, nil)
}

func NewNotFoundError(resource string) *AppError {
	return NewAppErrorWithDetails(ErrCodeNotFound, "Resource not found", resource, nil)
}

func NewAlreadyExistsError(resource string) *AppError {
	return NewAppErrorWithDetails(ErrCodeAlreadyExists, "Resource already exists", resource, nil)
}

func NewInsufficientFundsError(required, available float64) *AppError {
	details := fmt.Sprintf("Required: %.2f, Available: %.2f", required, available)
	return NewAppErrorWithDetails(ErrCodeInsufficientFunds, "Insufficient funds", details, nil)
}

func NewInvalidOrderError(details string) *AppError {
	return NewAppErrorWithDetails(ErrCodeInvalidOrder, "Invalid order", details, nil)
}

func NewOrderNotFoundError(orderID string) *AppError {
	return NewAppErrorWithDetails(ErrCodeOrderNotFound, "Order not found", orderID, nil)
}

func NewPositionNotFoundError(symbol string) *AppError {
	return NewAppErrorWithDetails(ErrCodePositionNotFound, "Position not found", symbol, nil)
}

// 外部API関連エラー
func NewAPIError(message string, cause error) *AppError {
	return NewAppError(ErrCodeAPIError, message, cause)
}

func NewSessionExpiredError() *AppError {
	return NewAppError(ErrCodeSessionExpired, "Session has expired", nil)
}

func NewRateLimitError() *AppError {
	return NewAppError(ErrCodeRateLimit, "Rate limit exceeded", nil)
}

func NewMarketClosedError() *AppError {
	return NewAppError(ErrCodeMarketClosed, "Market is closed", nil)
}
