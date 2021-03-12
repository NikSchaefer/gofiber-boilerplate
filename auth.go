package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// User Auth Model
type User struct {
	gorm.Model
	ID       guuid.UUID `gorm:"primaryKey" json:"id"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Sessions []Session  `gorm:"foreignKey:UserRefer" gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Session Model for the user
type Session struct {
	gorm.Model
	Sessionkey guuid.UUID `gorm:"primaryKey"`
	UserRefer  guuid.UUID
}

// Initalize and set the authentication and authorization routes
func AuthRoutes(router fiber.Router, db *gorm.DB) {
	auth := router.Group("/auth", securityMiddleware)
	auth.Post("/login", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}

		foundUser := User{}
		queryUser := User{Username: json.Username}
		err := db.First(&foundUser, &queryUser).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User not Found")
		}
		if foundUser.Password != json.Password {
			return c.Status(401).SendString("Incorrect Password")
		}
		newSession := Session{UserRefer: foundUser.ID, Sessionkey: guuid.New()}
		CreateErr := db.Create(&newSession).Error
		if CreateErr != nil {
			return c.Status(500).SendString("Creation Error")
		}
		return c.Status(200).JSON(newSession)
	})

	auth.Post("/logout", func(c *fiber.Ctx) error {
		json := new(Session)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := Session{}
		if json.Sessionkey == empty.Sessionkey {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		return c.SendStatus(200)
	})
	auth.Post("/create", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		newUser := User{
			Username: json.Username,
			Password: json.Password,
			ID:       guuid.New(),
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err != gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Already Exists")
		}
		db.Create(&newUser)
		return c.SendStatus(200)
	})
	auth.Post("/user", func(c *fiber.Ctx) error {
		user := User{}
		myUser := User{Username: "NikSchaefer"}
		Sessions := []Session{}
		db.First(&user, &myUser)
		db.Model(&user).Association("Sessions").Find(&Sessions)
		user.Sessions = Sessions
		return c.JSON(user)
	})
	auth.Post("/delete", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		if json.Password != foundUser.Password {
			return c.Status(401).SendString("Invalid Credentials")
		}
		db.Model(&foundUser).Association("Sessions").Clear()
		createErr := db.Delete(&foundUser).Error
		if createErr != nil {
			fmt.Println(createErr)
		}
		return c.SendStatus(200)
	})
	auth.Post("/update", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		return c.SendStatus(200)
	})

}

func securityMiddleware(c *fiber.Ctx) error {
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Download-Options", "noopen")
	c.Set("Strict-Transport-Security", "max-age=5184000")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Accepts("application/json")
	return c.Next()
}
