# Documentação: Carregamento de Ambiente Remoto

O Nanogo suporta o carregamento de variáveis de ambiente a partir de fontes remotas, como AWS Secrets Manager e HashiCorp Vault, permitindo centralizar e proteger segredos de produção.

## Visão geral do carregamento remoto
O carregamento remoto é realizado por loaders específicos:
- **AWSLoader** (`pkg/env/env_aws_loader.go`): Carrega segredos do AWS Secrets Manager.
- **VaultLoader** (`pkg/env/provider_vault.go`): Carrega segredos do HashiCorp Vault.

A escolha do provedor é feita via variável `ENV_PROVIDER` (`AWS` ou `VAULT_HASHICORP`).

## Configuração de provedores remotos
### AWS Secrets Manager
Defina as variáveis:
```
ENV_PROVIDER=AWS
AWS_SECRET_MANAGER_REGION=us-east-1
AWS_SECRET_MANAGER_NAME=meu-segredo
```

### HashiCorp Vault
Defina as variáveis:
```
ENV_PROVIDER=VAULT_HASHICORP
VAULT_HASHICORP_HOST=https://vault.example.com
VAULT_HASHICORP_TOKEN=s.xxxxxxx
VAULT_HASHICORP_SECRET_PATH=secret/data/myapp
```

## Exemplos de uso
```go
import (
    "github.com/caiomarcatti12/nanogo/pkg/env"
    "github.com/caiomarcatti12/nanogo/pkg/i18n"
)

i18nService, _ := i18n.Factory("./translations")
_ = env.Loader(i18nService) // Carrega variáveis conforme ENV_PROVIDER
```

## Segurança e boas práticas
- Nunca exponha tokens ou segredos em código-fonte.
- Use roles e permissões restritas para acesso aos segredos.
- Prefira variáveis de ambiente para configuração de credenciais.

## Tratamento de erros
- Falhas de autenticação, segredos ausentes ou formato inválido resultam em erro fatal internacionalizado.
- O loader valida a presença de todas as variáveis obrigatórias antes de tentar carregar os segredos.

---

> Consulte também: [Gerenciamento de Variáveis de Ambiente (env)](./features/env.md)
