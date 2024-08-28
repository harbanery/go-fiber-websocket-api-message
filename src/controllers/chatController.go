package controllers

import (
	"gofiber-chat-api/src/middlewares"
	"gofiber-chat-api/src/models"
	"gofiber-chat-api/src/services"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateRoom(c *fiber.Ctx) error {
	var bodyRequest models.Chat
	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	chat := middlewares.XSSMiddleware(&bodyRequest).(*models.Chat)

	if err := models.CreateChat(chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create chat room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chat room created successfully.",
	})
}

func GetRooms(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Logout error",
		})
	}

	userID, ok := user["id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Logout error",
		})
	}

	chats := models.SelectChatsbyUserID(&userID)
	if len(chats) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "no content",
			"statusCode": 202,
			"message":    "Product is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chats OK",
		"data":       chats,
	})
}

func ConversationChat(c *websocket.Conn) {
	hub := c.Locals("hub").(*services.Hub)

	// Mendapatkan chat_id dari query parameter
	chatID := c.Params("id") // Misalnya, URL: /ws/chat/:id

	chat := models.SelectChatbyID(&chatID)
	if chat.ID == "" {
		// Jika tidak ada chat_id, tutup koneksi
		_ = c.WriteMessage(websocket.CloseMessage, []byte{})
		_ = c.Close()
		return
	}

	// Mendapatkan atau membuat Room berdasarkan id
	room := hub.GetRoom(chatID)

	// Mendaftarkan koneksi ke Room
	room.Register <- c

	// Pastikan koneksi di-unregister dan ditutup saat selesai
	defer func() {
		room.Unregister <- c
		_ = c.Close()
	}()

	for {
		// var message models.Message
		// Membaca pesan dari klien
		// if err := c.ReadJSON(&message); err != nil {
		// 	log.Println("Read message error:", err)
		// 	break
		// }

		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			break
		}

		user, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			log.Println("Read message error: jwt user")
			return
		}

		userID, ok := user["id"].(string)
		if !ok {
			log.Println("Read message error: get id from jwt user")
			return
		}

		if chatMember := models.SelectChatUserbyChatUserID(&chatID, &userID); chatMember.ID == "" {
			log.Println("Read message error: this user is no entry")
			return
		}

		message := models.Message{
			ChatID: chat.ID,
			UserID: userID, // Replace with the actual user ID
			Body:   string(msg),
			Status: "NONE",
		}

		// Mengisi ChatID dan UserID dari token JWT (asumsi autentikasi telah dilakukan)
		// Misalnya:
		// userID := c.Locals("user_id").(string)
		// chat.UserID = userID
		// chat.ChatID = ChatID

		// Menyimpan pesan chat ke dalam database
		// Pastikan ChatID diisi
		if err := models.CreateMessage(&message); err != nil {
			// Anda dapat menangani error sesuai kebutuhan
			log.Println("Failed to save message:", err)
			break
		}

		// Menyiarkan pesan ke semua klien dalam Room
		room.Broadcast <- message
	}
}
