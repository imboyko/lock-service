package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/imboyko/lock-service/internal/storage/models"
)

type Storage interface {
	// Возвращает все установленные блокировки
	GetAll(context.Context) ([]models.Lock, error)
	// Возвращает блокировку по идентификатору
	GetById(context.Context, string) (models.Lock, error)
	// Удаляет блокировку по идентификатору
	DeleteById(context.Context, string) error
	// Добавляет или обновляет блокировку
	Save(context.Context, models.Lock) error
}

func NewRouter(s Storage) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// TODO Bearer auth, inject context logger and username

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
