# Context Manager

Abstração para armazenar e recuperar valores em contexto seguro entre goroutines.

## Principais APIs
- `NewSafeContextManager() ISafeContextManager` cria a instância singleton.
- Métodos: `SetValues`, `CreateValue`, `GetValue`.

## Variáveis de Ambiente
Nenhuma.

## Exemplo de Uso
```go
cm := context_manager.NewSafeContextManager()
cm.SetValues(cm.CreateValue("key", "val"), func(){
    // código com contexto
})
```
