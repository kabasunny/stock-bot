// internal/infrastructure/client/util.go
package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// SendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
// internal/infrastructure/client/util.go

// SendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func SendRequest(
	req *http.Request,
	maxRetries int,
	logger *zap.Logger,
) (map[string]interface{}, error) { // 引数をシンプルに
	var response map[string]interface{}

	// retryDoに渡す関数
	retryFunc := func(client *http.Client, decodeFunc func([]byte, interface{}) error) (*http.Response, error) {
		//timeoutコンテキストを作成
		req, cancel := withContextAndTimeout(req, 60*time.Second)
		defer cancel()

		resp, err := client.Do(req) //clientは、http.Client{}
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return resp, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // 読み込み終わったらすぐにクローズ
		if err != nil {
			return resp, fmt.Errorf("response body read error: %w", err)
		}

		logRequestAndResponse(req, body, logger)

		if err := decodeFunc(body, &response); err != nil {
			return resp, fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
		return resp, nil
	}

	decodeFunc := func(body []byte, v interface{}) error {
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
		if err != nil {
			return fmt.Errorf("shift-jis decode error: %w", err)
		}
		return json.Unmarshal(bodyUTF8, v) // UTF-8 でデコード
	}

	resp, err := retryDo(retryFunc, maxRetries, 2*time.Second, &http.Client{}, decodeFunc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return response, nil
}

// ConvertResponse は、map[string]interface{} をレスポンスDTOに変換する
func ConvertResponse[T any](respMap map[string]interface{}) (*T, error) {
	if resultCode, ok := respMap["sResultCode"].(string); ok && resultCode != "0" {
		resultText := ""
		if rt, ok := respMap["sResultText"].(string); ok {
			resultText = rt
		}
		return nil, errors.Errorf("API error: ResultCode=%s, ResultText=%s", resultCode, resultText)
	}

	resBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response map to json")
	}
	var res T
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response json")
	}
	return &res, nil
}

// structToMapString は、構造体を map[string]string に変換する
func structToMapString(data interface{}) (map[string]string, error) {
	params := make(map[string]string)
	reqBytes, err := json.Marshal(data) // 一度JSONに変換
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request to json")
	}
	var reqMap map[string]interface{}                         // interface{} のmap
	if err := json.Unmarshal(reqBytes, &reqMap); err != nil { // JSONからmapに変換
		return nil, errors.Wrap(err, "failed to unmarshal request json to map")
	}
	for k, v := range reqMap {
		if s, ok := v.(string); ok { // string の場合のみ
			params[k] = s
		}
		// string 以外は無視 (float64 などが入ってくる可能性がある)
	}
	return params, nil
}

// formatSDDate は、time.Time を "YYYY.MM.DD-HH:MM:SS.TTT" 形式の文字列に変換します。
func formatSDDate(t time.Time) string {
	return t.Format("2006.01.02-15:04:05.000")
}

// retryDo, withContextAndTimeout, logRequestAndResponse は変更なし (省略)
// 省略したretryDo, withContextAndTimeout, logRequestAndResponse は、前に提示したコードと同じです。

// retryDo は、指定された関数をリトライする(変更なし)
func retryDo(fn func(*http.Client, func([]byte, interface{}) error) (*http.Response, error), maxRetries int, interval time.Duration, client *http.Client, decodeFunc func([]byte, interface{}) error) (*http.Response, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		resp, err := fn(client, decodeFunc) //clientを渡す
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if i < maxRetries-1 { // 最後のリトライでなければ
			time.Sleep(interval) // ちょっと待つ
		}
	}
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// withContextAndTimeout は、リクエストにコンテキストとタイムアウトを設定する
func withContextAndTimeout(req *http.Request, timeout time.Duration) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), timeout) //timeoutを設定
	return req.WithContext(ctx), cancel
}

// logRequestAndResponse は、リクエストとレスポンスをログに出力する
func logRequestAndResponse(req *http.Request, respBody []byte, logger *zap.Logger) {
	// リクエスト情報のログ出力
	logger.Debug("Request:",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", req.Header), // リクエストヘッダー全体
	)

	// リクエストボディのログ出力 (存在する場合)
	if req.Body != nil {
		// リクエストボディを読み取り、ログ出力後に再度設定
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 再設定

		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), bodyBytes)
		if err == nil {
			logger.Debug("Request Body (UTF-8):", zap.String("body", string(bodyUTF8)))
		} else {
			logger.Error("Failed to decode request body", zap.Error(err))
		}
	}

	// レスポンスボディのログ出力
	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), respBody)
	if err == nil {
		logger.Debug("Raw Response Body (UTF-8, one line per JSON):")
		scanner := bufio.NewScanner(strings.NewReader(string(bodyUTF8)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue // 空行はスキップ
			}
			var js map[string]interface{}
			if json.Unmarshal([]byte(line), &js) == nil {
				logger.Debug("", zap.String("line", line)) // JSON として有効な行のみ出力
			} else {
				logger.Debug("Invalid JSON:", zap.String("line", line))
			}
		}
		if err := scanner.Err(); err != nil {
			logger.Error("Error scanning response body", zap.Error(err))
		}
	} else {
		logger.Error("Failed to decode response body", zap.Error(err))
	}
}

// RetryDo は、HTTP リクエストをリトライ付きで実行する(変更なし)
func RetryDo(
	retryFunc func(*http.Client, func([]byte, interface{}) error) (*http.Response, error), // decodeFuncの型修正
	maxRetries int,
	initialBackoff time.Duration,
	client *http.Client, // http.Client を引数で渡す
	decodeFunc func([]byte, interface{}) error, // デコード関数を引数で渡す []byteに変更
) (*http.Response, error) {
	var resp *http.Response
	var err error

	for retries := 1; retries <= maxRetries; retries++ {
		resp, err = retryFunc(client, decodeFunc)

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil // 成功時: エラーがなく、ステータスコードが200の場合
		}

		if retries < maxRetries {
			// 指数バックオフを計算
			// 回数が増すごとに間隔が広くなる
			// 初期遅延時間に対して2の乗数でリトライ間隔を増加 (例: 2秒, 4秒, 8秒...)
			backoff := time.Duration(math.Pow(2, float64(retries))) * initialBackoff
			// 計算したリトライ間隔の時間だけ待機
			time.Sleep(backoff)

			// レスポンスが存在し、かつそのボディがまだ閉じられていない場合は閉じる
			// これはリソースリークを防ぐための重要なステップ
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}
	}

	if resp != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %v, last status code: %d", maxRetries+1, err, resp.StatusCode)
	}
	return nil, fmt.Errorf("HTTP request failed after %d retries: last error: %w", maxRetries+1, err)
}
