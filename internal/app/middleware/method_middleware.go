// Package middleware Description: Middleware для проверки разрешенных методов запроса.
package middleware

import (
	"encoding/json"
	"merch-shop/internal/pkg/errors"
	"net/http"
)

// MethodMiddleware checks if the request method is allowed
func MethodMiddleware(next http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range allowedMethods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		jsonErr := errors.NewErrorResponse("Недопустимый HTTP метод")
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			return
		}
	})
}
