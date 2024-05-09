/*
  - Copyright 2023 Caio Matheus Marcatti Calimério
    *
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
    *
  - http://www.apache.org/licenses/LICENSE-2.0
    *
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
*/
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/caiomarcatti12/nanogo/v2/config/rsql"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoORM[T interface{}] interface {
	SetCollection(collectionName string)
	Insert(document T) (uuid.UUID, error)
	Update(document T) (bool, error)
	Delete(document T) (bool, error)
	DeleteById(id uuid.UUID) (bool, error)
	FindById(id uuid.UUID) (*T, error)
	FindAll() ([]T, error)
	RawQuery(query bson.M, sort bson.M, limit int64, skip int64) ([]T, int64, error)
	RawQueryParseRsql(filter rsql.QueryFilter) ([]T, int64, error)
}

type MongoORM[T interface{}] struct {
	logger     log.ILog
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoORM[T interface{}](logger log.ILog, connector IMongoConnector) IMongoORM[T] {
	err := connector.ConnectMongoDB()

	if err != nil {
		logger.Error(err.Error())
	}

	return &MongoORM[T]{
		db:     connector.GetClientDB(),
		logger: logger,
	}
}

func (r *MongoORM[T]) SetCollection(collectionName string) {
	r.collection = r.db.Collection(collectionName)
}

func (r *MongoORM[T]) Insert(document T) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := uuid.NewRandom()

	if err != nil {
		r.logger.Error("erro ao gerar UUID: " + err.Error())
		return uuid.Nil, errors.New("erro ao gerar UUID: " + err.Error())
	}

	// Adicionando o _id ao mapa
	err = r.SetID(&document, _id)

	if err != nil {
		r.logger.Error(err.Error())
		return _id, err
	}

	// Inserindo o documento na coleção
	_, err = r.collection.InsertOne(ctx, document)
	if err != nil {
		r.logger.Error(err.Error())
		return _id, err
	}

	return _id, nil
}

func (r *MongoORM[T]) Update(document T) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := r.GetID(&document)

	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	filter := bson.D{{"_id", _id}}
	updateDoc := bson.M{"$set": document}

	_, err = r.collection.UpdateOne(ctx, filter, updateDoc)

	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	return true, nil

}

func (r *MongoORM[T]) Delete(document T) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := r.GetID(&document)

	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	filter := bson.D{{"_id", _id}}

	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	return true, nil
}

func (r *MongoORM[T]) DeleteById(uuid uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", uuid}}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (r *MongoORM[T]) FindById(id uuid.UUID) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	var result T

	err := r.collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}

		r.logger.Error(err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *MongoORM[T]) FindAll() ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	// Verificar se a lista de resultados está vazia
	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func (r *MongoORM[T]) RawQueryParseRsql(filter rsql.QueryFilter) ([]T, int64, error) {
	query, sort, limit, skip, err := RsqlConvertToBson(filter)

	if err != nil {
		return nil, 0, err
	}

	return r.RawQuery(query, sort, limit, skip)
}

func (r *MongoORM[T]) RawQuery(query bson.M, sort bson.M, limit int64, skip int64) ([]T, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configurando as opções de pesquisa
	findOptions := options.Find()
	if sort != nil {
		findOptions.SetSort(sort)
	}
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, 0, err
	}

	return results, count, nil
}

func (r *MongoORM[T]) RawQueryCount(query bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return count, nil
}

func (r *MongoORM[T]) SetID(entity interface{}, id uuid.UUID) error {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Se for um ponteiro, obtém o valor ao qual ele aponta
	}

	if val.Kind() == reflect.Struct {
		idField := val.FieldByName("ID")
		if idField.IsValid() && idField.Type() == reflect.TypeOf(uuid.UUID{}) {
			idField.Set(reflect.ValueOf(id))
			return nil
		}
	}

	return fmt.Errorf("field ID not found, cannot be set, or type mismatch")
}

func (r *MongoORM[T]) GetID(entity interface{}) (uuid.UUID, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Se for um ponteiro, obtém o valor ao qual ele aponta
	}

	// Checa se o valor é uma struct e se o tipo é ModuleEntity
	if val.Kind() == reflect.Struct {
		idField := val.FieldByName("ID")
		if idField.IsValid() && idField.Type() == reflect.TypeOf(uuid.UUID{}) {
			return idField.Interface().(uuid.UUID), nil
		}
	}

	return uuid.Nil, fmt.Errorf("field ID not found or type mismatch")
}
