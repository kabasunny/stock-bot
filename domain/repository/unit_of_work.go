package repository

import (
	"context"
	"stock-bot/domain/event"
)

// UnitOfWork はトランザクション境界とドメインイベントの管理を行う
type UnitOfWork interface {
	// Begin はトランザクションを開始
	Begin(ctx context.Context) error

	// Commit はトランザクションをコミットし、ドメインイベントを発行
	Commit(ctx context.Context) error

	// Rollback はトランザクションをロールバック
	Rollback(ctx context.Context) error

	// WithTransaction はトランザクション内で処理を実行するヘルパー関数
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// AddDomainEvent はドメインイベントを追加
	AddDomainEvent(event event.DomainEvent)

	// GetDomainEvents は蓄積されたドメインイベントを取得
	GetDomainEvents() []event.DomainEvent

	// ClearDomainEvents はドメインイベントをクリア
	ClearDomainEvents()

	// IsInTransaction はトランザクション中かどうかを判定
	IsInTransaction() bool

	// Repository accessors
	OrderRepository() OrderRepository
	PositionRepository() PositionRepository
	MasterRepository() MasterRepository
	StrategyRepository() StrategyRepository
}

// AggregateRoot は集約ルートのインターフェース
type AggregateRoot interface {
	// GetDomainEvents は集約が発行するドメインイベントを取得
	GetDomainEvents() []event.DomainEvent

	// ClearDomainEvents はドメインイベントをクリア
	ClearDomainEvents()

	// AddDomainEvent はドメインイベントを追加
	AddDomainEvent(event event.DomainEvent)
}
