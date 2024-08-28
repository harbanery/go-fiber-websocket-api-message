package middlewares

import (
	"gofiber-chat-api/src/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func AllowUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("hub", services.GetHub())
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
