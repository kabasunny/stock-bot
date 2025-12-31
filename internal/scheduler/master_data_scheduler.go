package scheduler

import (
	"context"
	"log/slog"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
	"time"
)

// MasterDataScheduler はマスターデータの定期更新を管理する
type MasterDataScheduler struct {
	masterUseCase app.MasterUseCase
	session       *client.Session
	logger        *slog.Logger
	ticker        *time.Ticker
	stopCh        chan struct{}
}

// NewMasterDataScheduler は新しいスケジューラーを作成する
func NewMasterDataScheduler(
	masterUseCase app.MasterUseCase,
	session *client.Session,
	logger *slog.Logger,
) *MasterDataScheduler {
	return &MasterDataScheduler{
		masterUseCase: masterUseCase,
		session:       session,
		logger:        logger,
		stopCh:        make(chan struct{}),
	}
}

// Start はスケジューラーを開始する（日次午前2時に実行）
func (s *MasterDataScheduler) Start() {
	s.logger.Info("Starting master data scheduler")

	// 次の午前2時までの時間を計算
	now := time.Now()
	next2AM := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(next2AM) {
		next2AM = next2AM.Add(24 * time.Hour)
	}

	// 初回実行までの待機時間
	initialDelay := next2AM.Sub(now)
	s.logger.Info("Master data scheduler will start at", "next_execution", next2AM)

	go func() {
		// 初回実行まで待機
		select {
		case <-time.After(initialDelay):
			s.updateMasterData()
		case <-s.stopCh:
			return
		}

		// 24時間間隔で定期実行
		s.ticker = time.NewTicker(24 * time.Hour)
		defer s.ticker.Stop()

		for {
			select {
			case <-s.ticker.C:
				s.updateMasterData()
			case <-s.stopCh:
				s.logger.Info("Master data scheduler stopped")
				return
			}
		}
	}()
}

// Stop はスケジューラーを停止する
func (s *MasterDataScheduler) Stop() {
	s.logger.Info("Stopping master data scheduler")
	close(s.stopCh)
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

// updateMasterData はマスターデータを更新する
func (s *MasterDataScheduler) updateMasterData() {
	s.logger.Info("Starting scheduled master data update")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := s.masterUseCase.DownloadAndStoreMasterData(ctx, s.session); err != nil {
		s.logger.Error("Failed to update master data", "error", err)
		return
	}

	s.logger.Info("Scheduled master data update completed successfully")
}

// TriggerManualUpdate は手動更新をトリガーする
func (s *MasterDataScheduler) TriggerManualUpdate(ctx context.Context) error {
	s.logger.Info("Manual master data update triggered")
	return s.masterUseCase.DownloadAndStoreMasterData(ctx, s.session)
}
