/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/rsql"
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
	RawQueryParseRsql(filter rsql.QueryFilter) ([]T, int64, error)
	RawQuery(query bson.M, sort bson.M, limit int64, skip int64) ([]T, int64, error)
}

type MongoORM[T interface{}] struct {
	logger     log.ILog
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoORM[T interface{}](db IDatabase, logger log.ILog) IMongoORM[T] {
	if client, ok := db.GetClient().(*mongo.Database); ok {
		return &MongoORM[T]{
			db:     client,
			logger: logger,
		}
	} else {
		panic("error constructing MongoORM: client is not a mongo.Client")
	}
}

func (r *MongoORM[T]) SetCollection(collectionName string) {
	r.collection = r.db.Collection(collectionName)
}

func (r *MongoORM[T]) Insert(document T) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := r.GetID(&document)

	if err != nil {
		r.logger.Error(err.Error())
		return uuid.Nil, err
	}

	if _id == uuid.Nil {
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
	}

	// Inserindo o documento na coleção
	_, err = r.collection.InsertOne(ctx, &document)
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
		// Tenta chamar o método GetID
		getIDMethod := val.MethodByName("GetID")
		if getIDMethod.IsValid() && getIDMethod.Type().NumIn() == 0 && getIDMethod.Type().NumOut() == 1 {
			existingID := getIDMethod.Call(nil)[0].Interface().(uuid.UUID)
			if existingID != (uuid.UUID{}) {
				return nil // ID já está preenchido, não força atualização
			}
		}

		// Tenta chamar o método SetID
		setIDMethod := val.MethodByName("SetID")
		if setIDMethod.IsValid() && setIDMethod.Type().NumIn() == 1 && setIDMethod.Type().In(0) == reflect.TypeOf(uuid.UUID{}) {
			setIDMethod.Call([]reflect.Value{reflect.ValueOf(id)})
			return nil
		}
	}

	return fmt.Errorf("field ID not found, cannot be set, or type mismatch")
}
func (r *MongoORM[T]) GetID(s interface{}) (uuid.UUID, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return uuid.Nil, fmt.Errorf("não é uma struct")
	}

	// Usar o método UnexportedField
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return uuid.Nil, fmt.Errorf("campo 'ID' não encontrado")
	}

	// Acessar o valor do campo não exportado
	id := reflect.NewAt(idField.Type(), unsafe.Pointer(idField.UnsafeAddr())).Elem()

	if id.Type() != reflect.TypeOf(uuid.UUID{}) {
		return uuid.Nil, fmt.Errorf("campo 'ID' não é do tipo uuid.UUID")
	}

	return id.Interface().(uuid.UUID), nil
}

func RsqlConvertToBson(qp rsql.QueryFilter) (bson.M, bson.M, int64, int64, error) {
	query := bson.M{}

	// Verificar se Filter foi fornecido
	if qp.Filter != "" {
		// Convertendo o filtro RSQL para bson.M
		conditions, err := rsql.Parse(qp.Filter)
		if err != nil {
			return nil, nil, 0, 0, err
		}

		for _, condition := range conditions {
			value := parseValue(condition.Value)
			switch condition.Operator {
			case "==":
				query[condition.Field] = value
			case "!=":
				query[condition.Field] = bson.M{"$ne": value}
			case "<>":
				query[condition.Field] = bson.M{"$ne": value}
			case "=gt=":
				query[condition.Field] = bson.M{"$gt": value}
			case "=gte=":
				query[condition.Field] = bson.M{"$gte": value}
			case "=lt=":
				query[condition.Field] = bson.M{"$lt": value}
			case "=lte=":
				query[condition.Field] = bson.M{"$lte": value}
			case "=like=":
				query[condition.Field] = bson.M{"$regex": value, "$options": "i"}
			}
		}
	}

	sort := bson.M{}
	// Verificar se SortParams foi fornecido
	if qp.SortParams != "" {
		// Convertendo o parâmetro de ordenação para bson.M
		sortFields := strings.Split(qp.SortParams, ":")
		if len(sortFields) == 2 {
			field := sortFields[0]
			direction := 0
			switch sortFields[1] {
			case "ASC":
				direction = 1
			case "DESC":
				direction = -1
			}
			if direction != 0 {
				sort[field] = direction
			}
		}
	}

	size := 15
	// Verificar se Size foi fornecido
	if qp.Size != 15 {
		size = int(qp.Size)
	}

	skip := 0
	// Verificar se Skip foi fornecido
	if qp.Skip != 0 {
		skip = int(qp.Skip)
	}

	return query, sort, int64(size), int64(skip), nil
}

func parseValue(value string) interface{} {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	valueUUID, err := uuid.Parse(value)

	if err == nil {
		return valueUUID
	}

	return value
}
