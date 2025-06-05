package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"pvz-service/internal/logging"
	"pvz-service/internal/metrics"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var pathTemplate string
		if route := mux.CurrentRoute(r); route != nil {
			var err error
			pathTemplate, err = route.GetPathTemplate()
			if err != nil {
				slog.WarnContext(ctx, err.Error())
			}
		}

		metrics.IncRestRequestsTotal(pathTemplate)

		requestId := uuid.New()
		slog.InfoContext(ctx, fmt.Sprintf("Start [%s] request processing", requestId.String()))
		start := time.Now()

		ctx = logging.WithLogRequestID(ctx, requestId)
		ctx = logging.WithLogRequestPath(ctx, r.URL.Path)
		ctx = logging.WithLogRequestMethod(ctx, r.Method)

		rw := &responseWriter{w, http.StatusOK}
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)

		timeServe := time.Since(start)
		ctx = r.Context()
		ctx = logging.WithLogRequestStatus(ctx, rw.statusCode)
		ctx = logging.WithLogRequestDuration(ctx, timeServe.String())

		slog.InfoContext(ctx, fmt.Sprintf("Ended [%s] request processing", requestId.String()))

		metrics.IncRestResponsesDuration(pathTemplate, r.Method, timeServe)
		metrics.IncRestResponsesStatusesTotal(pathTemplate, rw.statusCode)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
