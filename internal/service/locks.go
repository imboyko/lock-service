package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/imboyko/lock-service/internal/models"
)

type LockStorage interface {
	// Возвращает все установленные блокировки
	GetAll(context.Context) ([]models.Lock, error)
	// Возвращает блокировку по идентификатору
	GetById(context.Context, string) (models.Lock, error)
	// Удаляет блокировку по идентификатору
	DeleteById(string) error
	// Добавляет или обновляет блокировку
	Save(models.Lock)
}

type LockService struct {
	store LockStorage
}

func (s *LockService) AddOrUpdate(ctx context.Context, l models.Lock) error {
	return errors.New("Not yet implemented")
}

func (s *LockService) GetAllLocks(ctx context.Context) ([]models.Lock, error) {
	return nil, errors.New("Not yet implemented")
}

func (s *LockService) GetLockById(ctx context.Context, id string) ([]models.Lock, error) {
	return nil, errors.New("Not yet implemented")
}

func (s *LockService) DeleteLockById(ctx context.Context, id string) error {
	return errors.New("Not yet implemented")
}

type ErrLockNotFound string

func (id ErrLockNotFound) Error() string {
	return fmt.Sprintf("блокировка с идентификатором %s не найдена", string(id))
}
