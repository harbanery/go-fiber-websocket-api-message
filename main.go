package main

import (
	"gofiber-chat-api/src/configs"
	"gofiber-chat-api/src/helpers"
	"gofiber-chat-api/src/routes"
	"gofiber-chat-api/src/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/joho/godotenv"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return "0.0.0.0:" + port

}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	}))

	services.InitHub()
	configs.ConnectToDB()
	helpers.Migration()
	routes.Endpoint(app)

	if err := app.Listen(getPort()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
