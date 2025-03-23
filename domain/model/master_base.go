// domain/model/master.go
package model

import "time"

// MasterBase は、すべてのマスタデータに共通するフィールドを定義
// 共通項目 (sCreateTime, sUpdateTime, sUpdateNumber, sDeleteFlag, sDeleteTime) に対応し
type MasterBase struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
