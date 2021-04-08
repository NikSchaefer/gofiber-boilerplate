package main

import (
	"log"
	"os"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "gorm.io/driver/postgres"
)

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST",
	}))

	database.ConnectDB()

	router.InitalizeRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

// Set Env variables for
// PORT=3000
// DATABASE_URL="host=localhost port=5432 user=postgres password= dbname= sslmode=disable"
// If deploying with Heroku these will be
// automatically set provided you have created a db add on
