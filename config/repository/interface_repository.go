package repository

import (
	"github.com/google/uuid"
)

type Repository[T any] interface {
	Insert(document *T) (*T, error)
	Update(document *T) (*T, error)
	Delete(document *T) (bool, error)
	FindById(id uuid.UUID) (*T, error)
	DeleteById(id uuid.UUID) (bool, error)
	FindAll() ([]*T, error)
}
