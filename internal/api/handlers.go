package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/imboyko/lock-service/internal/logger"
	"github.com/imboyko/lock-service/internal/storage"
	"github.com/imboyko/lock-service/internal/storage/models"
)

func HandleGetAll(s *storage.RedisStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := logger.GetCtxLogger(r.Context())
		locks, err := s.GetAll(r.Context())
		if err != nil {
			oplog.Error("get all locks failed", logger.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := renderJson(w, locks); err != nil {
			oplog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func HandleGetById(s *storage.RedisStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		oplog := logger.GetCtxLogger(r.Context()).With(slog.String("lockId", id))

		lock, err := s.GetById(r.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				oplog.Warn(err.Error())
				w.WriteHeader(http.StatusGone)
			} else {
				oplog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := renderJson(w, lock); err != nil {
			oplog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func HandlePutById(s *storage.RedisStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id := chi.URLParam(r, "id")
		oplog := logger.GetCtxLogger(r.Context()).With(slog.String("lockId", id))

		l := models.Lock{Id: id, Timestamp: time.Now(), Username: "TODO Username"} // TODO lock.Username
		err := s.Save(r.Context(), l)
		if err != nil {
			oplog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		oplog.Info("lock set")
	})
}

func HandleDeleteById(s *storage.RedisStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		oplog := logger.GetCtxLogger(r.Context()).With(slog.String("lockId", id))

		err := s.DeleteById(r.Context(), id)
		if err != nil {
			oplog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		oplog.Info("lock released")
	})
}

func renderJson[T any](w http.ResponseWriter, val T) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(val)
}
