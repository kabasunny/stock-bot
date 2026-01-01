package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"stock-bot/domain/model"
	"time"
)

// HTTPStrategyClient は戦略管理APIとHTTP通信するクライアント
type HTTPStrategyClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewHTTPStrategyClient は新しいHTTPStrategyClientを作成
func NewHTTPStrategyClient(baseURL string, logger *slog.Logger) *HTTPStrategyClient {
	return &HTTPStrategyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// StrategyResponse は戦略APIのレスポンス型
type StrategyResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Description *string                `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"`
	RiskLimits  map[string]interface{} `json:"risk_limits"`
	Statistics  map[string]interface{} `json:"statistics"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	CreatedBy   *string                `json:"created_by,omitempty"`
}

// StrategyListResponse は戦略一覧のレスポンス型
type StrategyListResponse struct {
	Strategies []*StrategyResponse `json:"strategies"`
}

// GetStrategy は指定されたIDの戦略を取得
func (c *HTTPStrategyClient) GetStrategy(ctx context.Context, id string) (*model.Strategy, error) {
	url := fmt.Sprintf("%s/strategies/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var strategyResp StrategyResponse
	if err := json.NewDecoder(resp.Body).Decode(&strategyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	strategy, err := c.convertToModel(&strategyResp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	c.logger.Info("strategy retrieved via HTTP", slog.String("id", id), slog.String("name", strategy.Name))
	return strategy, nil
}

// ListActiveStrategies はアクティブな戦略一覧を取得
func (c *HTTPStrategyClient) ListActiveStrategies(ctx context.Context) ([]*model.Strategy, error) {
	url := fmt.Sprintf("%s/strategies/active", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var listResp StrategyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	strategies := make([]*model.Strategy, len(listResp.Strategies))
	for i, strategyResp := range listResp.Strategies {
		strategy, err := c.convertToModel(strategyResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert strategy %d: %w", i, err)
		}
		strategies[i] = strategy
	}

	c.logger.Info("active strategies retrieved via HTTP", slog.Int("count", len(strategies)))
	return strategies, nil
}

// UpdateStatistics は戦略の統計情報を更新
func (c *HTTPStrategyClient) UpdateStatistics(ctx context.Context, strategyID string, pl float64, isWin bool) error {
	url := fmt.Sprintf("%s/strategies/%s/statistics", c.baseURL, strategyID)

	payload := map[string]interface{}{
		"pl":     pl,
		"is_win": isWin,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	c.logger.Info("strategy statistics updated via HTTP",
		slog.String("strategy_id", strategyID),
		slog.Float64("pl", pl),
		slog.Bool("is_win", isWin))

	return nil
}

// ActivateStrategy は戦略をアクティブ化
func (c *HTTPStrategyClient) ActivateStrategy(ctx context.Context, strategyID string) error {
	url := fmt.Sprintf("%s/strategies/%s/activate", c.baseURL, strategyID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	c.logger.Info("strategy activated via HTTP", slog.String("strategy_id", strategyID))
	return nil
}

// DeactivateStrategy は戦略を非アクティブ化
func (c *HTTPStrategyClient) DeactivateStrategy(ctx context.Context, strategyID string) error {
	url := fmt.Sprintf("%s/strategies/%s/deactivate", c.baseURL, strategyID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	c.logger.Info("strategy deactivated via HTTP", slog.String("strategy_id", strategyID))
	return nil
}

// convertToModel はAPIレスポンスをドメインモデルに変換
func (c *HTTPStrategyClient) convertToModel(resp *StrategyResponse) (*model.Strategy, error) {
	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, resp.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	strategy := &model.Strategy{
		ID:        resp.ID,
		Name:      resp.Name,
		Type:      model.StrategyType(resp.Type),
		Status:    model.StrategyStatus(resp.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if resp.Description != nil {
		strategy.Description = *resp.Description
	}

	if resp.CreatedBy != nil {
		strategy.CreatedBy = *resp.CreatedBy
	}

	// Config, RiskLimits, Statisticsの変換は必要に応じて実装
	// 現在は基本フィールドのみ対応

	return strategy, nil
}
