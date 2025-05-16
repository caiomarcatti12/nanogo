# Banco de Dados

Camada de abstração para conexão a bancos, atualmente com suporte ao MongoDB.

## Principais APIs
- `Factory(env, logger) IDatabase` seleciona o provider usando `DATABASE_PROVIDER`.
- Interfaces `IDatabase` e `MongoORM` para operações comuns.

## Variáveis de Ambiente
- `DATABASE_PROVIDER` (ex: `MONGODB`)
- `MONGO_PROTOCOL`
- `MONGO_AUTH_DBNAME`
- `MONGO_USERNAME`
- `MONGO_PASSWORD`
- `MONGO_HOST`
- `MONGO_PORT`
- `MONGO_DATABASE`
- `MONGO_URI`

## Exemplo de Uso
```go
db := db.Factory(envAdapter, logger)
client := db.GetClient()
```
