package handlers

import (
	"encoding/json"
	"time"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/badoux/checkmail"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User model.User
type Session model.Session
type Product model.Product

func GetUser(sessionid guuid.UUID) (User, int) {
	db := database.DB
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
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	db := database.DB
	data := new(LoginRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := User{}
	query := User{Username: data.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User not Found")
	}
	if !comparePasswords(found.Password, []byte(data.Password)) {
		return c.Status(fiber.StatusBadRequest).SendString("Incorrect Password")
	}
	session := Session{UserRefer: found.ID, Expires: SessionExpires(), Sessionid: guuid.New()}
	db.Create(&session)
	c.Cookie(&fiber.Cookie{
		Name:     "sessionid",
		Expires:  SessionExpires(),
		Value:    session.Sessionid.String(),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(session)
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	data := new(Session)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	session := Session{}
	query := Session{Sessionid: data.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("Session Not Found")
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.SendStatus(fiber.StatusOK)
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	data := new(User)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	password := hashAndSalt([]byte(data.Password))
	err := checkmail.ValidateFormat(data.Email)
	if err != nil {
		return c.Status(400).SendString("Invalid Email Format")
	}
	new := User{
		Username: data.Username,
		Password: password,
		Email:    data.Email,
		ID:       guuid.New(),
	}
	found := User{}
	query := User{Username: data.Username}
	err = db.First(&found, &query).Error
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
	user := c.Locals("user").(User)
	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	type DeleteUserRequest struct {
		password string
	}
	db := database.DB
	data := new(DeleteUserRequest)
	user := c.Locals("user").(User)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if !comparePasswords(user.Password, []byte(data.password)) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Password")
	}
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)
	c.ClearCookie("sessionid")
	return c.SendStatus(fiber.StatusOK)
}

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		NewPassword string `json:"newPassword"`
	}
	db := database.DB
	user := c.Locals("user").(User)
	data := new(ChangePasswordRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if !comparePasswords(user.Password, []byte(data.NewPassword)) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Password")
	}
	user.Password = hashAndSalt([]byte(data.NewPassword))
	db.Save(&user)
	return c.SendStatus(fiber.StatusOK)
}

func hashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

// Universal date the Session Will Expire
func SessionExpires() time.Time {
	return time.Now().Add(5 * 24 * time.Hour)
}
