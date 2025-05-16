# Tratamento de Erros

Define estruturas padronizadas para retorno de erros em APIs.

## Principais APIs
- Estrutura `CustomError` com campos `Code`, `Message` e `Details`.
- Funções auxiliares como `InvalidPayload`.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
err := errors.InvalidPayload("campo x é obrigatório")
return err
```
