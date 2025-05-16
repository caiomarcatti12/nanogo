# RSQL Parser

Utilitário para interpretar consultas RSQL simples e gerar filtros.

## Principais APIs
- `Parse(query string) ([]Condition, error)` retorna condições estruturadas.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
conds, err := rsql.Parse("name==john;age=like=3")
```
