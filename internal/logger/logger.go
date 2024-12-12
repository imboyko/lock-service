package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	h := slog.NewTextHandler(os.Stderr, nil)
	l := slog.New(h)
	return l
}

func Error(err error) slog.Attr {
	return slog.Any("error", err)
}
