package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/imboyko/lock-service/internal/logger"
)

func ContextLogger(l *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			l := l.With(slog.String("requestId", middleware.GetReqID(r.Context())))

			r = r.WithContext(logger.SetCtxLogger(r.Context(), l))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
