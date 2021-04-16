package model

import guuid "github.com/google/uuid"

type Product struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
	UserRefer guuid.UUID `json:"-"`
	Value     string     `json:"value"`
	Name      string     `json:"name"`
}
