// internal/infrastructure/client/dto/master/request/download_master.go
package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// DownloadMasterRequest は、マスタ情報ダウンロードの要求を表すDTO
type DownloadMasterRequest struct {
	dto.RequestBase        // 共通フィールドを埋め込む
	CLMID           string `json:"sCLMID"`                 // 機能ID (固定値: "CLMEventDownload")
	TargetCLMID     string `json:"sTargetCLMID,omitempty"` // 対象機能ID (カンマ区切りで複数指定可能、空文字列の場合は全マスタ情報)
}
