package api

import (
	"time"

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
	Expires   time.Time  `json:"-"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
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

// Set Routes and Migrations for Models
func Initalize(router fiber.Router, db *gorm.DB) {
	db.AutoMigrate(&User{}, &Session{}, &Product{})
	AuthRoutes(router, db)
	ProductRoutes(router, db)
}

func GetUser(sessionid guuid.UUID, db *gorm.DB) (User, int) {
	query := Session{Sessionid: sessionid}
	found := Session{}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, fiber.StatusNotFound
	}
	user := User{}
	usrQuery := User{ID: found.UserRefer}
	err = db.First(&user, &usrQuery).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, fiber.StatusNotFound
	}
	return user, 0
}

func JsonMiddleware(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.Next()
}

// func AuthenticatedMiddleware(c *fiber.Ctx) error {
// 	json := new(Session)
// 	if err := c.BodyParser(json); err != nil {
// 		return c.SendStatus(fiber.StatusBadRequest)
// 	}
// 	user, status := GetUser(json.Sessionid, database)
// 	if status != 0 {
// 		return c.SendStatus(status)
// 	}
// 	c.Locals("user", user)
// 	return c.Next()
// }
