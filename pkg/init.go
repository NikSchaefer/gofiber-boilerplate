package pkg

import (
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/analytics"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

// Only autoMigrate when needed
const autoMigrate = true

// Only init when they are needed instead of global init
func InitializeServices() {
	analytics.InitAnalytics()
	validator.InitializeValidator()


	database.InitializeDB(autoMigrate)
}
