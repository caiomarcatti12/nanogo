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
	"fmt"

	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

type OperationType string

type IDatabase interface {
	Connect() error
	GetClient() interface{}
	Disconnect()
}

func Factory(env env.IEnv, logger log.ILog) IDatabase {
	provider := env.GetEnv("DATABASE_PROVIDER", "MONGODB")

	switch provider {
	case "MONGODB":
		credential := MongoCredential{
			protocol:   env.GetEnv("MONGO_PROTOCOL", "mongodb"),
			dbAuthName: env.GetEnv("MONGO_AUTH_DBNAME", "admin"),
			username:   env.GetEnv("MONGO_USERNAME", ""),
			password:   env.GetEnv("MONGO_PASSWORD", ""),
			host:       env.GetEnv("MONGO_HOST", ""),
			port:       env.GetEnv("MONGO_PORT", "27017"),
			database:   env.GetEnv("MONGO_DATABASE", ""),
			uri:        env.GetEnv("MONGO_URI", ""),
		}
		instance := NewInstaceMongoDB(credential, logger)

		err := instance.Connect()

		if err != nil {
			panic(err)
		}

		return instance
	case "CLICKHOUSE":
		credential := ClickhouseCredential{
			Addr:     env.GetEnv("CLICKHOUSE_ADDR", "localhost:9000"),
			Username: env.GetEnv("CLICKHOUSE_USERNAME", "default"),
			Password: env.GetEnv("CLICKHOUSE_PASSWORD", ""),
			Database: env.GetEnv("CLICKHOUSE_DATABASE", "default"),
		}
		instance := NewInstanceClickhouse(credential, logger)

		if err := instance.Connect(); err != nil {
			panic(err)
		}

		return instance
	default:
		panic(fmt.Sprintf("invalid database provider: %s", provider))
	}
}
