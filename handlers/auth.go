package handlers

import (
	"fmt"
	"time"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ChangePassword struct {
	User
	NewPassword string `json:"sessionid"`
}

type User model.User
type Session model.Session
type Product model.Product

var db *gorm.DB = database.DB

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


func Login(c *fiber.Ctx) error {
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	empty := User{}
	if json.Username == empty.Username || empty.Password == json.Password {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Data Sent")
	}

	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User not Found")
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
		return c.Status(fiber.StatusBadRequest).SendString("Incorrect Password")
	}
	session := Session{UserRefer: found.ID, Expires: SessionExpires(), Sessionid: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Creation Error")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "sessionid",
		Expires:  SessionExpires(),
		Value:    session.Sessionid.String(),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(session)
}

func Logout(c *fiber.Ctx) error {
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if json.Sessionid == new(Session).Sessionid {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Data Sent")
	}
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("Session Not Found")
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.SendStatus(fiber.StatusOK)
}

func CreateUser(c *fiber.Ctx) error {
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	empty := User{}
	if json.Username == empty.Username || empty.Password == json.Password {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Data Sent")
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
		return c.Status(fiber.StatusBadRequest).SendString("User Already Exists")
	}
	db.Create(&new)
	session := Session{UserRefer: new.ID, Sessionid: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Creation Error")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "sessionid",
		Expires:  time.Now().Add(5 * 24 * time.Hour),
		Value:    session.Sessionid.String(),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(session)
}
func GetUserInfo(c *fiber.Ctx) error {
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user, status := GetUser(json.Sessionid, db)
	if status != 0 {
		return c.SendStatus(status)
	}
	var Products []model.Product = []model.Product{}
	db.Model(&user).Association("Products").Find(&Products)
	user.Products = Products
	return c.JSON(user)
}
func DeleteUser(c *fiber.Ctx) error {
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	empty := User{}
	if json.Username == empty.Username || empty.Password == json.Password {
		return c.Status(401).SendString("Invalid Data Sent")
	}
	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User Not Found")
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Credentials")
	}
	db.Model(&found).Association("Sessions").Delete()
	db.Model(&found).Association("Products").Delete()
	db.Delete(&found)
	c.ClearCookie("sessionid")
	return c.SendStatus(fiber.StatusOK)
}
func UpdateUser(c *fiber.Ctx) error {
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	empty := User{}
	if json.Username == empty.Username || empty.Password == json.Password {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Data Sent")
	}
	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User Not Found")
	}
	return c.SendStatus(fiber.StatusOK)
}
func ChangePasswordRoute(c *fiber.Ctx) error {
	json := new(ChangePassword)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User Not Found")
	}
	if !comparePasswords(found.Password, []byte(json.NewPassword)) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Password")
	}
	found.Password = hashAndSalt([]byte(json.Password))
	db.Save(&found)
	return c.SendStatus(fiber.StatusOK)
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
	return err == nil
}

// Universal date the Session Will Expire
func SessionExpires() time.Time {
	return time.Now().Add(5 * 24 * time.Hour)
}
