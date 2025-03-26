// internal/infrastructure/client/master_data_client_impl.go
package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type masterDataClientImpl struct {
	client *TachibanaClient
	logger *zap.Logger // 追加
}

func (m *masterDataClientImpl) DownloadMasterData(ctx context.Context, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(m.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMEventDownload"
	req.P_no = m.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	// URLクエリパラメータに設定
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	fmt.Println("---------------------------------")
	// req.URL をデコードして表示
	decodedURL, _ := url.QueryUnescape(httpReq.URL.String())
	// logger.Debug("Decoded URL:", zap.String("decodedUrl", decodedURL))
	fmt.Println("Decoded URL:", decodedURL)
	fmt.Println("---------------------------------")

	// 4. リクエストの送信 (SendRequestを直接使わず、専用の処理を行う)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "download master data failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
	}

	// 5. 配信されるマスタデータを受信する
	res := &response.ResDownloadMaster{}

	// 6. タイムアウトを設定
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second) // 180秒でタイムアウト
	defer cancel()

	// ファイルを作成
	file, err := os.Create("raw_response.txt")
	if err != nil {
		fmt.Println("ファイル作成エラー:", err)
		//return nil, fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()

	// 7. レスポンスボディを1行ずつ読み込む
	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if err := ctx.Err(); err != nil {
			//return nil, errors.Wrap(err, "context canceled") // タイムアウトチェック
			break
		}

		// ファイルに書き込む
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("ファイル書き込みエラー:", err)
			//return nil, fmt.Errorf("ファイル書き込みエラー: %w", err)
		}

		// Shift-JIS から UTF-8 への変換
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), []byte(line))
		if err != nil {
			m.logger.Warn("shift-jis decode error", zap.Error(err))
			continue
		}

		// JSON としてデコード
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(bodyUTF8), &item); err != nil {
			m.logger.Warn("Failed to unmarshal line", zap.Error(err), zap.String("line", line))
			continue // デコードに失敗したらスキップ
		}

		// sCLMID キーの存在確認
		sCLMID, ok := item["sCLMID"].(string)
		if !ok {
			m.logger.Warn("sCLMID not found in response item", zap.Any("item", item))
			continue
		}

		// sCLMID の値に応じて処理
		switch sCLMID {
		case "CLMSystemStatus":
			var systemStatus response.ResSystemStatus
			if err := convertMapToStruct(item, &systemStatus, ""); err != nil {
				m.logger.Error("failed to map SystemStatus", zap.Error(err))
				continue
			}
			res.SystemStatus = systemStatus

		case "CLMEventDownloadComplete":
			fmt.Println("CLMEventDownloadComplete")
			return res, nil // 正常終了
		default:
			m.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	if err := scanner.Err(); err != nil {
		return res, errors.Wrap(err, "scanner error")
	}
	m.logger.Info("DownloadMasterData: タイムアウト") // ログ出力
	return res, errors.New("タイムアウト")
}
