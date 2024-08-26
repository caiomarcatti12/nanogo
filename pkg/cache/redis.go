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
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	"github.com/gomodule/redigo/redis"
)

var (
	once     sync.Once
	instance RedisCache
)

type RedisCache struct {
	redisAddr string
	namespace string
	redisPass string
	pool      redis.Pool
	logger    log.ILog
}

func NewInstanceRedis(env env.IEnv, logger log.ILog) (ICache, error) {
	once.Do(func() {
		redisAddr := env.GetEnv("REDIS_ADDR")
		redisNamespace := env.GetEnv("REDIS_NAMESPACE")
		redisPass := env.GetEnv("REDIS_PASSWORD", "")

		instance = RedisCache{
			redisAddr: redisAddr,
			namespace: redisNamespace,
			redisPass: redisPass,
			logger:    logger,
		}
	})

	return &instance, nil

}

func (r *RedisCache) Connect() error {

	r.logger.Info("Connecting to Redis...")

	r.pool = redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				r.redisAddr,
				redis.DialPassword(r.redisPass),
			)
		},
	}

	return nil
}

func (r *RedisCache) Get(key string) (string, error) {
	conn := r.pool.Get()

	value, err := redis.String(conn.Do("GET", key))

	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RedisCache) GetDecoded(key string, dest interface{}) error {
	value, err := r.Get(key)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(value), &dest)

	if err != nil {
		return fmt.Errorf("error while trying to decode value: %w", err)
	}

	return nil
}

func (r *RedisCache) Set(key string, value interface{}, ttl ...int) error {
	conn := r.pool.Get()

	valueStr, err := r.stringifyValue(value)

	if err != nil {
		return fmt.Errorf("error while trying to stringify value: %w", err)
	}

	_, err = conn.Do("SET", key, valueStr)

	if err != nil {
		return fmt.Errorf("error while trying to set key: %w", err)
	}

	if len(ttl) > 0 && ttl[0] > 0 {
		_, err = conn.Do("EXPIRE", key, ttl[0])

		if err != nil {
			return fmt.Errorf("error while trying to set ttl: %w", err)
		}
	}

	return nil
}

func (r *RedisCache) Remove(key string) error {
	conn := r.pool.Get()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("error while trying to remove key: %w", err)
	}

	return nil
}

func (r *RedisCache) Disconnect() {
	r.pool.Close()
}

func (r *RedisCache) stringifyValue(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", errors.New("failed to stringify value: " + err.Error())
	}
	return string(bytes), nil
}
