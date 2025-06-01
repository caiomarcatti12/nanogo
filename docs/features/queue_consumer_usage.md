# Como Consumir uma Fila (Queue Consumer) no Nanogo

O consumo de filas no Nanogo segue o mesmo padrão de injeção de dependências e separação de responsabilidades dos services, mas envolve alguns pontos específicos relacionados à configuração, assinatura de eventos e tratamento de mensagens.

## Estrutura Básica de um Consumer

O consumer é responsável por processar mensagens recebidas de uma fila (ex: RabbitMQ). O ideal é definir uma interface para o consumer, facilitando testes e desacoplamento.

```go
package consumer

type IOrderConsumer interface {
    HandleOrderCreated(event OrderCreatedEvent) error
}

type OrderConsumer struct {
    // Dependências injetadas (ex: services, repos, logger)
    orderService service.IOrderService
}

func NewOrderConsumer(orderService service.IOrderService) IOrderConsumer {
    return &OrderConsumer{orderService: orderService}
}

func (c *OrderConsumer) HandleOrderCreated(event OrderCreatedEvent) error {
    // Processa o evento recebido da fila
    return c.orderService.ProcessOrder(event.OrderID)
}
```

## Registrando o Consumer no DI

```go
import "github.com/caiomarcatti12/nanogo/pkg/di"

di.RegisterFactory(NewOrderConsumer)
```

## Assinando a Fila e Consumindo Mensagens

O registro do consumer na fila é feito geralmente no bootstrap do projeto ou em um módulo de inicialização:

```go
import (
    "github.com/caiomarcatti12/nanogo/pkg/event"
    "github.com/caiomarcatti12/nanogo/pkg/di"
)

func init() {
    dispatcher, _ := di.Get[event.IEventDispatcher]()
    consumer, _ := di.Get[IOrderConsumer]()
    dispatcher.Subscribe("order.created", consumer.HandleOrderCreated)
}
```

## Exemplo de Definição de Evento

```go
type OrderCreatedEvent struct {
    OrderID string `json:"order_id"`
}
```

## Boas Práticas
- Sempre use interfaces para consumers.
- Injete dependências via construtor.
- Mantenha o handler de evento pequeno e delegue lógica para services.
- Trate erros e implemente retentativas se necessário.
- Use logs para rastrear o processamento das mensagens.

---

> Consulte também: [Documentação do Event Dispatcher](./features/event.md) e [Service Usage](./service_usage.md)
