// Package utils Description: Файл содержит функции для работы с JWT токенами.
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT - функция для генерации JWT
func GenerateJWT(userID uint, secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Время истечения токена (24 часа)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
