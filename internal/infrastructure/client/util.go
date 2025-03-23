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

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// formatSDDate は、time.Time を "YYYY.MM.DD-HH:MM:SS.TTT" 形式の文字列に変換します。
func formatSDDate(t time.Time) string {
	return t.Format("2006.01.02-15:04:05.000")
}

// SendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func SendRequest(req *http.Request, maxRetries int, tc *TachibanaClient, logger *zap.Logger) (map[string]interface{}, error) {
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
		// fmt.Print(body)
		resp.Body.Close() // 読み込み終わったらすぐにクローズ
		if err != nil {
			return resp, fmt.Errorf("response body read error: %w", err)
		}

		//loggerを使って、リクエストとレスポンスをログに出力
		logRequestAndResponse(req, body, logger)

		if err := decodeFunc(body, &response); err != nil {
			return resp, fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
		return resp, nil
	}

	// デコード関数を定義 (Shift-JIS から UTF-8 への変換)
	decodeFunc := func(body []byte, v interface{}) error { // 引数を io.Reader から []byte に変更
		// Shift-JISからUTF-8への変換
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
		if err != nil {
			return fmt.Errorf("shift-jis decode error: %w", err)
		}
		return json.Unmarshal(bodyUTF8, v) // UTF-8 でデコード
	}

	// reqのTimeoutを使うので、ここではClientを生成しない
	resp, err := retryDo(retryFunc, maxRetries, 2*time.Second, &http.Client{}, decodeFunc) //空のClientを渡す
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return response, nil
}

// retryDo は、指定された関数をリトライする
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

// retryDo は、HTTP リクエストをリトライ付きで実行する
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
