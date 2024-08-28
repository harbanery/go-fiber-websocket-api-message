package controllers

import (
	"gofiber-chat-api/src/helpers"
	"gofiber-chat-api/src/middlewares"
	"gofiber-chat-api/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var bodyRequest models.User
	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&bodyRequest).(*models.User)

	if existUser := models.SelectUserfromEmail(&user.Email); existUser.ID != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Email already exists",
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Password error",
		})
	}

	user.Password = string(hashPassword)
	if err = models.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "User created successfully.",
	})
}

func Login(c *fiber.Ctx) error {
	var bodyRequest models.User
	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&bodyRequest).(*models.User)
	existUser := models.SelectUserfromEmail(&user.Email)
	if existUser.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Invalid password",
		})
	}

	payload := map[string]interface{}{
		"id":    existUser.ID,
		"email": existUser.Email,
	}

	token, err := helpers.GenerateToken(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Login successfully",
		"email":      existUser.Email,
		"id":         existUser.ID,
		"token":      token,
	})
}

func Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Logout error",
		})
	}

	email, ok := user["email"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "unauthorized",
			"statusCode": 401,
			"message":    "Logout error",
		})
	}

	if existUser := models.SelectUserfromEmail(&email); existUser.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 201,
		"message":    "Logout successfully",
	})
}
