# Event Dispatcher

Sistema simples para publicação e consumo de eventos.

## Principais APIs
- `Factory(env, log, i18n) IEventDispatcher` seleciona o dispatcher.
- Implementação `InMemoryBroker` para testes e uso leve.

## Variáveis de Ambiente
- `EVENT_DISPATCHER` define o provedor (ex: `IN_MEMORY`).

## Exemplo de Uso
```go
dispatcher := event.Factory(envAdapter, logger, i18n)
dispatcher.RegisterConsumer(event.EventConsumer{Channel: "teste", Key: "demo", IHandler: MyHandler, HandlerFunc: "Handle"})
dispatcher.Dispatch(event.Event{Channel: "teste", Key: "demo", Data: payload})
```
