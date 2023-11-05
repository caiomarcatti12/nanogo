### Biblioteca MetricManager

A biblioteca `MetricManager` facilita a criação e gerenciamento de métricas dentro de aplicativos Go, integrando-se com o [Prometheus](https://prometheus.io/). Ela permite que você defina e atualize métricas que podem ser expostas e consumidas por um sistema de monitoramento Prometheus.

#### Configuração

Para utilizar o `MetricManager`, configure as seguintes variáveis de ambiente para definir a base das métricas e o endpoint para exposição das métricas:

```sh
ENABLE_PROMETHEUS=false # Habilita ou desabilita a exposição de métricas
PROMETHEUS_PREFIX= # Prefixo opcional para todas as métricas
PROMETHEUS_ROUTE=/metrics # Endpoint para acessar as métricas
```

- `ENABLE_PROMETHEUS`: Define se o endpoint de métricas estará ativo.
- `PROMETHEUS_PREFIX`: Um prefixo opcional que é adicionado a todas as métricas para evitar colisões de nomes.
- `PROMETHEUS_ROUTE`: O caminho HTTP onde as métricas podem ser coletadas pelo Prometheus.

#### Inicialização

Para iniciar o `MetricManager` e permitir a coleta de métricas, use o seguinte exemplo:

```go
package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/metric_manager"
)

func main() {
	env.LoadEnv()
	metrics := metric_manager.NewMetricManager()
	// Configure as métricas conforme necessário aqui
}
```

#### Uso da Biblioteca

A biblioteca `MetricManager` oferece métodos para criar e manipular diferentes tipos de métricas. Abaixo estão os métodos disponíveis e exemplos de como usá-los:

- `CreateMetric(metricType MetricType, name, help string, labelKeys LabelsKeys)`: Este método permite criar uma nova métrica com um tipo específico, nome, descrição e chaves de rótulo associadas.

```go
labelKeys := metric_manager.LabelsKeys{"url", "method"}
metricManager.CreateMetric(metric_manager.Counter, "http_request_count", "Total HTTP request count", labelKeys)
```

- `IncrementCounter(name string, labelValues Labels) error`: Este método incrementa um contador baseado no nome e valores de rótulos fornecidos.

```go
labels := metric_manager.Labels{"url": "/login", "method": "POST"}
err := metricManager.IncrementCounter("http_request_count", labels)
```

- `SetGauge(name string, value float64, labelValues Labels) error`: Este método define o valor de um gauge.

```go
labels := metric_manager.Labels{"host": "server01"}
err := metricManager.SetGauge("concurrent_sessions", 42, labels)
```

- `ObserveHistogram(name string, value float64, labelValues Labels) error`: Este método observa um valor, que será contabilizado em um histograma.

```go
labels := metric_manager.Labels{""endpoint"": "/api/resource"}
err := metricManager.ObserveHistogram("response_sizes", 512, labels)
```

- `ObserveSummary(name string, value float64, labelValues Labels) error`: Este método observa um valor, que será contabilizado em um resumo.

```go
labels := metric_manager.Labels{""endpoint"": "/api/resource"}
err := metricManager.ObserveSummary("request_durations", 0.350, labels)
```

Certifique-se de verificar e tratar erros retornados por esses métodos para garantir a correta monitorização e registro das métricas.
