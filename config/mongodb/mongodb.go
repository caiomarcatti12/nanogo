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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

var (
	clientInstance *mongo.Client
	once           sync.Once
	err            error
)

func ConnectMongoDB() (*mongo.Database, error) {
	once.Do(func() {
		mongoURI := env.GetEnv("MONGO_URI", "")
		var clientOptions *options.ClientOptions

		if mongoURI != "" {
			clientOptions = options.Client().ApplyURI(mongoURI)
		} else {
			dbnameAuth := env.GetEnv("MONGO_AUTH_DBNAME", "admin")
			username := env.GetEnv("MONGO_USERNAME")
			password := env.GetEnv("MONGO_PASSWORD")
			host := env.GetEnv("MONGO_HOST")
			port := env.GetEnv("MONGO_PORT", "27017")

			connectionURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
			clientOptions = options.Client().ApplyURI(connectionURI).SetAuth(options.Credential{AuthSource: dbnameAuth})
		}

		clientInstance, err = mongo.Connect(context.Background(), clientOptions)

		log.Info("Connected to MongoDB!")

	})

	if err != nil {
		return nil, err
	}

	dbname := env.GetEnv("MONGO_DBNAME")

	return clientInstance.Database(dbname), nil
}
