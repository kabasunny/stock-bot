// internal/infrastructure/client/master_data_client_impl.go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type masterDataClientImpl struct {
	client *TachibanaClientImpl
	logger *zap.Logger // 追加
}

func (m *masterDataClientImpl) DownloadMasterData(ctx context.Context, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

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

	// URLエンコード (GETリクエスト)
	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
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
	// ctx, cancel := context.WithTimeout(ctx, 180*time.Second) // 180秒でタイムアウト
	// defer cancel()

	// ファイルを作成
	file, err := os.Create("raw_response.txt")
	if err != nil {
		fmt.Println("ファイル作成エラー:", err)
		//return nil, fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()

	// 7. レスポンスボディを文字列として読み込む
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	bodyString := string(bodyBytes)

	// JSONオブジェクトの区切り文字で分割
	jsonStrings := strings.Split(bodyString, "}{")

	// 分割された各部分を個別のJSONオブジェクトとしてデコード
	for i, jsonString := range jsonStrings {
		// 最初の要素と最後の要素は、それぞれ '{' と '}' で始まる/終わるように調整
		if i == 0 {
			jsonString = jsonString + "}"
		} else if i == len(jsonStrings)-1 {
			jsonString = "{" + jsonString
		} else {
			jsonString = "{" + jsonString + "}"
		}

		// Shift-JIS から UTF-8 への変換
		bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), []byte(jsonString))
		if err != nil {
			m.logger.Warn("shift-jis decode error", zap.Error(err))
			continue
		}

		// JSONとしてデコード
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(bodyUTF8), &item); err != nil {
			m.logger.Warn("Failed to unmarshal line", zap.Error(err), zap.String("line", jsonString))
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

		case "CLMDateZyouhou":
			var dateInfo response.ResDateInfo
			if err := convertMapToStruct(item, &dateInfo, ""); err != nil {
				m.logger.Error("failed to map DateInfo", zap.Error(err))
				continue
			}
			res.DateInfo = append(res.DateInfo, dateInfo)

		case "CLMYobine":
			var tickRule response.ResTickRule
			if err := convertMapToStruct(item, &tickRule, ""); err != nil {
				m.logger.Error("failed to map TickRule", zap.Error(err))
				continue
			}
			res.TickRule = append(res.TickRule, tickRule)

		case "CLMUnyouStatus":
			var operationStatus response.ResOperationStatus
			if err := convertMapToStruct(item, &operationStatus, ""); err != nil {
				m.logger.Error("failed to map OperationStatus", zap.Error(err))
				continue
			}
			res.OperationStatus = append(res.OperationStatus, operationStatus)

		case "CLMUnyouStatusKabu":
			var operationStatusStock response.ResOperationStatus
			if err := convertMapToStruct(item, &operationStatusStock, ""); err != nil {
				m.logger.Error("failed to map OperationStatusKabu", zap.Error(err))
				continue
			}
			res.OperationStatusStock = append(res.OperationStatusStock, operationStatusStock)

		case "CLMUnyouStatusHasei":
			var operationStatusDerivative response.ResOperationStatus
			if err := convertMapToStruct(item, &operationStatusDerivative, ""); err != nil {
				m.logger.Error("failed to map OperationStatusHasei", zap.Error(err))
				continue
			}
			res.OperationStatusDerivative = append(res.OperationStatusDerivative, operationStatusDerivative)

		case "CLMIssueMstKabu":
			var stockMaster response.ResStockMaster
			if err := convertMapToStruct(item, &stockMaster, ""); err != nil {
				m.logger.Error("failed to map StockMaster", zap.Error(err))
				continue
			}
			res.StockMaster = append(res.StockMaster, stockMaster)

		case "CLMIssueSizyouMstKabu":
			var stockMarketMaster response.ResStockMarketMaster
			if err := convertMapToStruct(item, &stockMarketMaster, ""); err != nil {
				m.logger.Error("failed to map StockMarketMaster", zap.Error(err))
				continue
			}
			res.StockMarketMaster = append(res.StockMarketMaster, stockMarketMaster)

		case "CLMIssueSizyouKiseiKabu":
			var stockIssueRegulation response.ResStockIssueRegulation
			if err := convertMapToStruct(item, &stockIssueRegulation, ""); err != nil {
				m.logger.Error("failed to map StockIssueRegulation", zap.Error(err))
				continue
			}
			res.StockIssueRegulation = append(res.StockIssueRegulation, stockIssueRegulation)

		case "CLMIssueMstSak":
			var futureMaster response.ResFutureMaster
			if err := convertMapToStruct(item, &futureMaster, ""); err != nil {
				m.logger.Error("failed to map FutureMaster", zap.Error(err))
				continue
			}
			res.FutureMaster = append(res.FutureMaster, futureMaster)

		case "CLMIssueMstOp":
			var optionMaster response.ResOptionMaster
			if err := convertMapToStruct(item, &optionMaster, ""); err != nil {
				m.logger.Error("failed to map OptionMaster", zap.Error(err))
				continue
			}
			res.OptionMaster = append(res.OptionMaster, optionMaster)

		case "CLMIssueSizyouKiseiHasei":
			var futureOptionRegulation response.ResFutureOptionRegulation
			if err := convertMapToStruct(item, &futureOptionRegulation, ""); err != nil {
				m.logger.Error("failed to map FutureOptionRegulation", zap.Error(err))
				continue
			}
			res.FutureOptionRegulation = append(res.FutureOptionRegulation, futureOptionRegulation)

		case "CLMDaiyouKakeme":
			var marginRate response.ResMarginRate
			if err := convertMapToStruct(item, &marginRate, ""); err != nil {
				m.logger.Error("failed to map MarginRate", zap.Error(err))
				continue
			}
			res.MarginRate = append(res.MarginRate, marginRate)

		case "CLMHosyoukinMst":
			var marginMaster response.ResMarginMaster
			if err := convertMapToStruct(item, &marginMaster, ""); err != nil {
				m.logger.Error("failed to map MarginMaster", zap.Error(err))
				continue
			}
			res.MarginMaster = append(res.MarginMaster, marginMaster)

		case "CLMOrderErrReason":
			var errorReason response.ResErrorReason
			if err := convertMapToStruct(item, &errorReason, ""); err != nil {
				m.logger.Error("failed to map ErrorReason", zap.Error(err))
				continue
			}
			res.ErrorReason = append(res.ErrorReason, errorReason)

		case "CLMEventDownloadComplete":
			fmt.Println("CLMEventDownloadComplete")
			return res, nil // 正常終了

		default:
			m.logger.Warn("Unknown master data type", zap.String("sCLMID", sCLMID))
		}
	}

	// タイムアウトまたはエラーが発生した場合、エラーを返す
	m.logger.Info("DownloadMasterData: タイムアウトまたはエラー") // ログ出力
	return res, errors.New("タイムアウトまたはエラー")
}

func (m *masterDataClientImpl) GetMasterDataQuery(ctx context.Context, req request.ReqGetMasterData) (*response.ResGetMasterData, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetMasterData"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get master data query failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetMasterData](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetNewsHeader(ctx context.Context, req request.ReqGetNewsHead) (*response.ResGetNewsHeader, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetNewsHead"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get news header failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetNewsHeader](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetNewsBody(ctx context.Context, req request.ReqGetNewsBody) (*response.ResGetNewsBody, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetNewsBody"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get news body failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetNewsBody](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetIssueDetail(ctx context.Context, req request.ReqGetIssueDetail) (*response.ResGetIssueDetail, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetIssueDetail"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get issue detail failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetIssueDetail](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetMarginInfo(ctx context.Context, req request.ReqGetMarginInfo) (*response.ResGetMarginInfo, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetSyoukinZan"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get margin info failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetMarginInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetCreditInfo(ctx context.Context, req request.ReqGetCreditInfo) (*response.ResGetCreditInfo, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetShinyouZan"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get credit info failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetCreditInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetMarginPremiumInfo(ctx context.Context, req request.ReqGetMarginPremiumInfo) (*response.ResGetMarginPremiumInfo, error) {
	if !m.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u := m.client.loginInfo.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMMfdsGetHibuInfo"
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
	requestURL := u + "?" + encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET に変更
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3, m.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get margin premium info failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResGetMarginPremiumInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
