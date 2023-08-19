package repository

import (
	"github.com/google/uuid"
)

type Repository interface {
	Insert(document interface{}) (interface{}, error)
	Update(document interface{}) (interface{}, error)
	Delete(document interface{}) (bool, error)
	Save(document interface{}) (interface{}, error)
	FindById(id *uuid.UUID) (interface{}, error)
	FindAll() ([]interface{}, error)
}
