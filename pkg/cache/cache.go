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
package cache

import (
	"fmt"

	"github.com/caiomarcatti12/nanogo/v1/pkg/env"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
)

type ICache interface {
	Connect() error
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl ...int) error
	Remove(key string) error
	Disconnect()
}

func Factory(env env.IEnv, logger log.ILog) ICache {
	cacheProvider := env.GetEnv("CACHE_PROVIDER", "REDIS")

	switch cacheProvider {
	case "REDIS":
		instance, err := NewInstanceRedis(env, logger)

		if err != nil {
			panic(err)
		}

		return instance
	default:
		panic(fmt.Sprintf("invalid cache provider: %s", cacheProvider))
	}
}
