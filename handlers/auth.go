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
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}

	found := User{}
	query := User{Username: data.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User not Found")
	}
	if !comparePasswords(found.Password, []byte(data.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
		})
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
		"message": "sucess",
		"data":    session,
	})
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	data := new(Session)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	session := Session{}
	query := Session{Sessionid: data.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Session Not Found",
		})
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}

func CreateUser(c *fiber.Ctx) error {
	type CreateUserRequest struct {
		Password string `json:"password"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	db := database.DB
	data := new(CreateUserRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	password := hashAndSalt([]byte(data.Password))
	err := checkmail.ValidateFormat(data.Email)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid Email Format",
		})
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
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "User Already Exists",
		})
	}
	db.Create(&new)
	session := Session{UserRefer: new.ID, Sessionid: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
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
		"message": "sucess",
		"data":    session,
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(User)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
		"data":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	type DeleteUserRequest struct {
		Password string `json:"password"`
	}
	db := database.DB
	data := new(DeleteUserRequest)
	user := c.Locals("user").(User)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	if !comparePasswords(user.Password, []byte(data.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
		})
	}
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}
	db := database.DB
	user := c.Locals("user").(User)
	data := new(ChangePasswordRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	if !comparePasswords(user.Password, []byte(data.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
		})
	}
	user.Password = hashAndSalt([]byte(data.NewPassword))
	db.Save(&user)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
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
