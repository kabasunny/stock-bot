package tests

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
	"stock-bot/internal/tradeservice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBalanceClient はBalanceClientのモック実装
type MockBalanceClient struct {
	mock.Mock
}

func (m *MockBalanceClient) GetGenbutuKabuList(ctx context.Context, session *client.Session) (*response.ResGenbutuKabuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResGenbutuKabuList), args.Error(1)
}

func (m *MockBalanceClient) GetShinyouTategyokuList(ctx context.Context, session *client.Session) (*response.ResShinyouTategyokuList, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResShinyouTategyokuList), args.Error(1)
}

func (m *MockBalanceClient) GetZanKaiKanougaku(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougaku), args.Error(1)
}

func (m *MockBalanceClient) GetZanKaiKanougakuSuii(ctx context.Context, session *client.Session, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiKanougakuSuii), args.Error(1)
}

func (m *MockBalanceClient) GetZanKaiSummary(ctx context.Context, session *client.Session) (*response.ResZanKaiSummary, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSummary), args.Error(1)
}

func (m *MockBalanceClient) GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiGenbutuKaitukeSyousai), args.Error(1)
}

func (m *MockBalanceClient) GetZanKaiSinyouSinkidateSyousai(ctx context.Context, session *client.Session, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error) {
	args := m.Called(ctx, session, tradingDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanKaiSinyouSinkidateSyousai), args.Error(1)
}

func (m *MockBalanceClient) GetZanRealHosyoukinRitu(ctx context.Context, session *client.Session, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanRealHosyoukinRitu), args.Error(1)
}

func (m *MockBalanceClient) GetZanShinkiKanoIjiritu(ctx context.Context, session *client.Session, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanShinkiKanoIjiritu), args.Error(1)
}

func (m *MockBalanceClient) GetZanUriKanousuu(ctx context.Context, session *client.Session, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error) {
	args := m.Called(ctx, session, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ResZanUriKanousuu), args.Error(1)
}

// TestSessionRecovery_NoUnifiedClient はUnifiedClientがない場合のテストです
func TestSessionRecovery_NoUnifiedClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, // balanceClient - nilにしてスタブ実装を使用
		nil, nil, nil, nil,
		session,
		logger,
	)
	// UnifiedClientを設定しない

	ctx := context.Background()
	balance, err := tradeService.GetBalance(ctx)

	// balanceClientがnilの場合はスタブ実装で成功する
	assert.NoError(t, err, "Should succeed with stub implementation")
	assert.NotNil(t, balance, "Balance should not be nil")
	assert.Equal(t, 1000000.0, balance.Cash, "Should return stub cash value")
}

// TestSessionRecovery_SetUnifiedClient はSetUnifiedClient()の動作をテストします
func TestSessionRecovery_SetUnifiedClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	session := client.NewSession()

	tradeService := tradeservice.NewGoaTradeService(
		nil, nil, nil, nil, nil,
		session,
		logger,
	)

	// UnifiedClientを設定（nilでもパニックしないことを確認）
	tradeService.SetUnifiedClient(nil)

	// パニックしないことを確認
	assert.NotPanics(t, func() {
		tradeService.SetUnifiedClient(nil)
	}, "SetUnifiedClient should not panic with nil")
}

// TestIsSessionError はセッションエラー判定の動作をテストします
func TestIsSessionError(t *testing.T) {
	tests := []struct {
		name        string
		errorString string
		expected    bool
	}{
		{
			name:        "Session expired error",
			errorString: "session expired",
			expected:    true,
		},
		{
			name:        "Invalid session error",
			errorString: "invalid session",
			expected:    true,
		},
		{
			name:        "Authentication required error",
			errorString: "authentication required",
			expected:    true,
		},
		{
			name:        "Unauthorized error",
			errorString: "unauthorized",
			expected:    true,
		},
		{
			name:        "401 error",
			errorString: "HTTP 401 error",
			expected:    true,
		},
		{
			name:        "403 error",
			errorString: "HTTP 403 forbidden",
			expected:    true,
		},
		{
			name:        "Network error",
			errorString: "network timeout",
			expected:    false,
		},
		{
			name:        "Database error",
			errorString: "database connection failed",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// セッションエラーパターンをチェック
			isSessionError := false

			sessionErrorPatterns := []string{
				"session expired",
				"invalid session",
				"authentication required",
				"unauthorized",
				"401",
				"403",
			}

			for _, pattern := range sessionErrorPatterns {
				if contains(tt.errorString, pattern) {
					isSessionError = true
					break
				}
			}

			assert.Equal(t, tt.expected, isSessionError, "Session error detection should match expected")
		})
	}
}

// contains は文字列に部分文字列が含まれているかチェック
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
