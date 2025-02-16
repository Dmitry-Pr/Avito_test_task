// Package utils Description: Файл содержит функции для работы с JWT токенами.
package utils

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey - секретный ключ для подписи JWT
var SecretKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY переменная среды не найдена")
	}
	SecretKey = []byte(secretKey)
}

// GenerateJWT - функция для генерации JWT
func GenerateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Время истечения токена (24 часа)

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
