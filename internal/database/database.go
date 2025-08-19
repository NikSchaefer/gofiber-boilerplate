package database

import (
	"context"
	"log"

	"github.com/NikSchaefer/go-fiber/config"
	"github.com/NikSchaefer/go-fiber/ent"
	_ "github.com/lib/pq"
)

var DB *ent.Client

func InitializeDB(autoMigrate bool) {
	client, err := ent.Open("postgres", config.GetDatabaseURL())
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
