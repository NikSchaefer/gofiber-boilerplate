package api

import (
	"fmt"

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
		fmt.Println(json.Product)
		newProduct := Product{
			ProductID: json.ProductID,
			UserRefer: usr.ID,
			Name:      json.Name,
			Value:     json.Value,
		}
		db.Create(&newProduct)
		return c.Status(200).JSON(usr)
	})
	route.Post("/read", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	route.Post("/update", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	route.Post("/delete", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
}
