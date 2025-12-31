package event

import (
	"time"
)

// DomainEvent はドメインイベントの基底インターフェース
type DomainEvent interface {
	// EventID はイベントの一意識別子を返す
	EventID() string
	// EventType はイベントの種別を返す
	EventType() string
	// OccurredAt はイベントの発生時刻を返す
	OccurredAt() time.Time
	// AggregateID は関連する集約のIDを返す
	AggregateID() string
	// Version はイベントのバージョンを返す
	Version() int
}

// BaseDomainEvent はドメインイベントの基底実装
type BaseDomainEvent struct {
	eventID     string
	eventType   string
	occurredAt  time.Time
	aggregateID string
	version     int
}

// NewBaseDomainEvent は新しい基底ドメインイベントを作成
func NewBaseDomainEvent(eventType, aggregateID string, version int) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     generateEventID(),
		eventType:   eventType,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
		version:     version,
	}
}

func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseDomainEvent) Version() int {
	return e.version
}

// generateEventID はイベントIDを生成
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString はランダム文字列を生成
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
