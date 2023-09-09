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
package redis

import (
	"encoding/json"
	"errors"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	logger "github.com/caiomarcatti12/nanogo/v2/config/log"
	"sync"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	redisAddr string
	namespace string
	redisPass string
	mu        sync.Mutex
	connected bool
	pool      *redis.Pool
}

// Variável global para a instância de RedisCache
var cacheInstance *RedisCache

func StartRedisCache() *RedisCache {
	redisAddr := env.GetEnv("REDIS_ADDR")
	redisNamespace := env.GetEnv("REDIS_NAMESPACE")
	redisPass := env.GetEnv("REDIS_PASSWORD", "")

	cacheInstance = &RedisCache{
		redisAddr: redisAddr,
		namespace: redisNamespace,
		redisPass: redisPass,
	}

	cacheInstance.connect()

	return cacheInstance
}

func (cache *RedisCache) connect() {
	if cache.connected {
		return
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.connected {
		return
	}

	cache.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				cache.redisAddr,
				redis.DialPassword(cache.redisPass),
			)
		},
	}

	cache.connected = true
}

func Set(key string, value interface{}, ttl ...int) error {
	conn := cacheInstance.pool.Get()
	defer conn.Close()

	// Convertendo value para string
	valueStr, err := stringifyValue(value)
	if err != nil {
		logger.Fatal("Erro ao converter valor para string:", err)
		return err
	}

	_, err = conn.Do("SET", key, valueStr)
	if err != nil {
		logger.Fatal("Erro ao definir valor no cache:", err)
		return err
	}

	// Se um TTL for fornecido, defina a expiração para a chave
	if len(ttl) > 0 && ttl[0] > 0 {
		_, err = conn.Do("EXPIRE", key, ttl[0])
		if err != nil {
			logger.Fatal("Erro ao definir TTL no cache:", err)
			return err
		}
	}

	return nil
}

func GetDecode[T any](key string, target *T) (*T, error) {
	conn := cacheInstance.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			logger.Fatal("Erro ao obter valor do cache:", err)
			return nil, err
		}
	}

	if value == "" {
		return target, nil
	}

	if err != nil {
		logger.Fatal("Erro ao fazer o decode do json:", err)
		return nil, err
	}

	if err := decodeToStruct(value, &target); err != nil {
		return nil, err
	}
	return target, nil
}

func Get(key string) (string, error) {
	conn := cacheInstance.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		logger.Fatal("Erro ao obter valor do cache:", err)
		return "", err
	}

	return value, nil
}
func Remove(key string) error {
	conn := cacheInstance.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		logger.Fatal("Erro ao remover a chave do cache:", err)
		return err
	}

	return nil
}

// decodeToStruct tenta decodificar uma string em um struct usando mapperstruct.Decode.
func decodeToStruct[T any](value string, target *T) error {

	err := json.Unmarshal([]byte(value), &target)
	if err != nil {
		return errors.New("failed to decode string to json: " + err.Error())
	}

	return nil
}

// stringifyValue tenta converter um valor do tipo interface{} para uma string usando JSON
func stringifyValue(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", errors.New("failed to stringify value: " + err.Error())
	}
	return string(bytes), nil
}
