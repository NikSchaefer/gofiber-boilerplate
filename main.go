package main

import (
	"fmt"

	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/internal/router"
	util "github.com/NikSchaefer/go-fiber/pkg"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitializeApp() *fiber.App {
	godotenv.Load()

	util.InitializeServices()

	app := *fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // replace with your domain (e.g. google.com)
		AllowHeaders:     "Origin, Content-Type, Accept, Cookie",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
	}))

	app.Use(logger.New())
	router.Initialize(&app)

	fmt.Println("App initialized")

	return &app
}

func main() {
	app := InitializeApp()

	defer database.CloseDB()
	err := app.Listen(":8000")
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
