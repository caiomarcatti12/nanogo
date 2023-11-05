### Biblioteca MetricManager

A biblioteca `MetricManager` facilita a criação e gerenciamento de métricas dentro de aplicativos Go, integrando-se com o [Prometheus](https://prometheus.io/). Ela permite que você defina e atualize métricas que podem ser expostas e consumidas por um sistema de monitoramento Prometheus.

#### Configuração

Para utilizar o `MetricManager`, configure as seguintes variáveis de ambiente para definir a base das métricas e o endpoint para exposição das métricas:

```sh
ENABLE_PROMETHEUS=false # Habilita ou desabilita a exposição de métricas
PROMETHEUS_PREFIX= # Prefixo opcional para todas as métricas
PROMETHEUS_ROUTE=/metrics # Endpoint para acessar as métricas
PROMETHEUS_PORT=9090 # Porta onde o servidor Prometheus estará ouvindo
```

- `ENABLE_PROMETHEUS`: Define se o endpoint de métricas estará ativo.
- `PROMETHEUS_PREFIX`: Um prefixo opcional que é adicionado a todas as métricas para evitar colisões de nomes.
- `PROMETHEUS_ROUTE`: O caminho HTTP onde as métricas podem ser coletadas pelo Prometheus.
- `PROMETHEUS_PORT`: A porta em que o servidor de métricas será exposto.

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
metricManager.CreateMetric(metric_manager.Counter, "http_request_count", "Total HTTP request count", []string{"url", "method"})
```

- `IncrementCounter(name string, labelValues Labels) error`: Este método incrementa um contador baseado no nome e valores de rótulos fornecidos.

```go
err := metricManager.IncrementCounter("http_request_count", map[string]string{"url": "/api/resource", "method": "GET"})
```

- `SetGauge(name string, value float64, labelValues Labels) error`: Este método define o valor de um gauge.

```go
err := metricManager.SetGauge("concurrent_sessions", 42, map[string]string{"host": "server01"})
```

- `ObserveHistogram(name string, value float64, labelValues Labels) error`: Este método observa um valor, que será contabilizado em um histograma.

```go
err := metricManager.ObserveHistogram("response_sizes", 512, map[string]string{"endpoint": "/api/resource"})
```

- `ObserveSummary(name string, value float64, labelValues Labels) error`: Este método observa um valor, que será contabilizado em um resumo.

```go
err := metricManager.ObserveSummary("request_durations", 0.350, map[string]string{"endpoint": "/api/resource"})
```

Certifique-se de verificar e tratar erros retornados por esses métodos para garantir a correta monitorização e registro das métricas.
