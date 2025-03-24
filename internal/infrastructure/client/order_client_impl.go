// internal/infrastructure/client/order_client_impl.go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"stock-bot/internal/infrastructure/client/dto/order/request"
	"stock-bot/internal/infrastructure/client/dto/order/response"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type orderClientImpl struct {
	client *TachibanaClient
	logger *zap.Logger
}

func (o *orderClientImpl) NewOrder(ctx context.Context, req request.ReqNewOrder) (*response.ResNewOrder, error) {
	if !o.client.loggined {
		return nil, errors.New("not logged in")
	}

	// 1. リクエストURLの作成
	u, err := url.Parse(o.client.loginInfo.RequestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}

	// 2. リクエストパラメータの作成
	req.CLMID = "CLMKabuNewOrder"
	req.P_no = o.client.getPNo()
	req.P_sd_date = formatSDDate(time.Now())
	req.SJsonOfmt = "4"

	// 構造体を map[string]string に変換
	params, err := structToMapString(req) //utilの関数
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
	respMap, err := SendRequest(httpReq, 3, o.logger)

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
	req.SJsonOfmt = "4"                      // JSON出力フォーマット

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
	respMap, err := SendRequest(httpReq, 3, o.logger)
	if err != nil {
		return nil, errors.Wrap(err, "correct order failed")
	}

	// 5. レスポンスの処理
	res, err := ConvertResponse[response.ResCorrectOrder](respMap)
	if err != nil {
		return nil, err
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
	req.SJsonOfmt = "4"

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

	respMap, err := SendRequest(httpReq, 3, o.logger)
	if err != nil {
		return nil, errors.Wrap(err, "cancel order failed")
	}

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
	req.SJsonOfmt = "4"

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

	respMap, err := SendRequest(httpReq, 3, o.logger)
	if err != nil {
		return nil, errors.Wrap(err, "cancel all order failed")
	}

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
	req.SJsonOfmt = "4"

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

	respMap, err := SendRequest(httpReq, 3, o.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get order list failed")
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
	req.SJsonOfmt = "4"

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
	respMap, err := SendRequest(httpReq, 3, o.logger)
	if err != nil {
		return nil, errors.Wrap(err, "get order list detail failed")
	}

	res, err := ConvertResponse[response.ResOrderListDetail](respMap)
	if err != nil {
		return nil, err
	}

	// ★★★ ここで型アサーションを使って aKessaiOrderTategyokuList にアクセス ★★★
	if kessaiList, ok := res.KessaiOrderTategyokuList.([]interface{}); ok {
		// kessaiList は []interface{} 型 (要素の型は不明)
		for _, kessai := range kessaiList {
			if _, ok := kessai.(map[string]interface{}); ok {
				// kessaiMap は map[string]interface{} 型
				// (例) kessaiMap["sKessaiTategyokuDay"] などで値にアクセスできる
				//      (ただし、値の型は interface{} なので、型アサーションが必要)

				//fmt.Printf("建日: %v\n", kessaiMap["sKessaiTategyokuDay"]) // 例 (型アサーションが必要)
			}
		}
	} else if kessaiList, ok := res.KessaiOrderTategyokuList.([]response.ResKessaiOrderTategyoku); ok {
		for _, kessai := range kessaiList {
			// kessai は response.ResKessaiOrderTategyoku 型
			fmt.Printf("建日: %s, 建単価: %s\n", kessai.KessaiTategyokuDay, kessai.KessaiTategyokuPrice)
		}

	} else {
		//  []interface{}でも、[]response.ResKessaiOrderTategyokuでもない場合
		fmt.Println("res.KessaiOrderTategyokuList is nil")
	}

	return res, nil
}
