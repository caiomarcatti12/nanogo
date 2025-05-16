# Injeção de Dependências

Container simples para registrar e resolver dependências.

## Principais APIs
- `Factory(i18n, log) IContainer` cria o singleton.
- `Register`, `RegisterAll`, `GetByFactory`, `GetByName`.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
container := di.Factory(i18n, logger)
container.Register(MyFactory)
svc, _ := container.GetByFactory(MyFactory)
```
