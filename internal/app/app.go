package app

import (
	"context"
	"log/slog"

	"github.com/imboyko/lock-service/internal/config"
	"github.com/imboyko/lock-service/internal/logger"
	"github.com/imboyko/lock-service/internal/storage"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()
	stor, err := storage.NewRedisStorage(ctx, cfg.Redis)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer stor.Close()

	log.Info("redis client connected", slog.String("redis.addr", cfg.Redis.Addr()))

	<-ctx.Done()

	log.Warn(ctx.Err().Error())

	return nil
}
