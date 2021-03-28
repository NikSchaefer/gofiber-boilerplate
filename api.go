package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Auth Model
type User struct {
	ID        guuid.UUID `gorm:"primaryKey" json:"-"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Sessions  []Session  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"sessions"`
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
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
}

type ChangePassword struct {
	User
	NewPassword string
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
		if !comparePasswords(foundUser.Password, []byte(json.Password)) {
			return c.Status(401).SendString("Incorrect Password")
		}
		newSession := Session{UserRefer: foundUser.ID, Sessionid: guuid.New()}
		CreateErr := db.Create(&newSession).Error
		if CreateErr != nil {
			return c.Status(500).SendString("Creation Error")
		}
		c.Cookie(&fiber.Cookie{
			Name:     "sessionid",
			Expires:  time.Now().Add(5 * 24 * time.Hour),
			Value:    newSession.Sessionid.String(),
			HTTPOnly: true,
		})
		return c.Status(200).JSON(newSession)
	})

	auth.Post("/logout", func(c *fiber.Ctx) error {
		json := new(Session)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		if json.Sessionid == new(Session).Sessionid {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		session := Session{}
		query := Session{Sessionid: json.Sessionid}
		err := db.First(&session, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("Session Not Found")
		}
		db.Delete(&session)
		c.ClearCookie("sessionid")
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
		pw := hashAndSalt([]byte(json.Password))
		newUser := User{
			Username: json.Username,
			Password: pw,
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
			return c.SendStatus(400)
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
		if !comparePasswords(foundUser.Password, []byte(json.Password)) {
			return c.Status(401).SendString("Invalid Credentials")
		}
		db.Model(&foundUser).Association("Sessions").Clear()
		createErr := db.Delete(&foundUser).Error
		if createErr != nil {
			fmt.Println(createErr)
		}
		c.ClearCookie("sessionid")
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
	auth.Post("/changepassword", func(c *fiber.Ctx) error {
		json := new(ChangePassword)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		if !comparePasswords(foundUser.Password, []byte(json.NewPassword)) {
			return c.Status(401).SendString("Invalid Password")
		}
		foundUser.Password = hashAndSalt([]byte(json.Password))
		db.Save(&foundUser)
		return c.SendStatus(200)
	})
}

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ProductRoutes(router fiber.Router, db *gorm.DB) {
	route := router.Group("/product", securityMiddleware)

	route.Post("/create", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	route.Post("/read", func(c *fiber.Ctx) error {
		usr, err := authenticate(c, db)
		if err != 0 {
			return c.SendStatus(err)
		}
		return c.Status(200).JSON(usr)
	})
	route.Post("/update", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	route.Post("/delete", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

}

func authenticate(c *fiber.Ctx, db *gorm.DB) (User, int) {
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		fmt.Println(err)
		return User{}, 400
	}
	query := Session{Sessionid: json.Sessionid}
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
