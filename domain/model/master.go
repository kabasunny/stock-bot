// domain/model/master.go
package model

import "time"

// MasterBase は、すべてのマスタデータに共通するフィールドを定義します。
type MasterBase struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
