# Mapper

Ferramentas para serialização e desserialização de structs e mapas.

## Principais APIs
- `Serialize(data interface{}) interface{}` converte structs para mapas.
- `Deserialize(src, dest)` popula estruturas a partir de mapas/JSON.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
m := mapper.Serialize(myStruct)
var dst MyStruct
mapper.Deserialize(m, &dst)
```
