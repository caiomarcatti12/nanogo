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
package queue

//  import (
// 	 "github.com/caiomarcatti12/nanogo/v3/pkg/env"
// 	 logger "github.com/caiomarcatti12/nanogo/v3/pkg/log"

// 	 "sync"

// 	 "github.com/gocraft/work"
// 	 "github.com/gomodule/redigo/redis"
//  )

//  type RedisQueueConsumerInterface interface {
// 	 Run(job *work.Job) error
//  }

//  type RedisConsumer struct {
// 	 Queue *RedisQueue
//  }

//  func RedisQueueConsumer(consumer RedisQueueConsumerInterface, queueName string) *RedisConsumer {
// 	 redisQueue := RedisQ

// 	 // Registre a função de consumo
// 	 redisQueue.WorkerPool.Job(queueName, consumer.Run)

// 	 // Inicia o processamento de jobs
// 	 redisQueue.WorkerPool.Start()

// 	 return &RedisConsumer{
// 		 Queue: redisQueue,
// 	 }
//  }

//  type RedisQueue struct {
// 	 Enqueuer   *work.Enqueuer
// 	 WorkerPool *work.WorkerPool
// 	 redisAddr  string
// 	 namespace  string
// 	 redisPass  string
// 	 mu         sync.Mutex // Mutex para sincronização de threads
// 	 connected  bool       // Marcador de conexão
//  }

//  // Variável global para a instância de RedisQueue
//  var RedisQ *RedisQueue

//  func StartRedisQueue() {
// 	 redisAddr := env.GetEnv("REDIS_ADDR")
// 	 redisNamespace := env.GetEnv("REDIS_NAMESPACE")
// 	 redisPass := env.GetEnv("REDIS_PASSWORD", "")

// 	 RedisQ = &RedisQueue{
// 		 redisAddr: redisAddr,
// 		 namespace: redisNamespace,
// 		 redisPass: redisPass,
// 	 }

// 	 RedisQ.connect()
//  }

//  func (rq *RedisQueue) connect() {
// 	 if rq.connected {
// 		 return
// 	 }

// 	 rq.mu.Lock()
// 	 defer rq.mu.Unlock()

// 	 if rq.connected {
// 		 return
// 	 }

// 	 redisPool := &redis.Pool{
// 		 Dial: func() (redis.Conn, error) {
// 			 return redis.Dial(
// 				 "tcp",
// 				 rq.redisAddr,
// 				 redis.DialPassword(rq.redisPass),
// 			 )
// 		 },
// 	 }

// 	 rq.Enqueuer = work.NewEnqueuer(rq.namespace, redisPool)
// 	 rq.WorkerPool = work.NewWorkerPool(struct{}{}, 10, rq.namespace, redisPool)

// 	 rq.connected = true
//  }

//  func Enqueue(queueName string, params map[string]interface{}) bool {

// 	 job, err := RedisQ.Enqueuer.Enqueue(queueName, work.Q(params))

// 	 if err != nil {
// 		 logger.Fatal("Erro ao enfileirar tarefa:", err)
// 		 return false
// 	 }

// 	 logger.Trace("Tarefa enfileirada com sucesso. ID do job:", job.ID)

// 	 return true
//  }
