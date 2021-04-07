package handlers

import (
	"github.com/NikSchaefer/go-fiber/database"
	guuid "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRequest struct {
	Product
	Sessionid guuid.UUID
}

func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user := c.Locals("user").(User)
	newProduct := Product{
		ProductID: guuid.New(),
		UserRefer: user.ID,
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
	db := database.DB
	user := c.Locals("user").(User)
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	Products := []Product{}
	db.Model(&user).Association("Products").Find(&Products)
	return c.Status(fiber.StatusOK).JSON(Products)
}
func GetProductById(c *fiber.Ctx) error {
	db := database.DB
	id, err := guuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id format")
	}
	product := Product{}
	query := Product{
		ProductID: id,
	}
	err = db.First(&product, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(User)
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	found := Product{}
	query := Product{
		Name:      json.Name,
		UserRefer: user.ID,
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
	db := database.DB
	user := c.Locals("user").(User)
	json := new(ProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	found := Product{}
	query := Product{
		Name:      json.Name,
		UserRefer: user.ID,
	}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	found.Value = json.Value
	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
