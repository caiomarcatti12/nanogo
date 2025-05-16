# Métricas

Integração com Prometheus para coleta de métricas de aplicações.

## Principais APIs
- `Factory(env, logger) IMetric` cria o provider.
- Métodos para criar e manipular `Counter`, `Gauge`, `Histogram` e `Summary`.

## Variáveis de Ambiente
- `METRIC_PROVIDER` (ex: `PROMETHEUS`).
- `PROMETHEUS_PREFIX` define prefixo dos nomes.

## Exemplo de Uso
```go
metric := metric.Factory(envAdapter, logger)
metric.CreateMetric(metric.Counter, "requests", "Total", []string{"path"})
metric.IncrementCounter("requests", metric.Labels{"path": "/"})
```
