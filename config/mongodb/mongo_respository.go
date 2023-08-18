package mongodb

import (
	"github.com/google/uuid"
)

// Repository defines the contract for the MongoDB repository.
type MongoRepository interface {
	Insert(document Model) (Model, error)
	Update(document Model) (Model, error)
	Delete(document Model) (bool, error)
	Save(document Model) (Model, error)
	FindById(id *uuid.UUID) (Model, error)
	FindAll() ([]Model, error)
	RawQuery(query bson.M) ([]Model, error)
}
