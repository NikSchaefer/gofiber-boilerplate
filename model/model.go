package model

type Model struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	CreatedAt int64 `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"-"`
}
