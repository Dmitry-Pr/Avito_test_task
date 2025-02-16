package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"merch-shop/internal/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

var ExcludedPaths = []string{"/api/auth"}

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
			http.Error(w, "Неавторизован", http.StatusUnauthorized)
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
			http.Error(w, "Неавторизован", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))

			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Неавторизован", http.StatusUnauthorized)
			return
		}
	})
}
