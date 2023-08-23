package main

import (
	"fmt"
	"github.com/caiomarcatti12/nanogo/v2/config/redis"
	"time"
)

func main() {
	// Iniciar a instância de cache do Redis
	redis.StartRedisCache()

	// Definir valor no cache sem TTL
	err := redis.Set("nome", "Alice")
	if err != nil {
		fmt.Println("Erro ao definir 'nome':", err)
	}

	// Definir valor no cache com TTL de 5 segundos
	err = redis.Set("tempNome", "Bob", 5)
	if err != nil {
		fmt.Println("Erro ao definir 'tempNome':", err)
	}

	// Obter valor do cache
	nome, err := redis.Get("nome")
	if err != nil {
		fmt.Println("Erro ao obter 'nome':", err)
	} else {
		fmt.Println("Nome:", nome)
	}

	// Aguarde 6 segundos e tente obter o 'tempNome' (deve expirar após 5 segundos)
	time.Sleep(6 * time.Second)
	tempNome, err := redis.Get("tempNome")
	if err != nil {
		fmt.Println("Erro ao obter 'tempNome':", err)
	} else {
		fmt.Println("TempNome:", tempNome)
	}
}
