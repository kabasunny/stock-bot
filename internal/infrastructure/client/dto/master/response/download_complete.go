// internal/infrastructure/client/dto/master/response/download_complete.go
package response

// ResDownloadComplete は、初期ダウンロード終了通知を表すDTO
type ResDownloadComplete struct {
	CLMID string `json:"sCLMID"` // 機能ID (CLMEventDownloadComplete)
}
