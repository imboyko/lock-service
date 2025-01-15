package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"

	"github.com/imboyko/lock-service/internal/logger"
)

type ctxKeys int

const usernameCtxKey ctxKeys = iota

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

func AuthMiddleware(secret string) func(next http.Handler) http.Handler {
	keyFn := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetCtxLogger(r.Context())

			t, err := request.ParseFromRequest(r, request.BearerExtractor{}, keyFn, []request.ParseFromRequestOption{}...)
			if err != nil {
				log.Warn("authentification failed", logger.Error(err))
				renderError(w, err.Error(), http.StatusUnauthorized)
				return
			}

			username, _ := t.Claims.GetSubject()
			r = r.WithContext(setUsernameCtx(r.Context(), username))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func setUsernameCtx(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameCtxKey, username)
}

func getUsernameCtx(ctx context.Context) string {
	return ctx.Value(usernameCtxKey).(string)
}
