// Package middleware Description: Middleware для проверки разрешенных методов запроса.
package middleware

import "net/http"

// MethodMiddleware checks if the request method is allowed
func MethodMiddleware(next http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range allowedMethods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}
