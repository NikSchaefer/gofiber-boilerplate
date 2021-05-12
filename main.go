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
		AllowOrigins: "*", // comma string format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	database.ConnectDB()

	router.Initalize(app)
	log.Fatal(app.Listen(":" + getenv("PORT", "3000")))
}

/*
ENV Variables:
will auto set to 3000 if not set
PORT=3000
this should be a connection string or url
DATABASE_URL="host=localhost port=5432 user=postgres password= dbname= sslmode=disable"
**
Docker Command for Postgres database:
docker run --name database -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:alpine

DB_URL Variable for docker database
DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
**
Docker build base image in first stage for development
docker build --target build -t base .
**
run dev container
docker run -p 3000:3000 --mount type=bind,source="C:\Users\schaefer\go\src\fiber",target=/go/src/app --name fiber -td base
**
rebuild and run package
docker exec -it web go run main.go
**
stop and remove container
docker stop fiber; docker rm fiber
*/
