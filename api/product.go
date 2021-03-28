package api

import (
	guuid "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProductRoutes(router fiber.Router, db *gorm.DB) {
	route := router.Group("/product", SecurityMiddleware)

	route.Post("/create", func(c *fiber.Ctx) error {
		json := new(ProductRequest)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		usr, err := GetUser(json.Sessionid, db)
		if err != 0 {
			return c.SendStatus(err)
		}
		newProduct := Product{
			ProductID: guuid.New(),
			UserRefer: usr.ID,
			Name:      json.Name,
			Value:     json.Value,
		}
		er := db.Create(&newProduct).Error
		if er != nil {
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	})
	route.Post("/read", func(c *fiber.Ctx) error {
		json := new(ProductRequest)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		usr, err := GetUser(json.Sessionid, db)
		if err != 0 {
			return c.SendStatus(err)
		}
		Products := []Product{}
		db.Model(&usr).Association("Products").Find(&Products)
		return c.Status(200).JSON(Products)
	})
	route.Post("/update", func(c *fiber.Ctx) error {
		json := new(ProductRequest)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		usr, err := GetUser(json.Sessionid, db)
		if err != 0 {
			return c.SendStatus(err)
		}

		return c.Status(200).JSON(usr)
	})
	route.Post("/delete", func(c *fiber.Ctx) error {
		json := new(ProductRequest)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(400)
		}
		usr, err := GetUser(json.Sessionid, db)
		if err != 0 {
			return c.SendStatus(err)
		}
		return c.Status(200).JSON(usr)

	})
}
