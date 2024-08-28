package models

import (
	"gofiber-chat-api/src/configs"
)

type User struct {
	Base
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func SelectUserfromEmail(email *string) *User {
	var user *User
	configs.DB.First(&user, "email = ?", &email)
	return user
}

func CreateUser(user *User) error {
	result := configs.DB.Create(&user)
	return result.Error
}
