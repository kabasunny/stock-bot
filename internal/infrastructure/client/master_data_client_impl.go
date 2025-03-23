// internal/infrastructure/client/master_data_client_impl.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	request_master "stock-bot/internal/infrastructure/client/dto/master/request"
	response_master "stock-bot/internal/infrastructure/client/dto/master/response"

	"github.com/cockroachdb/errors"
)

type masterDataClientImpl struct {
	client *TachibanaClient
}

// GetMasterData は、各種マスタ情報をリアルタイム配信でダウンロードします。
func (m *masterDataClientImpl) GetMasterData(ctx context.Context, req request_master.ReqGetMasterData) (*response_master.ResGetMasterData, error) {
	fmt.Println("GetMasterData")
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}
	// 1. リクエストURLの作成
	requestURL := m.client.loginInfo.MasterURL + "/master/download" // masterURLを使用

	// 2. リクエストボディの作成
	req.CLMID = "CLMDownLoad"
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}
	// 3. HTTPリクエストの作成
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	// Content-Type 設定
	httpReq.Header.Set("Content-Type", "application/json")

	// 4. リクエストの送信
	m.client.p_NoMu.Lock()
	m.client.p_no++
	httpReq.Header.Set("p_no", fmt.Sprintf("%d", m.client.p_no))
	m.client.p_NoMu.Unlock()
	resp, err := http.DefaultClient.Do(httpReq) // ここで実際にリクエストを送信
	if err != nil {
		return nil, errors.Wrap(err, "failed to send http request")
	}
	defer resp.Body.Close()

	// 5. レスポンスの処理
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res response_master.ResGetMasterData // DTO (Data Transfer Object) を定義
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	// 6. レスポンスを返す（必要であればドメインモデルに変換）
	return &res, nil
}

func (m *masterDataClientImpl) GetMasterDataQuery(ctx context.Context, req request_master.ReqGetMasterData) (*response_master.ResGetMasterData, error) {
	fmt.Println("GetMasterDataQuery")
	if !m.client.loggined {
		return nil, fmt.Errorf("not logged in")
	}
	// 1. リクエストURLの作成
	u, err := url.Parse(m.client.loginInfo.RequestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request URL: %w", err)
	}
	u.Path += "/master/inquiry/master" // 正しいエンドポイント

	// 2. リクエストボディの作成
	req.CLMID = "CLMMfdsGetMasterData"
	requestBody, err := json.Marshal(req)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 3. HTTPリクエストの作成
	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 4. リクエストの送信
	m.client.p_NoMu.Lock()
	m.client.p_no++
	httpReq.Header.Set("p_no", fmt.Sprintf("%d", m.client.p_no))
	m.client.p_NoMu.Unlock()
	resp, err := http.DefaultClient.Do(httpReq) // http.DefaultClient を使用

	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	// 5. レスポンスの処理
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res response_master.ResGetMasterData
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// 6. レスポンスを返す
	return &res, nil
}

func (m *masterDataClientImpl) GetNewsHeader(ctx context.Context, req request_master.ReqGetNewsHead) (*response_master.ResGetNewsHeader, error) {
	fmt.Println("Dummy GetNewsHeader")
	return nil, nil
}
func (m *masterDataClientImpl) GetNewsBody(ctx context.Context, req request_master.ReqGetNewsBody) (*response_master.ResGetNewsBody, error) {
	fmt.Println("Dummy GetNewsBody")
	return nil, nil
}
func (m *masterDataClientImpl) GetIssueDetail(ctx context.Context, req request_master.ReqGetIssueDetail) (*response_master.ResGetIssueDetail, error) {
	fmt.Println("Dummy GetIssueDetail")
	return nil, nil
}
func (m *masterDataClientImpl) GetMarginInfo(ctx context.Context, req request_master.ReqGetMarginInfo) (*response_master.ResGetMarginInfo, error) {
	fmt.Println("Dummy GetMarginInfo")
	return nil, nil
}
func (m *masterDataClientImpl) GetCreditInfo(ctx context.Context, req request_master.ReqGetCreditInfo) (*response_master.ResGetCreditInfo, error) {
	fmt.Println("Dummy GetCreditInfo")
	return nil, nil
}

func (m *masterDataClientImpl) GetMarginPremiumInfo(ctx context.Context, req request_master.ReqGetMarginPremiumInfo) (*response_master.ResGetMarginPremiumInfo, error) {
	fmt.Println("Dummy GetMarginPremiumInfo")
	return nil, nil
}
