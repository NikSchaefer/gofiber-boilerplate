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
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	user := c.Locals("user").(User)
	newProduct := Product{
		UserRefer: user.ID,
		Name:      data.Name,
		Value:     data.Value,
	}
	err := db.Create(&newProduct).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}
func GetProducts(c *fiber.Ctx) error {
	db := database.DB
	Products := []Product{}
	db.Model(&model.Product{}).Order("ID asc").Limit(100).Find(&Products)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
		"data":    Products,
	})
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
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
		"data":    product,
	})
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
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	param := c.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid ID Format",
		})
	}
	found := Product{}
	query := Product{
		ID:        id,
		UserRefer: user.ID,
	}
	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Product Not Found",
		})
	}
	if data.Name != "" {
		found.Name = data.Name
	}
	if data.Value != "" {
		found.Value = data.Value
	}
	db.Save(&found)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
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
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}
