package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/imboyko/lock-service/internal/config"
	"github.com/imboyko/lock-service/internal/storage/models"
)

var ErrNotFound = errors.New("lock not found")

func NewRedisStorage(ctx context.Context, cfg config.Redis) (*RedisStorage, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		DB:       cfg.Db,
		Username: cfg.Username,
		Password: cfg.Password,
	})

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis server: %w", err)
	}

	return &RedisStorage{
		rdb: c,
		ttl: 60 * time.Second,
	}, nil
}

type RedisStorage struct {
	rdb *redis.Client
	ttl time.Duration
}

// Вызывает Close() у redis.Client
func (s *RedisStorage) Close() error {
	return s.rdb.Close()
}

// Получает все ключи, соответствующие шаблону lock:*
func (s *RedisStorage) GetAll(ctx context.Context) ([]models.Lock, error) {
	var (
		locks   []models.Lock = make([]models.Lock, 0)
		pattern lockId        = "*"
	)

	keys, _, err := s.rdb.Scan(ctx, 0, pattern.key(), 0).Result()
	if err != nil {
		return locks, fmt.Errorf("failed to scan keys: %w", err)
	}

	cmds, err := s.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, key := range keys {
			p.HGetAll(ctx, key)
		}
		return nil
	})
	if err != nil {
		return locks, fmt.Errorf("failed to exec pipe: %w", err)
	}

	for _, cmd := range cmds {
		var val models.Lock
		cmd.(*redis.MapStringStringCmd).Scan(&val) // TODO handle error
		locks = append(locks, val)
	}

	return locks, nil
}

// Возвращает значение ключа lock:id
func (s *RedisStorage) GetById(ctx context.Context, id string) (models.Lock, error) {
	var res models.Lock

	key := lockId(id).key()

	tx := s.rdb.TxPipeline()
	existsCmd := tx.Exists(ctx, lockId(id).key())
	hgetallCmd := tx.HGetAll(ctx, key)

	if _, err := tx.Exec(ctx); err != nil {
		return res, fmt.Errorf("failed to exec transaction: %w", err)
	}

	exists, err := existsCmd.Result()
	if err != nil {
		return res, fmt.Errorf("failed to determine the key existing: %w", err)
	} else if exists == 0 {
		return res, ErrNotFound
	}

	if err := hgetallCmd.Scan(&res); err != nil {
		return res, fmt.Errorf("failed to scan value: %w", err)
	}

	return res, nil
}

func (s *RedisStorage) DeleteById(ctx context.Context, id string) error {
	if _, err := s.rdb.Del(ctx, lockId(id).key()).Result(); err != nil && err != redis.Nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

func (s *RedisStorage) Save(ctx context.Context, l models.Lock) error {
	key := lockId(l.Id).key()

	p := s.rdb.TxPipeline()
	p.HSet(ctx, key, l)
	p.Expire(ctx, key, s.ttl)

	if _, err := p.Exec(ctx); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	return nil
}

type lockId string

// Возвращает значение ключа для данной блокировки
func (id lockId) key() string {
	return fmt.Sprintf("lock:%s", id)
}
