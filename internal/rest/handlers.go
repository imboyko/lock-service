package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/imboyko/lock-service/internal/models"
	"github.com/imboyko/lock-service/internal/storage"
)

type getter interface {
	GetAll(context.Context) ([]models.Lock, error)
}

func HandleGetAll(s getter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locks, err := s.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(locks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

type getterById interface {
	GetById(context.Context, string) (models.Lock, error)
}

func HandleGetById(s getterById) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lock, err := s.GetById(r.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				w.WriteHeader(http.StatusGone)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(lock); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

type saver interface {
	Save(context.Context, models.Lock) error
}

func HandlePutById(s saver) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO получение Username из контекста r.Context()
		l := models.Lock{Id: id, Timestamp: time.Now(), Username: "TODO Username"}
		err := s.Save(r.Context(), l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

type deleterById interface {
	DeleteById(context.Context, string) error
}

func HandleDeleteById(s deleterById) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.DeleteById(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
