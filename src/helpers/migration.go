package helpers

import (
	"gofiber-chat-api/src/configs"
	"gofiber-chat-api/src/models"
	"log"
)

func Migration() {
	err := configs.DB.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.ChatUser{},
		&models.Message{},
	)

	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
}
