// Package middleware Description: Middleware для логирования метода запроса, пути и времени обработки.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// LogsMiddleware логирует метод запроса, путь и время обработки.
func LogsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)
	})
}
