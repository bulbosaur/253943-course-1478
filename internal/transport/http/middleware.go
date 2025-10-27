package srv

import (
	"context"
	"lyceum/logger"
	"net/http"

	"github.com/google/uuid"
)

func LoggingMiddleware(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqIDHeader := r.Header.Get("x-request-id")
			if reqIDHeader != "" {
				logger.WithRequestID(r.Context(), reqIDHeader)
			} else {
				logger.WithRequestID(r.Context(), uuid.NewString())
			}

			next.ServeHTTP(w, r)
		})
	}
}