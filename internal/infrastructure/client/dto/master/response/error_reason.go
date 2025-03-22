// internal/infrastructure/client/dto/master/response/error_reason.go
package response

// ErrorReason は、取引所エラー等理由コードの情報を表すDTOです。
type ErrorReason struct {
	CLMID     string `json:"sCLMID"`         // 機能ID (CLMOrderErrReason)
	ErrorCode string `json:"sErrReasonCode"` // 取引所エラー等理由コード
	ErrorText string `json:"sErrReasonText"` // 取引所エラー等理由テキスト
}
