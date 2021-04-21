package database

import (
	"log"
	"os"

	"github.com/NikSchaefer/go-fiber/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error // define error here to prevent overshadowing the global ConnectDB

	// DATABASE_URL should be a connection string or url of a postgres database
	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&model.User{}, &model.Session{}, &model.Product{})
}
