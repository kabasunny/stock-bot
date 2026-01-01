package repository

import (
	"encoding/json"
)

// marshalJSON はオブジェクトをJSON文字列に変換
func marshalJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// unmarshalJSON はJSON文字列をオブジェクトに変換
func unmarshalJSON(data string, v interface{}) error {
	if data == "" {
		return nil
	}
	return json.Unmarshal([]byte(data), v)
}
