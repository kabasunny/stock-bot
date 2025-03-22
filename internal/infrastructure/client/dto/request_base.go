// internal/infrastructure/client/dto/request_base.go
package dto

// RequestBase は、APIリクエストDTOの共通フィールドを定義する構造体です。
type RequestBase struct {
	P_no      string `json:"p_no"`      // p_no (連番、クライアント側で設定)
	P_sd_date string `json:"p_sd_date"` // システム日付 (YYYYMMDD形式、クライアント側で設定)
	SJsonOfmt string `json:"sJsonOfmt"` // JSON出力フォーマット (固定値: "4")
}
