# Suporte ao NATS

O Nanogo possibilita utilizar o servidor de mensageria NATS de forma simples. Defina `QUEUE_PROVIDER=NATS` e configure `NATS_URL` para apontar para sua instância.

A publicação e o consumo seguem a interface `IQueue` já utilizada para RabbitMQ. Exemplo básico:

```go
queueCfg := queue.NatsQueue{
    Name:       "demo.subject",
    QueueGroup: "demo-group",
}

_ = queueManager.Configure(queueCfg)

consumer := NewDemoConsumer(logger)
go queueManager.Consume(&queueCfg, consumer)

msg := DemoMessage{ID: "1", Text: "hello"}
_ = queueManager.Publish(queueCfg.Name, "", msg)
```
