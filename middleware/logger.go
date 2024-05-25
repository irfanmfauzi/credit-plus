package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.StatusCode = statusCode
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := responseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		startTime := time.Now()
		duration := time.Since(startTime)

		next.ServeHTTP(&rec, r)

		logger := slog.Default().With(
			"url", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status_code", rec.StatusCode,
			"status", http.StatusText(rec.StatusCode),
		)

		if rec.StatusCode < 400 {
			logger.Info("logging_http_request")
		} else if rec.StatusCode >= 400 && rec.StatusCode < 500 {
			logger.Warn("logging_http_request")
		} else {
			logger.Error("logging_http_request")
		}

	})
}
