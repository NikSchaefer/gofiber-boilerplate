package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
	data := new(Product)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user := c.Locals("user").(User)
	newProduct := Product{
		UserRefer: user.ID,
		Name:      data.Name,
		Value:     data.Value,
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
	db.Model(&model.Product{}).Order("ID asc").Limit(100).Find(&Products)
	return c.Status(fiber.StatusOK).JSON(Products)
}
func GetProductById(c *fiber.Ctx) error {
	db := database.DB
	param := c.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id format")
	}
	product := Product{}
	query := Product{ID: id}
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
	data := new(UpdateProductRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	param := c.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id format")
	}
	found := Product{}
	query := Product{
		ID:        id,
		UserRefer: user.ID,
	}
	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).SendString("Product Not Found")
	}
	if data.Name != "" {
		found.Name = data.Name
	}
	if data.Value != "" {
		found.Value = data.Value
	}
	db.Save(&found)
	return c.SendStatus(fiber.StatusOK)
}
func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(User)
	param := c.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id format")
	}
	found := Product{}
	query := Product{
		ID:        id,
		UserRefer: user.ID,
	}
	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
