package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Initalize and set the authentication and authorization routes
func AuthRoutes(router fiber.Router, db *gorm.DB) {
	route := router.Group("/auth", SecurityMiddleware)
	route.Post("/login", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}

		found := User{}
		query := User{Username: json.Username}
		err := db.First(&found, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User not Found")
		}
		if !comparePasswords(found.Password, []byte(json.Password)) {
			return c.Status(401).SendString("Incorrect Password")
		}
		session := Session{UserRefer: found.ID, Sessionid: guuid.New()}
		err = db.Create(&session).Error
		if err != nil {
			return c.Status(500).SendString("Creation Error")
		}
		c.Cookie(&fiber.Cookie{
			Name:     "sessionid",
			Expires:  time.Now().Add(5 * 24 * time.Hour),
			Value:    session.Sessionid.String(),
			HTTPOnly: true,
		})
		return c.Status(200).JSON(session)
	})

	route.Post("/logout", func(c *fiber.Ctx) error {
		json := new(Session)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
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
	route.Post("/create", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		pw := hashAndSalt([]byte(json.Password))
		new := User{
			Username: json.Username,
			Password: pw,
			ID:       guuid.New(),
		}
		found := User{}
		query := User{Username: json.Username}
		err := db.First(&found, &query).Error
		if err != gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Already Exists")
		}
		db.Create(&new)
		return c.SendStatus(200)
	})
	route.Post("/user", func(c *fiber.Ctx) error {
		json := new(Session)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		user, err := GetUser(json.Sessionid, db)
		if err != 0 {
			return c.SendStatus(err)
		}
		Sessions := []Session{}
		Products := []Product{}
		db.Model(&user).Association("Sessions").Find(&Sessions)
		db.Model(&user).Association("Products").Find(&Products)
		user.Products = Products
		user.Sessions = Sessions
		return c.JSON(user)
	})
	route.Post("/delete", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		found := User{}
		query := User{Username: json.Username}
		err := db.First(&found, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		if !comparePasswords(found.Password, []byte(json.Password)) {
			return c.Status(401).SendString("Invalid Credentials")
		}
		db.Model(&found).Association("Sessions").Delete()
		db.Model(&found).Association("Products").Delete()
		db.Delete(&found)
		c.ClearCookie("sessionid")
		return c.SendStatus(200)
	})
	route.Post("/update", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		found := User{}
		query := User{Username: json.Username}
		err := db.First(&found, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		return c.SendStatus(200)
	})
	route.Post("/changepassword", func(c *fiber.Ctx) error {
		json := new(ChangePassword)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		found := User{}
		query := User{Username: json.Username}
		err := db.First(&found, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		if !comparePasswords(found.Password, []byte(json.NewPassword)) {
			return c.Status(401).SendString("Invalid Password")
		}
		found.Password = hashAndSalt([]byte(json.Password))
		db.Save(&found)
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
