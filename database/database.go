package database

import (
	"log"
	"os"

	"github.com/NikSchaefer/go-fiber/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() {
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = database
	db.AutoMigrate(&model.User{}, &model.Session{}, &model.Product{})
}
