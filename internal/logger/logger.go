package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxValues int

const ctxLoggerKey ctxValues = 0

func New() *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(h)

	slog.SetDefault(l)

	return l
}

func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

// Усланавливает логгер l в контекст ctx
func SetCtxLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, l)
}

// Возвращает логгер из контекста ctx.
// Если значение не установлено, то вернет slog.Default()
func GetCtxLogger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLoggerKey).(*slog.Logger); ok {
		return l
	} else {
		return slog.Default()
	}
}
