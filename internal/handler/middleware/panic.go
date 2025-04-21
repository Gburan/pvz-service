package middleware

import (
	"log"
	"net/http"
	"time"
)

func PanicMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("New request",
					"method", r.Method,
					"url", r.URL.Path,
					"time", time.Since(start),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
