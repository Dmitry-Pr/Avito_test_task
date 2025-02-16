// Package middleware Description: Middleware для проверки JWT токена
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"merch-shop/internal/pkg/errors"
	"net/http"
	"strings"

	"merch-shop/internal/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

// ExcludedPaths - пути, которые не нужно проверять на наличие JWT токена
var ExcludedPaths = []string{"/api/auth"}

type contextKey string // Define a new type

const (
	// UserIDKey - ключ id пользователя для контекста
	UserIDKey contextKey = "user_id"
)

// AuthMiddleware - middleware для проверки JWT токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range ExcludedPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			jsonErr := errors.NewErrorResponse("Не авторизован")
			err := json.NewEncoder(w).Encode(jsonErr)
			if err != nil {
				return
			}
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неизвестный метод подписи токена: %v", token.Method)
			}
			return utils.SecretKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			jsonErr := errors.NewErrorResponse("Не авторизован")
			err := json.NewEncoder(w).Encode(jsonErr)
			if err != nil {
				return
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			jsonErr := errors.NewErrorResponse("Не авторизован")
			err := json.NewEncoder(w).Encode(jsonErr)
			if err != nil {
				return
			}
			return
		}
	})
}
