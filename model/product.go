package model

import guuid "github.com/google/uuid"

type Product struct {
	ProductID guuid.UUID `gorm:"primaryKey" json:"productid"`
	UserRefer guuid.UUID `json:"-"`
	Value     string     `json:"value"`
	Name      string     `json:"name"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
}
