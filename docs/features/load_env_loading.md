# Documentação: Carregamento de Ambiente Local

O Nanogo permite o carregamento de variáveis de ambiente a partir de arquivos locais `.env`, facilitando a configuração de diferentes ambientes (desenvolvimento, produção, etc) de forma segura e flexível.

## Visão geral do carregamento local
O carregamento local é realizado principalmente pelo loader `FileEnvLoader`, localizado em `pkg/env/env_file_loader.go`. Ele utiliza a biblioteca `godotenv` para ler arquivos `.env` e popular as variáveis de ambiente do sistema.

O método central é:

```go
func (f *FileEnvLoader) Load() error
```

Este método busca o arquivo `.env` em múltiplos caminhos possíveis, incluindo:
- Diretório do executável
- `configs/.env`
- Caminho definido pela variável `ENV_FILE_PATH`

Se o arquivo for encontrado, as variáveis são carregadas e ficam disponíveis via `os.Getenv` ou pelo serviço `Env` do framework.

## Como configurar arquivos .env
Exemplo de arquivo `.env` (veja `configs/.env.example`):

```
APP_NAME=Application
VERSION=1.0.0
ENV=dev
ENV_PROVIDER=local
WEB_SERVER_HOST=0.0.0.0
WEB_SERVER_PORT=8080
```

## Como utilizar loaders locais
O carregamento é feito automaticamente ao inicializar o framework, mas pode ser chamado manualmente:

```go
import (
    "github.com/caiomarcatti12/nanogo/pkg/env"
    "github.com/caiomarcatti12/nanogo/pkg/i18n"
)

i18nService, _ := i18n.Factory("./translations")
_ = env.Loader(i18nService) // Carrega variáveis conforme ENV_PROVIDER
```

O loader padrão é o arquivo `.env`, mas pode ser alterado via `ENV_PROVIDER`.

## Exemplos de código
Acesso seguro às variáveis usando o serviço `Env`:

```go
import "github.com/caiomarcatti12/nanogo/pkg/env"

value := envService.GetEnv("APP_NAME")
port := envService.GetEnv("WEB_SERVER_PORT", "8080")
```

## Possíveis erros e soluções
- **Arquivo .env não encontrado:** O loader procura em múltiplos caminhos, mas se não encontrar, dispara erro fatal internacionalizado.
- **Variável obrigatória ausente:** O método `GetEnv` dispara panic se a variável não existir e não houver valor padrão.
- **Permissões de arquivo:** Certifique-se de que o arquivo `.env` está acessível pelo processo.

## Boas práticas
- Nunca versionar arquivos `.env` com segredos reais.
- Utilize `.env.example` para documentar as variáveis necessárias.
- Defina `ENV_PROVIDER=local` para garantir o uso do loader local.

---

> Consulte também: [Gerenciamento de Variáveis de Ambiente (env)](./features/env.md)
