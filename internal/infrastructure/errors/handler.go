package errors

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse はHTTPエラーレスポンスの構造
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail はエラーの詳細情報
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// HTTPErrorHandler はHTTPエラーハンドラー
type HTTPErrorHandler struct {
	logger *slog.Logger
}

// NewHTTPErrorHandler は新しいHTTPErrorHandlerを作成
func NewHTTPErrorHandler(logger *slog.Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{
		logger: logger,
	}
}

// HandleError はエラーをHTTPレスポンスに変換
func (h *HTTPErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	var ok bool

	// AppErrorかどうかを判定
	if appErr, ok = GetAppError(err); !ok {
		// 通常のerrorの場合はInternalErrorに変換
		appErr = NewInternalError(err)
	}

	// ログ出力
	h.logError(r.Context(), appErr, r)

	// HTTPレスポンスを作成
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	// レスポンスヘッダーを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)

	// JSONレスポンスを送信
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode error response", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// logError はエラーをログに出力
func (h *HTTPErrorHandler) logError(ctx context.Context, appErr *AppError, r *http.Request) {
	logLevel := h.getLogLevel(appErr.Code)

	attrs := []slog.Attr{
		slog.String("error_code", string(appErr.Code)),
		slog.String("error_message", appErr.Message),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", appErr.HTTPStatus),
	}

	if appErr.Details != "" {
		attrs = append(attrs, slog.String("details", appErr.Details))
	}

	if appErr.Cause != nil {
		attrs = append(attrs, slog.String("cause", appErr.Cause.Error()))
	}

	// リクエストIDがあれば追加
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}

	switch logLevel {
	case slog.LevelError:
		h.logger.LogAttrs(ctx, slog.LevelError, "HTTP error occurred", attrs...)
	case slog.LevelWarn:
		h.logger.LogAttrs(ctx, slog.LevelWarn, "HTTP warning occurred", attrs...)
	case slog.LevelInfo:
		h.logger.LogAttrs(ctx, slog.LevelInfo, "HTTP info", attrs...)
	}
}

// getLogLevel はエラーコードに応じたログレベルを取得
func (h *HTTPErrorHandler) getLogLevel(code ErrorCode) slog.Level {
	switch code {
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeNetwork:
		return slog.LevelError
	case ErrCodeTimeout, ErrCodeAPIError:
		return slog.LevelWarn
	case ErrCodeInvalidInput, ErrCodeNotFound, ErrCodeUnauthorized:
		return slog.LevelInfo
	default:
		return slog.LevelWarn
	}
}

// Middleware はHTTPミドルウェアとしてエラーハンドリングを提供
func (h *HTTPErrorHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// パニックをキャッチしてエラーに変換
		defer func() {
			if recovered := recover(); recovered != nil {
				var err error
				if e, ok := recovered.(error); ok {
					err = e
				} else {
					err = NewInternalError(nil)
				}
				h.HandleError(w, r, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
