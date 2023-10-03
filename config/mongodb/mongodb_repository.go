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
package mongodb

import (
	"context"
	"errors"
	"github.com/caiomarcatti12/nanogo/v2/config/rsql"
	"github.com/caiomarcatti12/nanogo/v2/config/util"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/caiomarcatti12/nanogo/v2/config/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository[T repository.Identifier] struct {
	collection mongo.Collection
	model      T
}

func NewMongoRepository[T repository.Identifier](collectionName string, model T) MongoRepository[T] {
	db, err := ConnectMongoDB()

	if err != nil {
		log.Fatal(err)
	}

	collection := db.Collection(collectionName)
	return MongoRepository[T]{collection: *collection, model: model}
}

func (r *MongoRepository[T]) Insert(document T) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uuid, err := uuid.NewRandom()
	if err != nil {
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), errors.New("erro ao gerar UUID: " + err.Error())
	}

	document.SetID(uuid)
	cp := document
	_, err = r.collection.InsertOne(ctx, cp)
	if err != nil {
		log.Errorf("Erro ao inserir documento: %v", err)
		return reflect.Zero(reflect.TypeOf(document)).Interface().(T), err
	}

	return document, nil
}

func (r *MongoRepository[T]) Update(document T) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", document.GetID()}}
	updateDoc := bson.M{"$set": document}
	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		log.Errorf("Erro ao atualizar documento: %v", err)
		return false, err
	}

	return true, nil
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
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Erro ao deletar documento: %v", err)
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (r *MongoRepository[T]) FindById(id uuid.UUID) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	var result T
	err := r.collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return util.ZeroReflectGeneric(r.model), nil
		}

		log.Errorf("Erro ao encontrar documento pelo ID: %v", err)
		return util.ZeroReflectGeneric(r.model), err
	}

	return result, nil
}

func (r *MongoRepository[T]) FindAll() ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		log.Errorf("Erro ao buscar todos os documentos: %v", err)
		return nil, err
	}

	// Verificar se a lista de resultados está vazia
	if len(results) == 0 {
		return []T{}, nil
	}

	return results, nil
}

func (r *MongoRepository[T]) RawQueryParseRsql(filter rsql.QueryFilter) ([]T, int64, error) {
	query, sort, limit, skip, err := RsqlConvertToBson(filter)

	if err != nil {
		return nil, 0, err
	}

	return r.RawQuery(query, sort, limit, skip)
}

func (r *MongoRepository[T]) RawQuery(query bson.M, sort bson.M, limit int64, skip int64) ([]T, int64, error) {
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
		log.Errorf("Erro ao contar os documentos: %v", err)
		return nil, 0, err
	}

	return results, count, nil
}

func (r *MongoRepository[T]) RawQueryCount(query bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		log.Errorf("Erro ao contar os documentos: %v", err)
		return 0, err
	}

	return count, nil
}
