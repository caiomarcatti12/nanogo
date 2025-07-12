# Queue AddConsumer Feature

O método `AddConsumer` foi adicionado ao package `queue` para simplificar o processo de registro de consumers, seguindo o mesmo padrão do `AddRoute` do webserver.

## Objetivo

Evitar o processo repetitivo de:
1. Registrar o handler no container DI
2. Configurar a fila
3. Criar o consumer
4. Iniciar o consumo

## Estruturas

### IConsumerHandler
```go
type IConsumerHandler interface {
    Handler(body interface{}, headers map[string]interface{}) error
}
```

### QueueConsumer
```go
type QueueConsumer struct {
    Queue   Queue       // Configuração da fila (NatsQueue, RabbitmqQueue, etc.)
    Handler interface{} // Factory function do consumer handler
}
```

## Uso

### Método Antigo (múltiplas etapas)
```go
func main() {
    nanogo.Bootstrap()

    // 1. Registrar o handler no DI
    if err := di.GetInstance().Register(NewDemoConsumer); err != nil {
        panic(err)
    }

    queueManager, err := di.Get[queue.IQueue]()
    if err != nil {
        panic(err)
    }

    logger, _ := di.Get[log.ILog]()

    // 2. Configurar a fila
    queueCfg := queue.NatsQueue{
        Name:       "demo.subject",
        QueueGroup: "demo-group",
    }
    _ = queueManager.Configure(queueCfg)

    // 3. Criar o consumer
    consumer := NewDemoConsumer(logger)

    // 4. Iniciar o consumo
    go func() {
        if err := queueManager.Consume(&queueCfg, consumer); err != nil {
            logger.Error("error consuming", "err", err)
        }
    }()
}
```

### Método Novo (uma única chamada)
```go
func main() {
    nanogo.Bootstrap()

    queueManager, err := di.Get[queue.IQueue]()
    if err != nil {
        panic(err)
    }

    // Configuração da fila
    queueCfg := queue.NatsQueue{
        Name:       "demo.subject",
        QueueGroup: "demo-group",
    }

    // Tudo em uma única chamada!
    consumer := queue.QueueConsumer{
        Queue:   &queueCfg,
        Handler: NewDemoConsumer,
    }

    if err := queueManager.AddConsumer(consumer); err != nil {
        panic(err)
    }
}
```

## Vantagens

1. **Simplicidade**: Reduz múltiplas etapas para uma única chamada
2. **Menos erros**: Menor chance de esquecer alguma etapa
3. **Consistência**: Segue o mesmo padrão do webserver
4. **Clareza**: Fica claro qual fila está associada a qual consumer

## Providers Suportados

- ✅ NATS
- ✅ RabbitMQ
- ❌ Redis (provider comentado)

## Implementação

O método `AddConsumer` executa internamente:
1. `Configure(consumer.Queue)` - configura a fila
2. `di.GetInstance().Register(consumer.Handler)` - registra no DI
3. `Consume(consumer.Queue, consumer.Handler)` - inicia o consumo

## Exemplo Completo

Veja o exemplo em `/examples/nats_addconsumer/main.go` para um caso de uso completo.
