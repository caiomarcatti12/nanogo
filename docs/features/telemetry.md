# Telemetria

Implementação de rastreamento distribuído com suporte ao OpenTelemetry ou memória.

## Principais APIs
- `Factory(env, log) ITelemetry` cria a implementação configurada.
- Métodos para criar spans: `CreateRootSpan`, `StartChildSpan`, `EndSpan` e `Shutdown`.

## Variáveis de Ambiente
- `TELEMETRY_ENABLE` (`true` ou `false`).
- `TELEMETRY_DISPATCHER` (`OPEN_TELEMETRY` ou `MEMORY`).
- `TELEMETRY_ENDPOINT` endereço do collector.

## Exemplo de Uso
```go
telemetry := telemetry.Factory(envAdapter, logger)
span := telemetry.CreateRootSpan("process")
defer telemetry.EndSpan(span, nil)
```
