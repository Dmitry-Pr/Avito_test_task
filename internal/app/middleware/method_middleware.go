package middleware

import "net/http"

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
