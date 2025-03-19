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
	"sync"

	"github.com/caiomarcatti12/nanogo/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once

type MongoClient struct {
	client     *mongo.Client
	clientDB   *mongo.Database
	logger     log.ILog
	credential MongoCredential
}

type MongoCredential struct {
	protocol   string
	dbAuthName string
	username   string
	password   string
	host       string
	port       string
	database   string
	uri        string
}

func NewInstaceMongoDB(credential MongoCredential, logger log.ILog) IDatabase {
	return &MongoClient{
		logger:     logger,
		credential: credential,
	}
}

func (m *MongoClient) Connect() error {
	var connErr error

	once.Do(func() {
		var clientOptions *options.ClientOptions

		m.logger.Trace("Connecting to MongoDB...")

		if m.credential.uri != "" {
			clientOptions = options.Client().ApplyURI(m.credential.uri)
		} else {
			connectionURI := fmt.Sprintf("%s://%s:%s@%s:%d", m.credential.protocol, m.credential.username, m.credential.password, m.credential.host, m.credential.port)

			if connectionURI == "://:@:" {
				connErr = errors.New("Connection URI cannot be empty.")
				return
			}
			clientOptions = options.Client().ApplyURI(connectionURI).SetAuth(options.Credential{AuthSource: m.credential.dbAuthName})
		}

		m.client, connErr = mongo.Connect(context.Background(), clientOptions)

		if connErr != nil {
			return
		}

		m.clientDB = m.client.Database(m.credential.database)
		m.logger.Trace("Connection to MongoDB established!")
	})

	if connErr != nil {
		m.logger.Error(connErr.Error())
	}

	return connErr
}

func (m *MongoClient) GetClient() interface{} {
	return m.clientDB
}

func (m *MongoClient) Disconnect() {
	m.logger.Info("Disconnecting from MongoDB...")
	m.client.Disconnect(context.Background())
}
