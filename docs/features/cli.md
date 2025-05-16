# CLI

Fornece utilidades para criação de aplicações de linha de comando.

## Principais APIs
- `Factory(log) ICLI` retorna a implementação de CLI.
- Métodos: `Add`, `Initialize`, `Execute`, `Get`, `GetAll`.

## Variáveis de Ambiente
Nenhuma específica.

## Exemplo de Uso
```go
cli := cli.Factory(logger)
cli.Add("--name", "", "Seu nome")
cli.Initialize()
value, _ := cli.Get("--name")
```
