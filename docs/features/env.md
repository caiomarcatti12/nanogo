# Gerenciamento de Variáveis de Ambiente (env)

Este pacote fornece uma solução estruturada e flexível para o gerenciamento de variáveis de ambiente em aplicações Go, permitindo carregar configurações de múltiplas fontes.

## Estrutura

O pacote `env` inclui:

* **Loader**: Interface principal para carregamento de variáveis de ambiente.
* **Env**: Classe principal para recuperação segura e tipada das variáveis.

## Funcionalidades Principais da Classe `Env`

A classe `Env` é a principal interface utilizada para acessar variáveis de ambiente no sistema de forma segura, tipada e internacionalizada. Ela realiza validações rigorosas garantindo que variáveis obrigatórias estejam corretamente definidas.

### Instanciando o Serviço `Env`

Exemplo de criação de uma instância:

```go
import (
"github.com/caiomarcatti12/nanogo/pkg/env"
"github.com/caiomarcatti12/nanogo/pkg/i18n"
)

func main() {
i18nService, err := i18n.Factory("./translations")
if err != nil {
panic(err)
}

envService := env.NewEnv(i18nService)
}
```

### Métodos

#### `GetEnv`

* Obtém uma variável de ambiente como `string`.
* **Validação:**

  * Nome da variável não pode ser vazio ou apenas espaços em branco.
  * Caso a variável não exista e nenhum valor padrão seja fornecido, dispara um `panic` com uma mensagem internacionalizada.

**Exemplo:**

```go
// Sem valor padrão
value := envService.GetEnv("DATABASE_URL")

// Com valor padrão
value := envService.GetEnv("DATABASE_URL", "postgres://localhost:5432")
```

#### `GetEnvBool`

* Obtém uma variável de ambiente como `bool`.
* Utiliza internamente o método `GetEnv`.
* Realiza conversão segura para booleano, retornando `false` se a conversão falhar.

**Exemplo:**

```go
isProduction := envService.GetEnvBool("PRODUCTION_MODE", "false")
```

## Tratamento de Erros

* A classe dispara `panic` em situações críticas como variável não encontrada ou nome de variável inválido.
* Mensagens de erro são obtidas através do serviço `i18n` permitindo internacionalização dos erros.

## Provedores de Variáveis

O pacote permite carregar variáveis de diferentes fontes, especificadas através da variável `ENV_PROVIDER`:

* **AWS Secrets Manager** (`ENV_PROVIDER=AWS`): Utiliza AWS Secrets Manager para armazenar e carregar variáveis de ambiente.
* **Arquivo `.env`** (`ENV_PROVIDER=ENV_FILE`): Carrega variáveis a partir de arquivos locais `.env`.
* **HashiCorp Vault** (`ENV_PROVIDER=VAULT_HASHICORP`): Carrega variáveis utilizando segredos armazenados no HashiCorp Vault.

### Configuração

Defina o provedor via variável de ambiente:

```shell
export ENV_PROVIDER=AWS # ou ENV_FILE ou VAULT_HASHICORP
```

Cada provedor é gerenciado por loaders específicos que podem ser estendidos ou customizados conforme necessidade.

### Configuração específica para Vault HashiCorp

Ao utilizar o provedor Vault HashiCorp, as seguintes variáveis de ambiente devem ser definidas:

```shell
export VAULT_HASHICORP_HOST="https://vault.example.com"
export VAULT_HASHICORP_TOKEN="s.xxxxxxx"
export VAULT_HASHICORP_SECRET_PATH="secret/data/myapp"
```

**Observações:**

* O carregador do Vault utiliza o método de autenticação via token.
* Espera que os segredos sejam armazenados em formato chave-valor usando o mecanismo KV versão 2.
* Caso alguma variável ou o segredo estejam ausentes ou em formato inválido, ocorrerá um erro fatal devidamente internacionalizado.
