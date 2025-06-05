package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func PanicMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				slog.Error(
					"Caught panic",
					"method", r.Method,
					"url", r.URL.Path,
					"time", time.Since(start),
					"stack trace", string(debug.Stack()),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
