package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	redis "github.com/caiomarcatti12/nanogo/v2/config/redis"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	//Inicializa a conexão com redis
	redis.StartRedisQueue()

	//Inicializa o NewCurlTask na fila my_job
	redis.RedisQueueConsumer(tasks.NewCurlTask(), "my_job")

	//Publica uma mensagem na fila my_job
	args := map[string]interface{}{"url": "http://example.com"}
	redis.Enqueue("CRONTAB_WORKER_CURL", args)

	//Aguarda o sinal de stop para parar a aplicação
	config.WaitSignalStop()
}
