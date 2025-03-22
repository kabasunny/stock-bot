// internal/infrastructure/client/dto/master/response/download_complete.go
package response

// DownloadComplete は、初期ダウンロード終了通知を表すDTOです。
type DownloadComplete struct {
	CLMID string `json:"sCLMID"` // 機能ID (CLMEventDownloadComplete)
}
