package model

import "gorm.io/gorm"

type IndexedModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
