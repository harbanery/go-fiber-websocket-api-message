package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	secretKey := os.Getenv("SECRET_KEY_JWT")
	claims := token.Claims.(jwt.MapClaims)

	for key, value := range payload {
		claims[key] = value
	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
