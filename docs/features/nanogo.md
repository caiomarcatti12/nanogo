# Bootstrap Nanogo

Pacote auxiliar para iniciar rapidamente a aplicação registrando todos os serviços básicos.

## Principais APIs
- `Bootstrap()` carrega ambiente, registra dependências e inicia serviços padrão.

## Variáveis de Ambiente
As mesmas exigidas pelos serviços registrados (env, log, db, etc.).

## Exemplo de Uso
```go
nanogo.Bootstrap()
```
