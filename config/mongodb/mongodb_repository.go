package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/caiomarcatti12/nanogo/v2/config/repository"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository[T repository.Model] struct {
	collection *mongo.Collection
	model      T
}

func NewMongoRepository[T repository.Model](collectionName string, model T) *MongoRepository[T] {
	db, err := ConnectMongoDB()

	if err != nil {
		log.Fatal(err)
	}

	collection := db.Collection(collectionName)
	return &MongoRepository[T]{collection: collection, model: model}
}

func (r *MongoRepository[T]) Insert(document T) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uuid, err := uuid.NewRandom()
	if err != nil {
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), errors.New("erro ao gerar UUID: " + err.Error())
	}

	document.SetID(&uuid)
	_, err = r.collection.InsertOne(ctx, document)
	if err != nil {
		log.Errorf("Erro ao inserir documento: %v", err)
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), err
	}

	return document, nil
}

func (r *MongoRepository[T]) Update(document T) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", document.GetID()}}
	updateDoc := bson.M{"$set": document}
	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		log.Errorf("Erro ao atualizar documento: %v", err)
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), err
	}

	return document, nil
}

func (r *MongoRepository[T]) Delete(document T) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", document.GetID()}}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Erro ao deletar documento: %v", err)
		return false, err
	}

	return true, nil
}

func (r *MongoRepository[T]) DeleteById(uuid uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", uuid}}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Erro ao deletar documento: %v", err)
		return false, err
	}

	return true, nil
}

func (r *MongoRepository[T]) Save(document T) (T, error) {
	if document.GetID() == nil {
		return r.Insert(document)
	} else {
		return r.Update(document)
	}
}

func (r *MongoRepository[T]) FindById(id *uuid.UUID) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	var result map[string]interface{}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Errorf("Erro ao encontrar documento pelo ID: %v", err)
		return reflect.Zero(reflect.TypeOf(r.model)).Interface().(T), err
	}

	document, err := r.decodeAndMap(result)
	if err != nil {
		return reflect.Zero(reflect.TypeOf(r.model)).Interface().(T), err
	}

	return document, nil
}

func (r *MongoRepository[T]) FindAll() ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		log.Errorf("Erro ao buscar todos os documentos: %v", err)
		return nil, err
	}

	// Verificar se a lista de resultados está vazia
	if len(results) == 0 {
		return []T{}, nil
	}

	var models []T
	for _, result := range results {
		document, err := r.decodeAndMap(result)
		if err != nil {
			return nil, err
		}
		models = append(models, document)
	}

	return models, nil
}

func (r *MongoRepository[T]) RawQuery(query bson.M) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var models []T
	for _, result := range results {
		document, err := r.decodeAndMap(result)
		if err != nil {
			return nil, err
		}
		models = append(models, document)
	}

	return models, nil
}

func (r *MongoRepository[T]) decodeAndMap(result map[string]interface{}) (T, error) {
	if idField, ok := result["_id"].(primitive.Binary); ok {
		convertedUUID, err := uuid.FromBytes(idField.Data)
		if err != nil {
			return reflect.Zero(reflect.TypeOf(r.model)).Interface().(T), err
		}
		result["_id"] = convertedUUID
	}

	var document T
	if err := mapstructure.Decode(result, &document); err != nil {
		return reflect.Zero(reflect.TypeOf(r.model)).Interface().(T), err
	}

	return document, nil
}
