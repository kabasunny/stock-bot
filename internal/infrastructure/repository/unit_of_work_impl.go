package repository

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/event"
	"stock-bot/domain/repository"

	"gorm.io/gorm"
)

// UnitOfWorkImpl はUnitOfWorkの実装
type UnitOfWorkImpl struct {
	db             *gorm.DB
	tx             *gorm.DB
	eventPublisher event.EventPublisher
	logger         *slog.Logger
	domainEvents   []event.DomainEvent
	inTransaction  bool

	// Repository instances
	orderRepo    repository.OrderRepository
	positionRepo repository.PositionRepository
	masterRepo   repository.MasterRepository
	strategyRepo repository.StrategyRepository
}

// NewUnitOfWork は新しいUnitOfWorkを作成
func NewUnitOfWork(
	db *gorm.DB,
	eventPublisher event.EventPublisher,
	logger *slog.Logger,
) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		db:             db,
		eventPublisher: eventPublisher,
		logger:         logger,
		domainEvents:   make([]event.DomainEvent, 0),
		inTransaction:  false,
	}
}

// Begin はトランザクションを開始
func (uow *UnitOfWorkImpl) Begin(ctx context.Context) error {
	if uow.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}

	uow.tx = uow.db.Begin()
	if uow.tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", uow.tx.Error)
	}

	uow.inTransaction = true
	uow.domainEvents = make([]event.DomainEvent, 0)

	// トランザクション用のリポジトリインスタンスを作成
	uow.orderRepo = NewOrderRepository(uow.tx)
	uow.positionRepo = NewPositionRepository(uow.tx, uow.orderRepo)
	uow.masterRepo = NewMasterRepository(uow.tx)
	// uow.strategyRepo = NewStrategyRepository(uow.tx) // 実装後に有効化

	uow.logger.Debug("transaction started")
	return nil
}

// Commit はトランザクションをコミットし、ドメインイベントを発行
func (uow *UnitOfWorkImpl) Commit(ctx context.Context) error {
	if !uow.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	// トランザクションをコミット
	if err := uow.tx.Commit().Error; err != nil {
		uow.logger.Error("failed to commit transaction", slog.Any("error", err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	uow.logger.Debug("transaction committed", slog.Int("domain_events", len(uow.domainEvents)))

	// ドメインイベントを発行
	if err := uow.publishDomainEvents(ctx); err != nil {
		uow.logger.Error("failed to publish domain events after commit", slog.Any("error", err))
		// トランザクションは既にコミット済みなので、イベント発行エラーは警告として扱う
	}

	// 状態をリセット
	uow.inTransaction = false
	uow.tx = nil
	uow.domainEvents = make([]event.DomainEvent, 0)
	uow.orderRepo = nil
	uow.positionRepo = nil
	uow.masterRepo = nil

	return nil
}

// Rollback はトランザクションをロールバック
func (uow *UnitOfWorkImpl) Rollback(ctx context.Context) error {
	if !uow.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if err := uow.tx.Rollback().Error; err != nil {
		uow.logger.Error("failed to rollback transaction", slog.Any("error", err))
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	uow.logger.Debug("transaction rolled back", slog.Int("discarded_events", len(uow.domainEvents)))

	// 状態をリセット
	uow.inTransaction = false
	uow.tx = nil
	uow.domainEvents = make([]event.DomainEvent, 0)
	uow.orderRepo = nil
	uow.positionRepo = nil
	uow.masterRepo = nil

	return nil
}

// AddDomainEvent はドメインイベントを追加
func (uow *UnitOfWorkImpl) AddDomainEvent(domainEvent event.DomainEvent) {
	uow.domainEvents = append(uow.domainEvents, domainEvent)
	uow.logger.Debug("domain event added",
		slog.String("event_type", domainEvent.EventType()),
		slog.String("event_id", domainEvent.EventID()),
		slog.String("aggregate_id", domainEvent.AggregateID()))
}

// GetDomainEvents は蓄積されたドメインイベントを取得
func (uow *UnitOfWorkImpl) GetDomainEvents() []event.DomainEvent {
	return uow.domainEvents
}

// ClearDomainEvents はドメインイベントをクリア
func (uow *UnitOfWorkImpl) ClearDomainEvents() {
	uow.domainEvents = make([]event.DomainEvent, 0)
}

// IsInTransaction はトランザクション中かどうかを判定
func (uow *UnitOfWorkImpl) IsInTransaction() bool {
	return uow.inTransaction
}

// Repository accessors

func (uow *UnitOfWorkImpl) OrderRepository() repository.OrderRepository {
	if uow.inTransaction {
		return uow.orderRepo
	}
	// トランザクション外では通常のリポジトリを返す
	return NewOrderRepository(uow.db)
}

func (uow *UnitOfWorkImpl) PositionRepository() repository.PositionRepository {
	if uow.inTransaction {
		return uow.positionRepo
	}
	// トランザクション外では通常のリポジトリを返す
	orderRepo := NewOrderRepository(uow.db)
	return NewPositionRepository(uow.db, orderRepo)
}

func (uow *UnitOfWorkImpl) MasterRepository() repository.MasterRepository {
	if uow.inTransaction {
		return uow.masterRepo
	}
	// トランザクション外では通常のリポジトリを返す
	return NewMasterRepository(uow.db)
}

// publishDomainEvents はドメインイベントを発行
func (uow *UnitOfWorkImpl) publishDomainEvents(ctx context.Context) error {
	if uow.eventPublisher == nil {
		uow.logger.Warn("no event publisher configured, skipping domain event publication")
		return nil
	}

	for _, domainEvent := range uow.domainEvents {
		if err := uow.eventPublisher.Publish(ctx, domainEvent); err != nil {
			return fmt.Errorf("failed to publish domain event %s: %w", domainEvent.EventID(), err)
		}
	}

	return nil
}

// WithTransaction はトランザクション内で処理を実行するヘルパー関数
func (uow *UnitOfWorkImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			uow.Rollback(ctx)
			panic(r)
		}
	}()

	if err := fn(ctx); err != nil {
		if rollbackErr := uow.Rollback(ctx); rollbackErr != nil {
			uow.logger.Error("failed to rollback after error",
				slog.Any("original_error", err),
				slog.Any("rollback_error", rollbackErr))
		}
		return err
	}

	return uow.Commit(ctx)
}
func (uow *UnitOfWorkImpl) StrategyRepository() repository.StrategyRepository {
	if uow.inTransaction {
		return uow.strategyRepo
	}
	// トランザクション外では通常のリポジトリを返す（実装後に有効化）
	// return NewStrategyRepository(uow.db)
	return nil // 暫定実装
}
