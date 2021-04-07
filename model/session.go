package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	Expires   time.Time  `json:"-"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
}
