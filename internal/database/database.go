package database

import (
	"context"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/NikSchaefer/go-fiber/ent"
)

var DB *ent.Client

func InitializeDB(autoMigrate bool) {
	client, err := ent.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	DB = client

	// TODO: Should move this to a lambda function that can be triggered
	// with the github action only run migrations if explicitly requested
	if autoMigrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}
}

func CloseDB() {
	DB.Close()
}
