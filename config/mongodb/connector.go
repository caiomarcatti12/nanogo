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
	"fmt"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	instance *MongoClient
	once     sync.Once
	err      error
)

type IMongoConnector interface {
	ConnectMongoDB() error
	GetClientDB() *mongo.Database
}

type MongoClient struct {
	client     *mongo.Client
	clientDB   *mongo.Database
	logger     log.ILog
	dbAuthName string
	username   string
	password   string
	host       string
	port       string
	database   string
	uri        string
	connected  bool
}

func NewMongoConnector(env env.IEnv, logger log.ILog) IMongoConnector {
	once.Do(func() {
		instance = &MongoClient{
			client:     nil,
			logger:     logger,
			dbAuthName: env.GetEnv("MONGO_AUTH_DBNAME", "admin"),
			username:   env.GetEnv("MONGO_USERNAME", ""),
			password:   env.GetEnv("MONGO_PASSWORD", ""),
			host:       env.GetEnv("MONGO_HOST", ""),
			port:       env.GetEnv("MONGO_PORT", "27017"),
			database:   env.GetEnv("MONGO_DATABASE", ""),
			uri:        env.GetEnv("MONGO_URI", ""),
			connected:  false,
		}
	})

	return instance
}

func (m *MongoClient) ConnectMongoDB() error {
	var clientOptions *options.ClientOptions

	if m.connected {
		return nil
	}

	m.logger.Info("Conectando no mongoDB...")
	if m.uri != "" {
		clientOptions = options.Client().ApplyURI(m.uri)
	} else {
		connectionURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", m.username, m.password, m.host, m.port)
		clientOptions = options.Client().ApplyURI(connectionURI).SetAuth(options.Credential{AuthSource: m.dbAuthName})
	}

	m.client, err = mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		m.logger.Error(err.Error())
		return err
	}

	m.clientDB = m.client.Database(m.database)
	m.logger.Info("Conexão estabelecida com sucesso")
	m.connected = true

	return nil
}

func (m *MongoClient) GetClientDB() *mongo.Database {
	return m.clientDB
}
