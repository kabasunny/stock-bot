// internal/infrastructure/client/util.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "stock-bot/internal/logger"

	"github.com/cockroachdb/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// SendPostRequest は、HTTP POSTリクエストを送信し、レスポンスをデコードする
func SendPostRequest(
	ctx context.Context,
	client *http.Client, // Add http.Client parameter
	urlStr string,
	reqData interface{},
	maxRetries int,
) (map[string]interface{}, error) {
	// リクエストデータをJSONにマーシャリング
	jsonBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request data")
	}

	// POSTリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create post request")
	}

	// リトライ時にボディを再読み込みできるように、GetBody を設定
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(jsonBody)), nil
	}

	// Content-Typeヘッダーを設定
	req.Header.Set("Content-Type", "application/json")

	// 既存のSendRequestを呼び出してリクエストを送信
	return SendRequest(client, req, maxRetries) // Pass client
}

// SendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func SendRequest(
	client *http.Client, // Add http.Client parameter
	req *http.Request,
	maxRetries int,
) (map[string]interface{}, error) { // logger引数を削除
	var response map[string]interface{}

	// retryDoに渡す関数
	retryFunc := func(client *http.Client, decodeFunc func([]byte, interface{}) error) (*http.Response, error) {
		// リトライのたびにボディを再生成する
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get request body for retry")
			}
			req.Body = body
		}

		//timeoutコンテキストを作成
		reqWithTimeout, cancel := withContextAndTimeout(req, 60*time.Second)
		defer cancel()

		resp, err := client.Do(reqWithTimeout) //clientは、http.Client{}
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
		logRequestAndResponse(req, body) // loggerを渡さない

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

	// retryDo を呼び出す際に、retryFunc と decodeFunc を渡す (変更なし)
	resp, err := retryDo(retryFunc, maxRetries, 2*time.Second, client, decodeFunc) // Use passed client
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return response, nil
}

// ConvertResponse は、map[string]interface{} をレスポンスDTOに変換する (再帰対応版)
func ConvertResponse[T any](respMap map[string]interface{}) (*T, error) {
	var res T
	if err := convertMapToStruct(respMap, &res, ""); err != nil { // key="" で初期呼び出し
		return nil, err
	}
	return &res, nil
}

// convertMapToStruct は、map[string]interface{} を構造体に変換する再帰関数
func convertMapToStruct(srcMap map[string]interface{}, dest interface{}, key string) error { // key を引数に追加
	destVal := reflect.ValueOf(dest)

	// ポインタでなければエラー
	if destVal.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	// nilポインタの場合もエラー
	if destVal.IsNil() {
		return errors.New("dest is nil pointer")
	}

	destVal = destVal.Elem() // ポインタの指す先の値を取得

	switch destVal.Kind() {
	case reflect.Struct:
		for i := 0; i < destVal.NumField(); i++ {
			field := destVal.Field(i)
			fieldType := destVal.Type().Field(i)
			currentKey := fieldType.Name // フィールド名をキーとして使用

			// jsonタグがあればそちらを優先
			if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				currentKey = jsonTag
			}

			if srcValue, ok := srcMap[currentKey]; ok {
				// srcValueをfieldの型に合うように変換してセット
				var err error
				if field.Kind() == reflect.Slice {
					// スライスの場合はconvertMapToStructをcurrentKeyを渡して呼び出し
					err = convertMapToStruct(srcMap, field.Addr().Interface(), currentKey)
				} else {
					// スライスでない場合は、今まで通りsetFieldValue
					err = setFieldValue(field, srcValue)
				}
				if err != nil {
					return errors.Wrapf(err, "failed to set field %s", currentKey)
				}
			}
		}
	case reflect.Map:
		// map[string]interface{} と map[string]string のみを処理。それ以外はエラー
		if destVal.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("unsupported map key type: %v", destVal.Type().Key().Kind())
		}

		// reflect.MapOf でmapの型を作成。interface{}を要素とするmap
		mapType := reflect.MapOf(destVal.Type().Key(), reflect.TypeOf((*interface{})(nil)).Elem())

		// reflect.MakeMap で指定された型のmapを作成
		newMap := reflect.MakeMap(mapType)

		for k, v := range srcMap {
			// 再帰的にmapのvalueを適切な型に変換
			val := reflect.New(reflect.TypeOf((*interface{})(nil)).Elem())

			// valのポインタを渡す
			if err := convertMapToStruct(map[string]interface{}{"tempKey": v}, val.Interface(), k); err != nil { // keyを渡す
				return errors.Wrapf(err, "failed to convert map value for key %s", k)
			}
			newMap.SetMapIndex(reflect.ValueOf(k), val.Elem().FieldByName("TempKey"))
		}
		destVal.Set(newMap)

	case reflect.Slice: // スライスの処理
		elemType := destVal.Type().Elem()
		srcSlice, ok := srcMap[key].([]interface{}) // keyを使ってスライスを取得
		if !ok {
			// スライスが見つからない場合 (空文字列 "" の場合など) は、
			// 空のスライスをセットして処理を続行
			emptySlice := reflect.MakeSlice(destVal.Type(), 0, 0)
			destVal.Set(emptySlice)
			return nil //早期return
		}

		newSlice := reflect.MakeSlice(destVal.Type(), len(srcSlice), len(srcSlice))
		for i, elem := range srcSlice {
			elemVal := reflect.New(elemType)
			if elemMap, ok := elem.(map[string]interface{}); ok {
				if err := convertMapToStruct(elemMap, elemVal.Interface(), ""); err != nil { // keyはここでは使用しない
					return err
				}
				newSlice.Index(i).Set(elemVal.Elem())
			} else {
				// map[string]interface{} に変換できない要素があった場合
				return fmt.Errorf("slice element is not a map: %v", reflect.TypeOf(elem))
			}
		}
		destVal.Set(newSlice)

	default: // その他の型 (ポインタ、interface{}など)
		return fmt.Errorf("unsupported type: %v", destVal.Type())
	}

	return nil
}

// setFieldValue は、reflect.Value に値をセットする (型変換付き)
func setFieldValue(field reflect.Value, value interface{}) error {
	val := reflect.ValueOf(value)

	switch field.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		} else {
			field.SetString(fmt.Sprintf("%v", value)) // 文字列に変換
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intValue int64
		switch v := value.(type) {
		case int:
			intValue = int64(v)
		case int8:
			intValue = int64(v)
		case int16:
			intValue = int64(v)
		case int32:
			intValue = int64(v)
		case int64:
			intValue = v
		case float32:
			intValue = int64(v)
		case float64:
			intValue = int64(v)
		case string:
			// 文字列からの変換も試みる
			_, err := fmt.Sscan(v, &intValue)
			if err != nil {
				return errors.Wrap(err, "failed to convert string to int")
			}
		default:
			return fmt.Errorf("unsupported type conversion to int: %v", reflect.TypeOf(value))
		}
		field.SetInt(intValue)

	case reflect.Float32, reflect.Float64:
		var floatValue float64
		switch v := value.(type) {
		case float32:
			floatValue = float64(v)
		case float64:
			floatValue = v
		case string: // 文字列からの変換
			_, err := fmt.Sscan(v, &floatValue)
			if err != nil {
				return errors.Wrap(err, "failed to convert string to float")
			}
		default:
			return fmt.Errorf("unsupported type conversion to float: %v", reflect.TypeOf(value))
		}
		field.SetFloat(floatValue)

	case reflect.Bool:
		var boolValue bool
		switch v := value.(type) {
		case bool:
			boolValue = v
		case string: // 文字列からの変換
			_, err := fmt.Sscan(v, &boolValue)
			if err != nil {
				return errors.Wrap(err, "failed to convert string to bool")
			}
		default:
			return fmt.Errorf("unsupported type conversion to bool: %v", reflect.TypeOf(value))
		}
		field.SetBool(boolValue)
	case reflect.Ptr:
		// ポインタの指す先の型を取得
		elemType := field.Type().Elem()
		// 新しいポインタを作成
		newPtr := reflect.New(elemType)
		// 再帰的に処理. "tempKey"というキーにバリューをいれて渡す
		if err := convertMapToStruct(map[string]interface{}{"tempKey": value}, newPtr.Interface(), "tempKey"); err != nil {
			// Keyがないことによるエラーは無視
			if err.Error() != "failed to set field tempKey: unsupported type: <nil>" {
				return errors.Wrap(err, "failed to convert pointer value")
			}
		}
		field.Set(newPtr)

	case reflect.Struct:
		// ネストされた構造体の場合、再帰的に処理
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if err := convertMapToStruct(nestedMap, field.Addr().Interface(), ""); err != nil { //keyは使用しない
				return err
			}
		} else {
			return fmt.Errorf("expected map for struct, got: %v", reflect.TypeOf(value))
		}

	case reflect.Interface:
		// インターフェースの場合、型アサーションして値を直接セット (型情報は失われる)
		field.Set(val)

	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}

	return nil
}

// structToMapString は、構造体を map[string]string に変換する (json.Marshal を使わない)
func structToMapString(data interface{}) (map[string]string, error) {
	params := make(map[string]string)
	if err := marshalToMap("", data, params); err != nil { // marshalToMapを呼び出す
		return nil, err
	}
	return params, nil
}

// marshalToMap は、data を map[string]string に変換する再帰関数
func marshalToMap(prefix string, data interface{}, params map[string]string) error {
	val := reflect.ValueOf(data)

	// ポインタの場合は指す先の値を取得
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil // nil の場合は何もしない
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := val.Type().Field(i)
			key := fieldType.Name

			// json タグの処理
			if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				if commaIndex := strings.Index(jsonTag, ","); commaIndex > 0 {
					key = jsonTag[:commaIndex] // ,omitemptyなどを考慮する場合
				} else {
					key = jsonTag
				}
			}

			// 埋め込みフィールドの場合はプレフィックスを付けない
			if fieldType.Anonymous {
				if err := marshalToMap(prefix, field.Interface(), params); err != nil {
					return err
				}
				continue
			}

			// フィールドにプレフィックスを付ける
			if prefix != "" {
				key = prefix + "." + key
			}

			// 再帰呼び出し or 値を文字列として追加
			switch field.Kind() {
			case reflect.Slice, reflect.Struct, reflect.Interface: // 再帰が必要な型
				if err := marshalToMap(key, field.Interface(), params); err != nil { //prefixをつけて再帰
					return err
				}
			default: // string, 数値, boolなど
				params[key] = fmt.Sprintf("%v", field.Interface()) // 文字列化して追加
			}
		}
	case reflect.Slice: //スライスの場合
		if prefix == "" { // トップレベルのスライスは非対応(structのフィールドである前提)
			return fmt.Errorf("unsupported top-level slice")
		}
		var sliceValues []string
		for i := 0; i < val.Len(); i++ {
			elemVal := val.Index(i)

			// スライスの各要素を JSON に変換
			elemBuf := new(bytes.Buffer)
			if err := marshalValue(elemBuf, elemVal); err != nil {
				return errors.Wrapf(err, "failed to marshal slice element at index %d", i)
			}

			sliceValues = append(sliceValues, elemBuf.String())
		}

		// JSON 配列形式の文字列を生成
		params[prefix] = "[" + strings.Join(sliceValues, ",") + "]"

	case reflect.Interface: // インターフェースの場合
		if val.IsNil() {
			return nil // nil の場合は何もしない
		}
		// 中身を取り出して再帰呼び出し
		return marshalToMap(prefix, val.Elem().Interface(), params)

	default: // string, 数値, boolなど
		if prefix == "" { // structのフィールドである必要あり
			return fmt.Errorf("unsupported top-level type: %v", val.Type())
		}
		params[prefix] = fmt.Sprintf("%v", val.Interface()) //文字列か
	}

	return nil
}

// marshalValue は、単一の値を JSON 形式の文字列に変換し、buf に書き込む (marshalJSON からのヘルパー関数)
func marshalValue(buf *bytes.Buffer, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return marshalValue(buf, val.Elem())
	case reflect.Struct:
		return marshalStruct(buf, val)
	case reflect.Slice, reflect.Array:
		return marshalSlice(buf, val)
	case reflect.String:
		buf.WriteString(strconv.Quote(val.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", val.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", val.Float())
	case reflect.Bool:
		fmt.Fprintf(buf, "%t", val.Bool())
	case reflect.Interface:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return marshalValue(buf, val.Elem())
	default:
		return fmt.Errorf("unsupported type: %v", val.Kind())
	}
	return nil
}

func marshalStruct(buf *bytes.Buffer, val reflect.Value) error {
	buf.WriteString("{")
	first := true
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		key := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIndex := strings.Index(jsonTag, ","); commaIndex > 0 {
				key = jsonTag[:commaIndex]
			} else {
				key = jsonTag
			}
		}

		// 埋め込みフィールドの場合は、prefixをつけない
		if fieldType.Anonymous {
			if err := marshalValue(buf, field); err != nil { // 再帰呼び出し
				return err
			}
			continue //以降の処理はスキップ
		}

		if !first {
			buf.WriteString(",")
		}
		first = false
		buf.WriteString(strconv.Quote(key))
		buf.WriteString(":")
		if err := marshalValue(buf, field); err != nil {
			return err
		}
	}
	buf.WriteString("}")
	return nil
}

func marshalSlice(buf *bytes.Buffer, val reflect.Value) error {
	buf.WriteString("[")
	for i := 0; i < val.Len(); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		if err := marshalValue(buf, val.Index(i)); err != nil {
			return err
		}
	}
	buf.WriteString("]")
	return nil
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
			slog.Debug("Retrying request",
				slog.Int("retry_count", i+1),
				slog.Any("error", err),
			)
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
func logRequestAndResponse(req *http.Request, respBody []byte) {
	// リクエスト情報のログ出力
	slog.Debug("Request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Any("headers", req.Header),
	)

	// req.URL をデコードして表示
	decodedURL, _ := url.QueryUnescape(req.URL.String())
	slog.Debug("Decoded URL", slog.String("decodedUrl", decodedURL))

	// UTF-8に変換を試みる（Shift_JISデコード）
	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), respBody)
	if err != nil {
		slog.Error("Failed to decode response body to UTF-8", slog.Any("error", err))
		return
	}

	// JSONを整形して出力
	var prettyBody bytes.Buffer
	err = json.Indent(&prettyBody, bodyUTF8, "", "    ")
	if err != nil {
		slog.Error("Failed to format JSON", slog.Any("error", err))
		slog.Debug("Raw Response Body (UTF-8)", slog.String("body", string(bodyUTF8)))
		return
	}

	slog.Debug("Formatted Response Body (UTF-8)", slog.String("body", prettyBody.String()))
}

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
