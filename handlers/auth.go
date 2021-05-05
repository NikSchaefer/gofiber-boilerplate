package handlers

import (
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
	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
			"sucess":  false,
		})
	}

	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
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
	return c.JSON(fiber.Map{
		"code":    200,
		"message": nil,
		"sucess":  true,
		"data":    session,
	})
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": nil,
		"sucess":  true,
	})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	password := hashAndSalt([]byte(json.Password))
	err := checkmail.ValidateFormat(json.Email)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid Email Address",
			"sucess":  false,
		})
	}
	new := User{
		Username: json.Username,
		Password: password,
		Email:    json.Email,
		ID:       guuid.New(),
	}
	found := User{}
	query := User{Username: json.Username}
	err = db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "User already exists",
			"sucess":  false,
		})
	}
	db.Create(&new)
	session := Session{UserRefer: new.ID, Sessionid: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "sessionid",
		Expires:  time.Now().Add(5 * 24 * time.Hour),
		Value:    session.Sessionid.String(),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{
		"code":    200,
		"message": nil,
		"sucess":  true,
		"data":    session,
	})
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
	json := new(DeleteUserRequest)
	user := c.Locals("user").(User)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	if !comparePasswords(user.Password, []byte(json.password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
			"sucess":  false,
		})
	}
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": nil,
		"sucess":  true,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}
	db := database.DB
	user := c.Locals("user").(User)
	json := new(ChangePasswordRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
			"sucess":  false,
		})
	}
	if !comparePasswords(user.Password, []byte(json.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
			"sucess":  false,
		})
	}
	user.Password = hashAndSalt([]byte(json.NewPassword))
	db.Save(&user)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": nil,
		"sucess":  true,
	})
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
