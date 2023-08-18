package repository

import (
	"github.com/google/uuid"
)

type Repository interface {
	Insert(document Model) (Model, error)
	Update(document Model) (Model, error)
	Delete(document Model) (bool, error)
	Save(document Model) (Model, error)
	FindById(id *uuid.UUID) (Model, error)
	FindAll() ([]Model, error)
}
