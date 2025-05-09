## 📦 `pkg/env`: Gerenciador de Variáveis de Ambiente

O pacote `env` é responsável por carregar e gerenciar as **variáveis de ambiente do sistema** de forma centralizada, segura e extensível. Ele permite injetar configurações no runtime e oferece suporte a múltiplos _providers_ (fontes).

---

### ⚙️ Provedores Suportados (`ENV_PROVIDER`)

Você pode configurar qual provider será utilizado por meio da variável `ENV_PROVIDER`.

| `ENV_PROVIDER` | Descrição                                                                 |
|----------------|---------------------------------------------------------------------------|
| `ENV_FILE`     | Carrega um arquivo `.env` de forma local (recomendado para desenvolvimento) |
| `AWS`          | Integra com o AWS Secrets Manager para buscar variáveis remotamente       |
| `OS`           | Usa as variáveis já disponíveis no ambiente do sistema operacional        |

---

### 🧠 Interface `IEnv`

```go
type IEnv interface {
	GetEnv(variable string, default_ ...string) string
	GetEnvBool(variable string, default_ ...string) bool
}
```

---

Excelente sugestão, Caio! Aqui está a **seção `📋 API de Uso` revisada**, com exemplos **tanto do uso direto** quanto via **inversão de dependência** (injeção no estilo Clean Architecture com o `di.Container` do seu projeto):

---

### 📋 API de Uso

#### ✅ Forma direta (modo mais simples)

Para obter o valor de uma variável:

```go
env := env.Factory(i18n) // inicializa com i18n já carregado

port := env.GetEnv("WEB_SERVER_PORT", "8080")
debug := env.GetEnvBool("DEBUG_MODE", "true")
```

---

#### 🧱 Com Inversão de Dependência (via DI Container)

O projeto `nanogo` usa o padrão de Injeção de Dependência com `di.IContainer`. Para obter a instância da interface `env.IEnv` corretamente resolvida:

```go
import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
)

func MyServiceConstructor(env env.IEnv) *MyService {
	return &MyService{env: env}
}
```

E no consumo:

```go
type MyService struct {
	env env.IEnv
}

func (s *MyService) DoSomething() {
	appName := s.env.GetEnv("APP_NAME")
	enabled := s.env.GetEnvBool("FEATURE_X_ENABLED", "false")
	// ...
}
```
---

### 🏗️ Carregamento

No bootstrap da aplicação (`nanogo.Bootstrap()`), o sistema detecta o provider:

```go
err := env.Loader(i18n)
```

A partir disso, ele carrega as variáveis correspondentes.

---

### 📌 Locais Esperados do `.env` (para `ENV_FILE`)

O loader tenta automaticamente:

```txt
../../configs/.env
../configs/.env
configs/.env
.env
```

---

### ☁️ Integração com AWS Secrets Manager

Quando `ENV_PROVIDER=AWS`, é necessário configurar o seguinte:

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `AWS_SECRET_MANAGER_REGION` | ✅ | Região AWS do Secret Manager |
| `AWS_SECRET_MANAGER_ACCESS_KEY` | ✅* | Access key da IAM User |
| `AWS_SECRET_MANAGER_SECRET_KEY` | ✅* | Secret key da IAM User |
| `AWS_SECRET_MANAGER_NAME` | ✅ | Nome do segredo a ser carregado |

> ✅* **Ou** configurar `AWS_ROLE_ARN` e `AWS_WEB_IDENTITY_TOKEN_FILE` via IAM Role para pods no EKS.

---

### 🧪 Testando Localmente

Crie um arquivo `.env` com:

```env
ENV_PROVIDER=ENV_FILE
ENV=development
APP_NAME=nanogo
LOG_LEVEL=DEBUG
WEB_SERVER_PORT=8080
WEBSERVER_ACCESS_LOG=true
WEBSERVER_ORIGINS=http://localhost:3000
CACHE_PROVIDER=REDIS
REDIS_ADDR=localhost:6379
REDIS_NAMESPACE=nanogo
```

Inicie a aplicação. Ela deve carregar normalmente, incluindo middlewares e rota de healthcheck.

---

### ✅ Lista de Variáveis Suportadas

#### 🌐 Comuns a todos os ambientes
| Variável              | Descrição                                 | Exemplo                     |
|-----------------------|-------------------------------------------|-----------------------------|
| `ENV_PROVIDER`        | Define o provedor: `ENV_FILE`, `AWS`, `OS` | `ENV_FILE`                  |
| `APP_NAME`            | Nome da aplicação                         | `nanogo`                    |
| `ENV`                 | Ambiente da aplicação                     | `development`, `production` |
| `LOG_LEVEL`           | Nível de log: `DEBUG`, `INFO`, etc.       | `DEBUG`                     |
| `VERSION`             | Versão da aplicação (opcional)            | `1.0.0`                     |

#### 🚀 Web Server
| Variável                      | Descrição                                  |
|-------------------------------|--------------------------------------------|
| `WEB_SERVER_PORT`             | Porta de execução do servidor HTTP         |
| `WEB_SERVER_HOST`             | Host a ser escutado                        |
| `WEB_SERVER_CERTIFICATE`      | Caminho do certificado TLS (HTTPS)         |
| `WEB_SERVER_KEY`              | Caminho da chave TLS (HTTPS)               |
| `WEB_SERVER_MAX_UPLOAD_SIZE`  | Limite máximo de upload (em MB)            |
| `WEBSERVER_ACCESS_LOG`        | `true` para logar payloads                 |
| `WEBSERVER_ORIGINS`           | Origens permitidas no CORS                 |
| `WEBSERVER_HEADERS`           | Cabeçalhos permitidos no CORS              |
| `WEBSERVER_METHODS`           | Métodos permitidos no CORS                 |

#### 🧠 Redis
| Variável            | Descrição                   |
|---------------------|-----------------------------|
| `CACHE_PROVIDER`    | Provedor de cache (REDIS)   |
| `REDIS_ADDR`        | Endereço Redis              |
| `REDIS_PASSWORD`    | Senha Redis (opcional)      |
| `REDIS_NAMESPACE`   | Prefixo para chaves Redis   |

#### 🐰 RabbitMQ
| Variável             | Descrição                           |
|----------------------|-------------------------------------|
| `QUEUE_PROVIDER`     | Tipo de fila: `RABBITMQ`            |
| `RABBITMQ_PROTOCOL`  | `amqp`                              |
| `RABBITMQ_USER`      | Usuário RabbitMQ                    |
| `RABBITMQ_PASSWORD`  | Senha RabbitMQ                      |
| `RABBITMQ_HOST`      | Host RabbitMQ                       |
| `RABBITMQ_PORT`      | Porta RabbitMQ                      |
| `RABBITMQ_VHOST`     | Virtual host usado                  |

#### 📈 Telemetria (OpenTelemetry)
| Variável              | Descrição                            |
|-----------------------|----------------------------------------|
| `TELEMETRY_ENABLE`    | Ativa ou desativa a telemetria        |
| `TELEMETRY_DISPATCHER`| `OPEN_TELEMETRY` ou `MEMORY`          |
| `TELEMETRY_ENDPOINT`  | Endpoint do collector gRPC            |

#### ☁️ AWS (caso `ENV_PROVIDER=AWS`)
| Variável                             | Descrição                                   |
|--------------------------------------|---------------------------------------------|
| `AWS_SECRET_MANAGER_REGION`          | Região do secret manager                    |
| `AWS_SECRET_MANAGER_ACCESS_KEY`      | Access key da IAM User                      |
| `AWS_SECRET_MANAGER_SECRET_KEY`      | Secret key da IAM User                      |
| `AWS_SECRET_MANAGER_NAME`            | Nome do segredo a ser carregado             |

> Obs: Para autenticação por role (EKS), utilizar:
> - `AWS_ROLE_ARN`
> - `AWS_WEB_IDENTITY_TOKEN_FILE`
