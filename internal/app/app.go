package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/imboyko/lock-service/internal/config"
	"github.com/imboyko/lock-service/internal/logger"
	"github.com/imboyko/lock-service/internal/rest"
	"github.com/imboyko/lock-service/internal/storage"
)

func Run(runCtx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()
	stor, err := storage.NewRedisStorage(runCtx, cfg.Redis)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer stor.Close()

	log.Info("redis client connected", slog.String("redis.addr", cfg.Redis.Addr()))

	controller := rest.NewRouter(stor)

	server := http.Server{
		Addr:    ":8080",
		Handler: controller,
	}

	g, gctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		log.Info("start listening and serving", slog.String("addr", ":8080"))
		return server.ListenAndServe()
	})

	g.Go(func() error {
		select {
		case <-runCtx.Done():
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			return server.Shutdown(stopCtx)

		case <-gctx.Done():
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Warn(err.Error())
			return nil
		} else {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}
