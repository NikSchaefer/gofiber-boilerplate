package handlers

import (
	guuid "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRequest struct {
	Product
	Sessionid guuid.UUID
}

func CreateProduct(c *fiber.Ctx) error {
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	usr, status := GetUser(json.Sessionid, db)
	if status != 0 {
		return c.SendStatus(status)
	}
	newProduct := Product{
		ProductID: guuid.New(),
		UserRefer: usr.ID,
		Name:      json.Name,
		Value:     json.Value,
	}
	err := db.Create(&newProduct).Error
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusOK)
}
func GetProduct(c *fiber.Ctx) error {
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	usr, err := GetUser(json.Sessionid, db)
	if err != 0 {
		return c.SendStatus(err)
	}
	Products := []Product{}
	db.Model(&usr).Association("Products").Find(&Products)
	return c.Status(fiber.StatusOK).JSON(Products)
}
func UpdateProduct(c *fiber.Ctx) error {
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	usr, status := GetUser(json.Sessionid, db)
	if status != 0 {
		return c.SendStatus(status)
	}
	found := Product{}
	query := Product{
		Name:      json.Name,
		UserRefer: usr.ID,
	}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).SendString("Product Not Found")
	}
	found.Value = json.Value
	db.Save(&found)
	return c.SendStatus(fiber.StatusOK)
}
func DeleteProduct(c *fiber.Ctx) error {
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	usr, status := GetUser(json.Sessionid, db)
	if status != 0 {
		return c.SendStatus(status)
	}
	found := Product{}
	query := Product{
		Name:      json.Name,
		UserRefer: usr.ID,
	}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	found.Value = json.Value
	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
