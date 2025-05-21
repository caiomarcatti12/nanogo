# API WebServer

O pacote `webserver` disponibiliza uma estrutura simples para criar APIs HTTP utilizando o roteador `gorilla/mux`. Ele integra automaticamente serviços de logging, internacionalização e telemetria, além de permitir o registro dinâmico de rotas e middlewares através do sistema de Dependency Injection (DI) do Nanogo.

## Estrutura

O pacote é composto pelos seguintes elementos principais:

- **IWebServer**: interface que define os métodos públicos do servidor.
- **WebServer**: implementação padrão exposta pelo factory.
- **Middlewares**: conjunto de middlewares já prontos, como CORS, extração de payload e correlação de logs.
- **Route**: estrutura utilizada para definir rotas (método, caminho e handler).

## Funcionalidades

- Inicialização simplificada via `Factory`/DI.
- Registro de rotas de forma tipada e com injeção automática de dependências.
- Inclusão de middlewares customizados ou dos já fornecidos pelo framework.
- Rotas de health check (`/healthz/livez`, `/healthz/readyz` e `/healthz/startupz`) configuradas por padrão.

## Uso Básico

### 1. Inicializando o servidor

```go
package main

import (
    "github.com/caiomarcatti12/nanogo/pkg/di"
    "github.com/caiomarcatti12/nanogo/pkg/nanogo"
    "github.com/caiomarcatti12/nanogo/pkg/webserver"
)

func main() {
    // prepara o container de dependências e serviços padrões
    nanogo.Bootstrap()

    ws, err := di.Get[webserver.IWebServer]()
    if err != nil {
        panic(err)
    }

    ws.Start()
}
```

### 2. Definindo rotas

```go
import (
    "net/http"
    "github.com/caiomarcatti12/nanogo/pkg/webserver/types"
)

func init() {
    ws, _ := di.Get[webserver.IWebServer]()

    ws.AddRoute(types.Route{
        Method:      http.MethodGet,
        Path:        "/hello/{name}",
        IHandler:    NewHelloController,
        HandlerFunc: "SayHello",
    })
}
```

### 3. Implementando handlers

```go
// HelloPayload mapeia dados de corpo, rota ou query string
type HelloPayload struct {
    Name string `json:"name"`
}

type HelloController struct{}

func NewHelloController() *HelloController { return &HelloController{} }

// SayHello será chamado de acordo com a rota registrada
func (hc *HelloController) SayHello(p HelloPayload) (interface{}, error) {
    return map[string]string{"message": "Hello " + p.Name}, nil
}
```

### 4. Middlewares

O servidor já inicia com os seguintes middlewares padrões:

- **CorsMiddleware** – configuração de CORS via variáveis de ambiente.
- **PayloadExtractorMiddleware** – extrai parâmetros de rota, query, body e uploads.
- **CorrelationIdMiddleware** – adiciona `X-Correlation-ID` às requisições.
- **TelemetryMiddleware** – cria spans de telemetria quando habilitado.

Novos middlewares podem ser adicionados através de `AddMidleware`.

## Variáveis de Ambiente

| Variável                       | Descrição                                               | Default |
|-------------------------------|---------------------------------------------------------|---------|
| WEB_SERVER_HOST               | Endereço de bind do servidor                            | `""`   |
| WEB_SERVER_PORT               | Porta do servidor                                       | `8080` |
| WEB_SERVER_CERTIFICATE        | Caminho do certificado TLS (habilita HTTPS)             | `""`   |
| WEB_SERVER_KEY                | Caminho da chave TLS                                    | `""`   |
| WEB_SERVER_MAX_UPLOAD_SIZE    | Tamanho máximo (MB) para uploads multipart              | `5`    |
| WEBSERVER_ORIGINS             | Lista de origens permitidas para CORS                   | `"*"`  |
| WEBSERVER_HEADERS             | Cabeçalhos permitidos para CORS                         | `"Content-Type"` |
| WEBSERVER_METHODS             | Métodos permitidos para CORS                            | `"GET,POST,PUT,DELETE"` |
| WEBSERVER_ACCESS_LOG          | Exibe log de entrada das requisições                    | `false` |

## Métodos Principais

- `AddMidleware(m middleware.IMiddleware)`: registra um middleware na cadeia de execução.
- `AddRoute(route types.Route)`: adiciona uma nova rota ao servidor.
- `Start()`: inicia o servidor utilizando HTTP ou HTTPS dependendo dos certificados.

## Testes Automatizados

```shell
go test ./...
```

## Boas Práticas

- Utilize `nanogo.Bootstrap()` para garantir que todos os serviços necessários estejam configurados.
- Mantenha handlers pequenos e focados em regras de negócio, delegando parsing de dados ao `PayloadExtractorMiddleware`.
- Defina variáveis de ambiente para personalizar o comportamento do servidor conforme o ambiente de execução.

