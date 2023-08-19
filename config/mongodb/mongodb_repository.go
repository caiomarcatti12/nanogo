package mongodb

import (
	"context"
	"reflect"
	"time"

	"github.com/caiomarcatti12/nanogo/v2/config/repository"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
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

	uuid, _ := uuid.NewRandom()
	document.SetID(&uuid)
	_, err := r.collection.InsertOne(ctx, document)

	if err != nil {
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), err
	}

	return document, err
}

func (r *MongoRepository[T]) Update(document T) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", document.GetID()}}
	updateDoc := bson.M{"$set": document}
	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
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
		return reflect.Zero(reflect.TypeOf(result)).Interface().(T), err
	}

	// Convert _id field back to uuid.UUID
	if idField, ok := result["_id"].(primitive.Binary); ok {
		convertedUUID, err := uuid.FromBytes(idField.Data)
		if err != nil {
			return reflect.Zero(reflect.TypeOf(result)).Interface().(T), err
		}
		result["id"] = convertedUUID
	}

	// We need a fresh instance for each document.
	outputModel := reflect.New(reflect.TypeOf(r.model).Elem()).Interface().(T)
	err = mapstructure.Decode(result, &outputModel)
	if err != nil {
		return reflect.Zero(reflect.TypeOf(result)).Interface().(T), err
	}

	return outputModel, nil
}

func (r *MongoRepository[T]) FindAll() ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	var outputModels []interface{}

	for _, result := range results {

		if idField, ok := result["_id"].(primitive.Binary); ok {
			convertedUUID, err := uuid.FromBytes(idField.Data)
			if err != nil {
				return nil, err
			}
			result["id"] = convertedUUID
		}

		model := reflect.New(reflect.TypeOf(r.model).Elem()).Interface().(T)
		err = mapstructure.Decode(result, &model)
		if err != nil {
			return nil, err
		}

		outputModels = append(outputModels, model)
	}

	return outputModels, nil
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

	var outputModels []T

	for _, result := range results {
		if idField, ok := result["_id"].(primitive.Binary); ok {
			convertedUUID, err := uuid.FromBytes(idField.Data)
			if err != nil {
				return nil, err
			}
			result["id"] = convertedUUID
		}

		model := reflect.New(reflect.TypeOf(r.model).Elem()).Interface().(T)
		err = mapstructure.Decode(result, &model)
		if err != nil {
			return nil, err
		}

		outputModels = append(outputModels, model)
	}

	return outputModels, nil
}

func decode(document interface{}) (map[string]interface{}, error) {
	var mapInterface map[string]interface{}
	err := mapstructure.Decode(document, &mapInterface)
	if err != nil {
		return nil, err
	}
	return mapInterface, nil
}

func encode(inputMap map[string]interface{}, outputModel interface{}) error {
	return mapstructure.Decode(inputMap, outputModel)
}
