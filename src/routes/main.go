package routes

import (
	"gofiber-chat-api/src/controllers"
	"gofiber-chat-api/src/middlewares"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Endpoint(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"uri":     ctx.Request().URI().String(),
			"path":    ctx.Path(),
			"message": "Server is running.",
		})
	})

	app.Get("/v1", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v1",
		})
	})

	app.Get("/v2", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v2",
		})
	})

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/logout", middlewares.JWTRestricted(), controllers.Logout)

	app.Post("/chat/room", controllers.CreateRoom)
	app.Get("/chat/room", middlewares.JWTRestricted(), controllers.GetRooms)
	app.Use("/chat", middlewares.AllowUpgrade)
	app.Get("/chat/:id", middlewares.JWTRestricted(), websocket.New(controllers.ConversationChat))

	// api := app.Group("/api")
}
