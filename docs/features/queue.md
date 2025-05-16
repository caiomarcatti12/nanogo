# Filas

Gerenciamento de filas com suporte ao RabbitMQ.

## Principais APIs
- `Factory(env, logger, metric, telemetry) IQueue` seleciona o provider.
- Métodos de `IQueue`: `Connect`, `Configure`, `Consume`, `Publish`, `Disconnect`.

## Variáveis de Ambiente
- `QUEUE_PROVIDER` (ex: `RABBITMQ`).
- `RABBITMQ_PROTOCOL`, `RABBITMQ_USER`, `RABBITMQ_PASSWORD`, `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_VHOST`.

## Exemplo de Uso
```go
queue := queue.Factory(envAdapter, logger, metric, telemetry)
queue.Connect()
```
