package handlers

import (
	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	guuid "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
	json := new(Product)
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
func GetProducts(c *fiber.Ctx) error {
	db := database.DB
	Products := []Product{}
	db.Model(&model.Product{}).Find(&Products)
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
	type UpdateProductRequest struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Sessionid string `json:"sessionid"`
	}
	db := database.DB
	user := c.Locals("user").(User)
	json := new(UpdateProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id, err := guuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Id Format")
	}
	found := Product{}
	query := Product{
		ProductID: id,
		UserRefer: user.ID,
	}
	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).SendString("Product Not Found")
	}
	if json.Name != "" {
		found.Name = json.Name
	}
	if json.Value != "" {
		found.Value = json.Value
	}
	db.Save(&found)
	return c.SendStatus(fiber.StatusOK)
}
func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(User)
	id, err := guuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Id Format")
	}

	found := Product{}
	query := Product{
		ProductID: id,
		UserRefer: user.ID,
	}
	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
