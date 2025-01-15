package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/imboyko/lock-service/internal/storage"
)

func NewRouter(s *storage.RedisStorage, l *slog.Logger, jwtSecret string) http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.Recoverer,
		middleware.CleanPath,
		cors.AllowAll().Handler,
		middleware.NoCache,
		ContextLogger(l),
	)

	r.Route("/locks", func(r chi.Router) {
		r.Get("/", HandleGetAll(s))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", HandleGetById(s))
			r.Group(func(r chi.Router) {
				r.Use(AuthMiddleware(jwtSecret))
				r.Put("/", HandlePutById(s))
				r.Delete("/", HandleDeleteById(s))
			})
		})
	})
	return r
}
