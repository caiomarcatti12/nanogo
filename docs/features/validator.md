# Validador

Funções de validação de estruturas usando a biblioteca `go-playground/validator`.

## Principais APIs
- `ValidateStruct(v interface{}) error` realiza validação de struct via tags.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
err := validator.ValidateStruct(req)
```
