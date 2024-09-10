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
		instance := NewInstaceMongoDB(env, logger)

		err := instance.Connect()

		if err != nil {
			panic(err)
		}

		return instance
	default:
		panic(fmt.Sprintf("invalid database provider: %s", provider))
	}
}
