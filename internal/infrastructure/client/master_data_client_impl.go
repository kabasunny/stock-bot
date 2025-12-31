// internal/infrastructure/client/master_data_client_impl.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv" // 追加
	"time"

	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type masterDataClientImpl struct {
	client *TachibanaClientImpl
}

func (m *masterDataClientImpl) DownloadMasterData(ctx context.Context, session *Session, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	// 1. リクエストURLの作成
	u := session.MasterURL

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMEventDownload"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10)
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
	decodedURL, _ := url.QueryUnescape(httpReq.URL.String())
	slog.Debug("Decoded URL", slog.String("decodedUrl", decodedURL))

	// 4. リクエストの送信 (SendRequestを直接使わず、専用の処理を行う)
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}
	resp, err := tempClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "download master data failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API のステータスコードが200以外のためエラー: %d", resp.StatusCode)
	}

	// 5. 配信されるマスタデータを受信する
	res := &response.ResDownloadMaster{}

		        // Shift-JISからUTF-8への変換リーダーを作成
				utf8Reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
			
				var buf bytes.Buffer
				chunk := make([]byte, 4096) // チャンクサイズ
			
				for {
					n, err := utf8Reader.Read(chunk)
					if n > 0 {
						buf.Write(chunk[:n]) // バッファに書き込む
					}
	
			// buf の内容から完全な JSON オブジェクトを抽出して処理
			// Pythonのサンプルに倣い、`}` を区切りとしてデコードを試みる
			for {
				jsonBytes := buf.Bytes()
				// JSONオブジェクトの終端 '}' を探す
				// ただし、ネストしたJSONオブジェクトの `}` が誤検知されないよう、
				// 最も外側の `}` を探す必要があるが、簡易的に最後の `}` を探す。
				// APIからのデータはトップレベルのJSONオブジェクトが連続して送られてくると仮定。
				idx := bytes.LastIndexByte(jsonBytes, '}')
				if idx == -1 {
					// 完全なJSONオブジェクトの終端が見つからない場合は、次のチャンクを待つ
					break
				}
	
				// 見つかった '}' の位置までのバイト列を取得し、デコードを試みる
				potentialJSON := jsonBytes[:idx+1]
				var item map[string]interface{}
				if unmarshalErr := json.Unmarshal(potentialJSON, &item); unmarshalErr == nil {
					// デコード成功。バッファから処理済みの部分を削除
					buf.Next(idx + 1)
	
					// sCLMID の値に応じて処理 (既存のswitch文)
					sCLMID, ok := item["sCLMID"].(string)
					if !ok {
						// slog.Warn("sCLMID not found in response item", slog.Any("item", item)) // デバッグ用
						continue
					}
					switch sCLMID {
					case "CLMSystemStatus":
						var systemStatus response.ResSystemStatus
						if err := convertMapToStruct(item, &systemStatus, ""); err != nil {
							slog.Error("failed to map SystemStatus", slog.Any("error", err))
							continue
						}
						res.SystemStatus = systemStatus
	
					case "CLMDateZyouhou":
						var dateInfo response.ResDateInfo
						if err := convertMapToStruct(item, &dateInfo, ""); err != nil {
							slog.Error("failed to map DateInfo", slog.Any("error", err))
							continue
						}
						res.DateInfo = append(res.DateInfo, dateInfo)
	
					case "CLMYobine":
						var tickRule response.ResTickRule
						if err := convertMapToStruct(item, &tickRule, ""); err != nil {
							slog.Error("failed to map TickRule", slog.Any("error", err))
							continue
						}
						res.TickRule = append(res.TickRule, tickRule)
	
					case "CLMUnyouStatus":
						var operationStatus response.ResOperationStatus
						if err := convertMapToStruct(item, &operationStatus, ""); err != nil {
							slog.Error("failed to map OperationStatus", slog.Any("error", err))
							continue
						}
						res.OperationStatus = append(res.OperationStatus, operationStatus)
	
					case "CLMUnyouStatusKabu":
						var operationStatusStock response.ResOperationStatus
						if err := convertMapToStruct(item, &operationStatusStock, ""); err != nil {
							slog.Error("failed to map OperationStatusKabu", slog.Any("error", err))
							continue
						}
						res.OperationStatusStock = append(res.OperationStatusStock, operationStatusStock)
	
					case "CLMUnyouStatusHasei":
						var operationStatusDerivative response.ResOperationStatus
						if err := convertMapToStruct(item, &operationStatusDerivative, ""); err != nil {
							slog.Error("failed to map OperationStatusHasei", slog.Any("error", err))
							continue
						}
						res.OperationStatusDerivative = append(res.OperationStatusDerivative, operationStatusDerivative)
	
					case "CLMIssueMstKabu":
						var stockMaster response.ResStockMaster
						if err := convertMapToStruct(item, &stockMaster, ""); err != nil {
							slog.Error("failed to map StockMaster", slog.Any("error", err))
							continue
						}
						res.StockMaster = append(res.StockMaster, stockMaster)
	
					case "CLMIssueSizyouMstKabu":
						var stockMarketMaster response.ResStockMarketMaster
						if err := convertMapToStruct(item, &stockMarketMaster, ""); err != nil {
							slog.Error("failed to map StockMarketMaster", slog.Any("error", err))
							continue
						}
						res.StockMarketMaster = append(res.StockMarketMaster, stockMarketMaster)
	
					case "CLMIssueSizyouKiseiKabu":
						var stockIssueRegulation response.ResStockIssueRegulation
						if err := convertMapToStruct(item, &stockIssueRegulation, ""); err != nil {
							slog.Error("failed to map StockIssueRegulation", slog.Any("error", err))
							continue
						}
						res.StockIssueRegulation = append(res.StockIssueRegulation, stockIssueRegulation)
	
					case "CLMIssueMstSak":
						var futureMaster response.ResFutureMaster
						if err := convertMapToStruct(item, &futureMaster, ""); err != nil {
							slog.Error("failed to map FutureMaster", slog.Any("error", err))
							continue
						}
						res.FutureMaster = append(res.FutureMaster, futureMaster)
	
					case "CLMIssueMstOp":
						var optionMaster response.ResOptionMaster
						if err := convertMapToStruct(item, &optionMaster, ""); err != nil {
							slog.Error("failed to map OptionMaster", slog.Any("error", err))
							continue
						}
						res.OptionMaster = append(res.OptionMaster, optionMaster)
	
					case "CLMIssueSizyouKiseiHasei":
						var futureOptionRegulation response.ResFutureOptionRegulation
						if err := convertMapToStruct(item, &futureOptionRegulation, ""); err != nil {
							slog.Error("failed to map FutureOptionRegulation", slog.Any("error", err))
							continue
						}
						res.FutureOptionRegulation = append(res.FutureOptionRegulation, futureOptionRegulation)
	
					case "CLMDaiyouKakeme":
						var marginRate response.ResMarginRate
						if err := convertMapToStruct(item, &marginRate, ""); err != nil {
							slog.Error("failed to map MarginRate", slog.Any("error", err))
							continue
						}
						res.MarginRate = append(res.MarginRate, marginRate)
	
					case "CLMHosyoukinMst":
						var marginMaster response.ResMarginMaster
						if err := convertMapToStruct(item, &marginMaster, ""); err != nil {
							slog.Error("failed to map MarginMaster", slog.Any("error", err))
							continue
						}
						res.MarginMaster = append(res.MarginMaster, marginMaster)
	
					case "CLMOrderErrReason":
						var errorReason response.ResErrorReason
						if err := convertMapToStruct(item, &errorReason, ""); err != nil {
							slog.Error("failed to map ErrorReason", slog.Any("error", err))
							continue
						}
						res.ErrorReason = append(res.ErrorReason, errorReason)
	
					case "CLMEventDownloadComplete":
						slog.Info("CLMEventDownloadComplete received, download finished.")
						return res, nil // 正常終了
	
					default:
						slog.Warn("Unknown master data type", slog.String("sCLMID", sCLMID))
					}
				} else {
					// デコード失敗（不完全なJSONか、他のエラー）。この `}` はJSONオブジェクトの終端ではない可能性があるため、
					// バッファをクリアせず、次のチャンクの読み込みを待つ
					// 複数行のJSONオブジェクトの場合に対応するため、ここでは break しない
					// ただし、このエラーログは非常に重要。
					slog.Debug("Failed to unmarshal potential JSON segment", slog.Any("error", unmarshalErr), slog.String("segment", string(potentialJSON)))
					// このJSONセグメントが解析できなかったので、次の可能性を探すか、さらにデータを読み込むためにループを抜ける
					break
				}
			}
	
			if err == io.EOF {
				break // EOF ならループ終了
			}
			if err != nil {
				slog.Error("Error reading response body", slog.Any("error", err))
				return nil, errors.Wrap(err, "error reading response body")
			}
		}
	
		// EOF に達したが、CLMEventDownloadComplete が受信されていない場合
		// 残りのバッファにも完了通知がないか確認し、あれば処理する
		finalBufContent := buf.Bytes()
		if bytes.Contains(finalBufContent, []byte("CLMEventDownloadComplete")) {
			slog.Info("CLMEventDownloadComplete found in final buffer, download finished.")
			return res, nil
		}
	
		slog.Error("DownloadMasterData stream finished without CLMEventDownloadComplete signal")
		return nil, errors.New("download master data stream finished without complete signal")}

func (m *masterDataClientImpl) GetMasterDataQuery(ctx context.Context, session *Session, req request.ReqGetMasterData) (*response.ResGetMasterData, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetMasterData"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get master data query with post failed")
	}

	res, err := ConvertResponse[response.ResGetMasterData](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetNewsHeader(ctx context.Context, session *Session, req request.ReqGetNewsHead) (*response.ResGetNewsHeader, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetNewsHead"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get news header with post failed")
	}

	res, err := ConvertResponse[response.ResGetNewsHeader](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetNewsBody(ctx context.Context, session *Session, req request.ReqGetNewsBody) (*response.ResGetNewsBody, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetNewsBody"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get news body with post failed")
	}

	res, err := ConvertResponse[response.ResGetNewsBody](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetIssueDetail(ctx context.Context, session *Session, req request.ReqGetIssueDetail) (*response.ResGetIssueDetail, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetIssueDetail"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get issue detail with post failed")
	}

	res, err := ConvertResponse[response.ResGetIssueDetail](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetMarginInfo(ctx context.Context, session *Session, req request.ReqGetMarginInfo) (*response.ResGetMarginInfo, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetSyoukinZan"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get margin info with post failed")
	}

	res, err := ConvertResponse[response.ResGetMarginInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetCreditInfo(ctx context.Context, session *Session, req request.ReqGetCreditInfo) (*response.ResGetCreditInfo, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetShinyouZan"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get credit info with post failed")
	}

	res, err := ConvertResponse[response.ResGetCreditInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *masterDataClientImpl) GetMarginPremiumInfo(ctx context.Context, session *Session, req request.ReqGetMarginPremiumInfo) (*response.ResGetMarginPremiumInfo, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	u, err := url.Parse(session.MasterURL) // sessionからMasterURLを取得
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse master URL from session")
	}

	req.CLMID = "CLMMfdsGetHibuInfo"
	req.P_no = strconv.FormatInt(int64(session.GetPNo()), 10) // sessionからp_noを取得
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	// 認証済みセッションのCookieJarを持つ一時的なhttp.Clientを作成
	tempClient := &http.Client{
		Jar: session.CookieJar,
	}

	respMap, err := SendRequest(tempClient, httpReq, 3) // tempClient を使用
	if err != nil {
		return nil, errors.Wrap(err, "get margin premium info with post failed")
	}

	res, err := ConvertResponse[response.ResGetMarginPremiumInfo](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
