// internal/infrastructure/client/order_client_impl.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
)

type orderClientImpl struct {
	client *TachibanaClientImpl
}

func (o *orderClientImpl) NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuNewOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	// --- CLMKabuHensaiData スライスを JSON 配列形式の文字列に変換 ---
	var hensaiDataJSON string
	if len(req.CLMKabuHensaiData) > 0 {
		hensaiBytes, err := json.Marshal(req.CLMKabuHensaiData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal CLMKabuHensaiData to JSON")
		}
		hensaiDataJSON = string(hensaiBytes)
	}

	// --- スライスを除いたリクエスト構造体を作成 ---
	tempReq := req
	tempReq.CLMKabuHensaiData = nil // スライスを一旦除外

	// --- スライスを除いた部分を map[string]string に変換 ---
	params, err := structToMapString(tempReq)
	if err != nil {
		return nil, err
	}

	// --- CLMKabuHensaiData が空でない場合は、キーと値を追加 ---
	if hensaiDataJSON != "" {
		params["aCLMKabuHensaiData"] = hensaiDataJSON
	}

	// ★★★ 変更: map[string]string から JSON 文字列を組み立てる ★★★
	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		// 文字列の場合は、"key":"value" の形式にする
		// 文字列でない場合は、"key":value の形式にする (aCLMKabuHensaiData は文字列ではない)
		if k == "aCLMKabuHensaiData" {
			buf.WriteString(fmt.Sprintf(`"%s":%s`, k, v)) // aCLMKabuHensaiDataは文字列ではないので、そのまま
		} else {
			buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v)) // キーと値をダブルクォートで囲む
		}

	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()

	encodedPayload := url.QueryEscape(string(payloadJSON)) // エンコードするのはここ
	u.RawQuery = encodedPayload

	// 3. HTTPリクエストの作成 (GET)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3)

	if err != nil {
		return nil, errors.Wrap(err, "new order failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResNewOrder](respMap) //utilの関数
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderClientImpl) NewOrderWithPost(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuNewOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	// --- NewOrder(GET)と全く同じペイロード作成ロジック ---
	var hensaiDataJSON string
	if len(req.CLMKabuHensaiData) > 0 {
		hensaiBytes, err := json.Marshal(req.CLMKabuHensaiData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal CLMKabuHensaiData to JSON")
		}
		hensaiDataJSON = string(hensaiBytes)
	}

	tempReq := req
	tempReq.CLMKabuHensaiData = nil

	params, err := structToMapString(tempReq)
	if err != nil {
		return nil, err
	}

	if hensaiDataJSON != "" {
		params["aCLMKabuHensaiData"] = hensaiDataJSON
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for k, v := range params {
		if !first {
			buf.WriteString(",")
		}
		first = false
		if k == "aCLMKabuHensaiData" {
			buf.WriteString(fmt.Sprintf(`"%s":%s`, k, v))
		} else {
			buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
		}
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()
	// --- ペイロード作成ロジックここまで ---

	// --- SendRequestを直接呼び出すロジック ---
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "new order failed")
	}
	// --- SendRequest呼び出しここまで ---

	res, err := ConvertResponse[response.ResNewOrder](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderClientImpl) CorrectOrder(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(o.client.loginInfo.RequestURL) // RequestURL を使用
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMKabuCorrectOrder"        // CLMID を設定
	req.P_no = o.client.getPNo()             // クライアントから p_no を取得
	req.P_sd_date = formatSDDate(time.Now()) // システム日付を設定
	req.JsonOfmt = "4"                       // JSON出力フォーマット

	// 構造体を map[string]string に変換
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

	// 4. リクエストの送信
	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "correct order failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResCorrectOrder](respMap)
	if err != nil {
		return nil, err
	}
	if res.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s (errno: %s)", res.ResultText, res.ResultCode)
	}

	return res, nil
}

func (o *orderClientImpl) CorrectOrderWithPost(ctx context.Context, req request.ReqCorrectOrder) (*response.ResCorrectOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCorrectOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	// --- NewOrderWithPostと同様の手動ペイロード作成ロジック ---
	// CorrectOrderにはCLMKabuHensaiDataのようなスライスがないため、NewOrderよりシンプル
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
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v)) // 全て文字列として追加
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()
	// --- ペイロード作成ロジックここまで ---

	// --- SendRequestを直接呼び出すロジック ---
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "correct order failed")
	}
	// --- SendRequest呼び出しここまで ---

	res, err := ConvertResponse[response.ResCorrectOrder](respMap)
	if err != nil {
		return nil, err
	}
	if res.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s (errno: %s)", res.ResultText, res.ResultCode)
	}

	return res, nil
}

func (o *orderClientImpl) CancelOrder(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error) {
	// ほぼ CorrectOrder と同様の実装 (CLMID, レスポンスの型が異なる)
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}
	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCancelOrder"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "cancel order failed")
	}

	res, err := ConvertResponse[response.ResCancelOrder](respMap)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *orderClientImpl) CancelOrderWithPost(ctx context.Context, req request.ReqCancelOrder) (*response.ResCancelOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMKabuCancelOrder"
	req.RequestBase.P_no = o.client.getPNo()
	req.RequestBase.P_sd_date = formatSDDate(time.Now())
	req.RequestBase.JsonOfmt = "4"

	// --- CorrectOrderWithPostと同様の手動ペイロード作成ロジック ---
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
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v)) // 全て文字列として追加
	}
	buf.WriteString("}")
	payloadJSON := buf.Bytes()
	// --- ペイロード作成ロジックここまで ---

	// --- SendRequestを直接呼び出すロジック ---
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil
	}

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "cancel order failed")
	}
	// --- SendRequest呼び出しここまで ---

	res, err := ConvertResponse[response.ResCancelOrder](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderClientImpl) CancelOrderAll(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error) {
	// ほぼ CorrectOrder, CancelOrder と同様の実装 (CLMID, レスポンスの型が異なる)
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}
	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}
	req.CLMID = "CLMKabuCancelOrderAll"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "cancel all order failed")
	}

	res, err := ConvertResponse[response.ResCancelOrderAll](respMap)
	if err != nil {
		return nil, err
	}

		return res, nil

	}

	

	func (o *orderClientImpl) CancelOrderAllWithPost(ctx context.Context, req request.ReqCancelOrderAll) (*response.ResCancelOrderAll, error) {

		if !o.client.loggined {

			return nil, errors.New("not logged in")

		}

	

		u, err := url.Parse(o.client.loginInfo.RequestURL)

		if err != nil {

			return nil, errors.Wrap(err, "failed to parse request URL")

		}

	

		req.CLMID = "CLMKabuCancelOrderAll"

		req.RequestBase.P_no = o.client.getPNo()

		req.RequestBase.P_sd_date = formatSDDate(time.Now())

		req.RequestBase.JsonOfmt = "4"

	

		// --- これまでのPOST版と同様の手動ペイロード作成ロジック ---

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

			buf.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v)) // 全て文字列として追加

		}

		buf.WriteString("}")

		payloadJSON := buf.Bytes()

		// --- ペイロード作成ロジックここまで ---

	

		// --- SendRequestを直接呼び出すロジック ---

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payloadJSON))

		if err != nil {

			return nil, errors.Wrap(err, "failed to create http request")

		}

		httpReq.Header.Set("Content-Type", "application/json")

		httpReq.GetBody = func() (io.ReadCloser, error) {

			return io.NopCloser(bytes.NewBuffer(payloadJSON)), nil

		}

	

		respMap, err := SendRequest(httpReq, 3)

		if err != nil {

			return nil, errors.Wrap(err, "cancel all order failed")

		}

		// --- SendRequest呼び出しここまで ---

	

		res, err := ConvertResponse[response.ResCancelOrderAll](respMap)

		if err != nil {

			return nil, err

		}

	

		return res, nil

	}

	

	func (o *orderClientImpl) GetOrderList(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error) {

	

		if !o.client.loggined {

	

			return nil, errors.New("not logged in")

	

		}

	

		u, err := url.Parse(o.client.loginInfo.RequestURL)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "failed to parse request URL")

	

		}

	

	

	

		req.CLMID = "CLMOrderList"

	

		req.P_no = o.client.getPNo()

	

		req.P_sd_date = formatSDDate(time.Now())

	

		req.JsonOfmt = "4"

	

	

	

		params, err := structToMapString(req)

	

		if err != nil {

	

			return nil, err

	

		}

	

	

	

		payloadJSON, err := json.Marshal(params)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "failed to marshal request payload")

	

		}

	

		encodedPayload := url.QueryEscape(string(payloadJSON))

	

		u.RawQuery = encodedPayload

	

	

	

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "failed to create http request")

	

		}

	

	

	

		respMap, err := SendRequest(httpReq, 3)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "get order list failed")

	

		}

	

	

	

		res, err := ConvertResponse[response.ResOrderList](respMap)

	

		if err != nil {

	

			return nil, err

	

		}

	

		return res, nil

	

	}

	

	

	

	func (o *orderClientImpl) GetOrderListWithPost(ctx context.Context, req request.ReqOrderList) (*response.ResOrderList, error) {

	

		if !o.client.loggined {

	

			return nil, errors.New("not logged in")

	

		}

	

	

	

		u, err := url.Parse(o.client.loginInfo.RequestURL)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "failed to parse request URL")

	

		}

	

	

	

		req.CLMID = "CLMOrderList"

	

		req.P_no = o.client.getPNo()

	

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

	

	

	

		respMap, err := SendRequest(httpReq, 3)

	

		if err != nil {

	

			return nil, errors.Wrap(err, "get order list with post failed")

	

		}

	

	

	

		res, err := ConvertResponse[response.ResOrderList](respMap)

	

		if err != nil {

	

			return nil, err

	

		}

	

		return res, nil

	

	}

// internal/infrastructure/client/order_client_impl.go
func (o *orderClientImpl) GetOrderListDetail(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMOrderListDetail"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.JsonOfmt = "4"

	params, err := structToMapString(req)
	if err != nil {
		return nil, err
	}

	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}
	encodedPayload := url.QueryEscape(string(payloadJSON))
	u.RawQuery = encodedPayload

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get order list detail failed")
	}

	res, err := ConvertResponse[response.ResOrderListDetail](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderClientImpl) GetOrderListDetailWithPost(ctx context.Context, req request.ReqOrderListDetail) (*response.ResOrderListDetail, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	req.CLMID = "CLMOrderListDetail"
	req.P_no = o.client.getPNo()
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

	respMap, err := SendRequest(httpReq, 3)
	if err != nil {
		return nil, errors.Wrap(err, "get order list detail with post failed")
	}

	res, err := ConvertResponse[response.ResOrderListDetail](respMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}