package api

import (
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// User Auth Model
type User struct {
	ID        guuid.UUID `gorm:"primaryKey" json:"-"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Sessions  []Session  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	Products  []Product  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"products"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}

// Session Model for the user
type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
}

// General Purpose Product Model
type Product struct {
	ProductID guuid.UUID `gorm:"primaryKey" json:"productid"`
	UserRefer guuid.UUID `json:"-"`
	Value     string     `json:"value"`
	Name      string     `json:"name"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
}

type ProductRequest struct {
	Product
	Sessionid guuid.UUID
}

type ChangePassword struct {
	User
	NewPassword string `json:"sessionid"`
}

func Initalize(router fiber.Router, db *gorm.DB) {
	db.AutoMigrate(&User{}, &Session{}, &Product{})
	AuthRoutes(router, db)
	ProductRoutes(router, db)
}

func GetUser(sessionid guuid.UUID, db *gorm.DB) (User, int) {
	query := Session{Sessionid: sessionid}
	found := Session{}
	err := db.First(&found, &query).Error
	if err != nil {
		return User{}, 401
	}
	user := User{}
	usrQuery := User{ID: found.UserRefer}
	err = db.First(&user, &usrQuery).Error
	if err != nil {
		return User{}, 401
	}
	return user, 0
}

func SecurityMiddleware(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.Next()
}
