package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"github.com/imboyko/lock-service/internal/config"
	"github.com/imboyko/lock-service/internal/models"
	"github.com/imboyko/lock-service/internal/storage"
)

func TestNewRedisStorage(t *testing.T) {
	t.Parallel()

	t.Run("successful connect", func(t *testing.T) {
		mr := miniredis.RunT(t)
		stor, err := storage.NewRedisStorage(
			context.Background(),
			config.Redis{
				Host: mr.Host(),
				Port: mr.Port(),
			})
		assert.NoError(t, err)
		assert.NoError(t, stor.Close())
	})

	t.Run("failed connect", func(t *testing.T) {
		_, err := storage.NewRedisStorage(
			context.Background(),
			config.Redis{
				Port: "9999",
			})
		assert.Error(t, err)
	})
}

func TestGetters(t *testing.T) {
	mock, stor, _ := createStorage(t)
	t.Cleanup(mock.Close)

	testSet := []models.Lock{
		{
			Id:        "666",
			Username:  "User 1",
			Timestamp: time.Now().Round(time.Second),
		},
		{
			Id:        "777",
			Username:  "User 2",
			Timestamp: time.Now().Round(time.Second),
		},
	}

	for _, l := range testSet {
		key := fmt.Sprintf("lock:%s", l.Id)
		mock.HSet(key, "id", l.Id)
		mock.HSet(key, "timestamp", l.Timestamp.Format(time.RFC3339))
		mock.HSet(key, "username", l.Username)
	}

	t.Run("GetById", func(t *testing.T) {
		exp := testSet[0]
		res, err := stor.GetById(context.Background(), exp.Id)
		if assert.NoError(t, err) {
			assert.Equal(t, exp, res)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		res, err := stor.GetAll(context.Background())
		if assert.NoError(t, err) {
			assert.Len(t, res, len(testSet))
		}
	})
}

func TestDeleteById(t *testing.T) {
	mock, stor, _ := createStorage(t)
	t.Cleanup(mock.Close)

	err := stor.DeleteById(context.Background(), "666")
	assert.NoError(t, err)
}

func createStorage(t *testing.T) (*miniredis.Miniredis, *storage.RedisStorage, error) {
	mr := miniredis.RunT(t)
	stor, err := storage.NewRedisStorage(
		context.Background(),
		config.Redis{
			Host: mr.Host(),
			Port: mr.Port(),
		})
	return mr, stor, err
}
