package main

import (
	"log"
	"os"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	database.ConnectDB()

	router.Initalize(app)

	log.Fatal(app.Listen(":" + getenv("PORT", "3000")))
}

// Set Env variables, PORT will auto set to 3000 if not set
// *
// PORT=3000
// DATABASE_URL="host=localhost port=5432 user=postgres password= dbname= sslmode=disable"
// *

// Docker Command for Postgres database
// *
// docker run --name database -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:alpine
// *
/// DB_URL Variable for docker Command
// *
// DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
