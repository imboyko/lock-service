package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/imboyko/lock-service/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := app.Run(ctx)
	if err != nil {
		os.Exit(1)
	}
}
