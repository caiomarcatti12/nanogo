# Cache

Gerencia sistemas de cache com suporte atual ao Redis.

## Principais APIs
- `Factory(env, logger) ICache` cria a instância conforme `CACHE_PROVIDER`.
- Métodos de `ICache`: `Connect`, `Get`, `Set`, `Remove`, `Disconnect`.

## Variáveis de Ambiente
- `CACHE_PROVIDER` (ex: `REDIS`)
- `REDIS_ADDR`
- `REDIS_NAMESPACE`
- `REDIS_PASSWORD`

## Exemplo de Uso
```go
cache := cache.Factory(envAdapter, logger)
cache.Connect()
cache.Set("key", "value")
val, _ := cache.Get("key")
```
