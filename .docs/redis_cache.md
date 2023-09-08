### Biblioteca Redis

A biblioteca Redis facilita a interação com um servidor Redis, permitindo armazenar dados em cache, o que pode melhorar significativamente a performance do seu aplicativo.

#### Configuração

Para configurar o cliente do Redis, você deve definir as seguintes variáveis de ambiente:

```sh
REDIS_ADDR=host.docker.internal
REDIS_NAMESPACE=namespace
REDIS_PASSWORD= (opcional)
```

- `REDIS_ADDR`: O endereço do seu servidor Redis.
- `REDIS_NAMESPACE`: Um namespace para diferenciar e organizar melhor as suas chaves.
- `REDIS_PASSWORD`: A senha do seu servidor Redis, se aplicável.

#### Inicialização

Para iniciar a conexão com o servidor Redis, use o código abaixo:

```go
package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/redis"
)

func main() {
	env.LoadEnv()
	redis.StartRedisCache()
}
```

#### Uso da Biblioteca

A biblioteca Redis oferece uma série de métodos que facilitam a interação com o servidor Redis. Abaixo estão os métodos disponíveis e exemplos de como usá-los:

- `redis.Set(key string, value interface{}, expiration time.Duration) error`: Este método permite adicionar um novo valor ao cache. A chave é uma string que identifica o valor, o valor é o próprio dado que você deseja armazenar e a expiração é a duração após a qual o valor será automaticamente removido do cache.

```go
err := redis.Set("nome_da_chave", valor, time.Duration(0))
```

- `redis.Get(key string) (string, error)`: Este método permite recuperar um valor do cache como uma string, usando sua chave.

```go
valor, err := redis.Get("nome_da_chave")
```

- `redis.GetDecode(key string, obj interface{}) error`: Este método permite recuperar um valor do cache e decodificá-lo diretamente em uma variável de um tipo específico.

```go
var app model.Application
err := redis.GetDecode("nome_da_chave", &app)
```

- `redis.Remove(key string) error`: Este método permite remover um valor do cache usando sua chave.

```go
err := redis.Remove("nome_da_chave")
```

Lembre-se de tratar possíveis erros retornados por esses métodos para garantir que seu aplicativo possa lidar adequadamente com situações de falha.