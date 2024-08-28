package configs

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	url := os.Getenv("POSTGRES_URL")

	var err error

	DB, err = gorm.Open(postgres.Open(url), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}
