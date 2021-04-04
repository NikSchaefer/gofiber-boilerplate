package api

import (
	guuid "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProductRoutes(router fiber.Router, db *gorm.DB) {
	route := router.Group("/product", JsonMiddleware)

	route.Post("/create", func(c *fiber.Ctx) error {
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
	})
	route.Post("/read", func(c *fiber.Ctx) error {
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
	})
	route.Post("/update", func(c *fiber.Ctx) error {
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
	})
	route.Post("/delete", func(c *fiber.Ctx) error {
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
			return c.Status(401).SendString("Product Not Found")
		}
		found.Value = json.Value
		db.Delete(&found)
		return c.SendStatus(fiber.StatusOK)
	})
}
