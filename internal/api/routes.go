package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/imboyko/lock-service/internal/storage"
)

func NewRouter(s *storage.RedisStorage, l *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.CleanPath, ContextLogger(l))

	// TODO Bearer auth, inject username

	r.Route("/locks", func(r chi.Router) {
		r.Get("/", HandleGetAll(s))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", HandleGetById(s))
			r.Put("/", HandlePutById(s))
			r.Delete("/", HandleDeleteById(s))
		})
	})
	return r
}
