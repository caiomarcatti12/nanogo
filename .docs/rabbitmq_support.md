### **Trabalhando com RabbitMQ no Nanogo**

Nesta seção, vamos descrever como você pode configurar e trabalhar com RabbitMQ usando o Nanogo.

#### **Configuração de Variáveis de Ambiente**

A biblioteca se conecta ao RabbitMQ usando as informações fornecidas através das variáveis de ambiente. Certifique-se de configurar corretamente as seguintes variáveis:

```env
RABBITMQ_HOST=host.rabbitmq
RABBITMQ_USER=user
RABBITMQ_PASSWORD=password
RABBITMQ_PORT=5672
RABBITMQ_VHOST=/
```

### **Simplicidade na Implementação com Nanogo**

A biblioteca Nanogo oferece uma camada adicional de abstração que facilita a integração com o RabbitMQ, encapsulando as operações complexas de maneira simples e direta. Isso permite que os desenvolvedores criem rapidamente aplicativos que podem publicar e consumir mensagens do RabbitMQ sem ter que lidar com os detalhes internos do RabbitMQ.

No exemplo abaixo, demonstramos como é simples iniciar um consumidor utilizando a biblioteca Nanogo. Aqui, encapsulamos a configuração da `Exchange` e da `Queue` em suas respectivas funções e criamos um consumidor que implementa a interface `Consumer`:

```go
package main

import (
	"log"

	"github.com/caiomarcatti12/nanogo/v2/config"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/rabbitmq"
)

type MyConsumer struct{}

func (mc *MyConsumer) Consume(body map[string]interface{}, headers map[string]interface{}) {
	log.Printf("Headers: %v", headers)
	log.Printf("Body: %s", body)
}

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	// Cria um consumidor da fila MyConsumer
	rabbitmq.Consume(exchange(), queue(), &MyConsumer{})

	config.WaitSignalStop()
}

func exchange() rabbitmq.Exchange {
	return rabbitmq.Exchange{
		Name:    "teste-exchange-go",
		Durable: true,
		Type:    "direct",
		AutoDel: false,
		NoWait:  false,
	}
}

func queue() rabbitmq.Queue {
	return rabbitmq.Queue{
		Name:       "teste-queue",
		Durable:    true,
		AutoDel:    false,
		Exclusive:  false,
		NoWait:     false,
		Parameters: nil,
	}
}
```

#### **Passos para Implementar**

1. **Definição do Consumidor**: Criamos uma estrutura `MyConsumer` que implementa a interface `Consumer`. O método `Consume` dessa estrutura será chamado cada vez que uma nova mensagem for recebida, fornecendo detalhes do corpo e dos cabeçalhos da mensagem.

2. **Carregando Variáveis de Ambiente**: Utilizamos `env.LoadEnv()` para carregar as variáveis de ambiente do arquivo `.env`.

3. **Definição da Exchange e da Queue**: Criamos funções `exchange` e `queue` para encapsular as configurações da exchange e da queue, respectivamente. Isso facilita a reutilização e a manutenção do código.

4. **Inicialização do Consumidor**: Na função `main`, inicializamos o consumidor chamando `rabbitmq.Consume` com a exchange e a queue que definimos, junto com uma instância de `MyConsumer`.

5. **Espera por Sinais de Parada**: Por último, usamos `config.WaitSignalStop()` para manter o aplicativo em execução, esperando por sinais de parada para encerrar graceiosamente.

Com isso, conseguimos ter uma aplicação concisa e bem estruturada, facilitando tanto o desenvolvimento quanto a manutenção do código. O Nanogo se torna um grande aliado para trabalhar com RabbitMQ em Go, simplificando significativamente o processo de setup e manipulação das mensagens do RabbitMQ.

### **Compreendendo as Estruturas e Constantes**

##### **ExchangeType**

A biblioteca define quatro tipos de exchanges através do tipo `ExchangeType`. São eles:

- `Direct`
- `Fanout`
- `Topic`
- `Headers`

##### **Exchange**

A estrutura `Exchange` permite definir as propriedades de uma exchange. Os campos incluem:

- `Name`: Nome da exchange.
- `Durable`: Determina se a exchange sobreviverá a reboots do broker.
- `Type`: O tipo da exchange, definido usando `ExchangeType`.
- `AutoDel`: Define se a exchange será deletada automaticamente quando não estiver mais sendo usada.
- `Internal`: Define se a exchange é interna e não pode ser publicada diretamente por publishers.
- `NoWait`: Define se a declaração da exchange será não-bloqueante.

##### **Queue**

A estrutura `Queue` permite definir as propriedades de uma fila (queue). Os campos são:

- `Name`: O nome da fila.
- `Key`: A chave de roteamento para a fila.
- `Durable`: Determina se a fila sobreviverá a reboots do broker.
- `AutoDel`: Define se a fila será deletada automaticamente quando não estiver mais sendo usada.
- `Exclusive`: Define se a fila é exclusiva para a conexão atual.
- `NoWait`: Define se a declaração da fila será não-bloqueante.
- `Parameters`: Permite definir parâmetros adicionais para a fila, usando uma tabela AMQP.
  Claro, peço desculpas pela confusão. Vamos manter as descrições junto com os exemplos de código para cada tópico.

### **Inicializando uma Conexão**

Para estabelecer uma conexão com o servidor RabbitMQ, é utilizado o método `NewInstanceRabbitmq`. Este método segue o padrão singleton, garantindo que apenas uma instância da conexão seja criada durante a execução do aplicativo. Ele se encarrega de buscar os detalhes de conexão do arquivo de configuração (como host, usuário, senha, etc.) e estabelece uma conexão com o servidor RabbitMQ.

Aqui está um exemplo de como você pode obter uma instância da conexão:

```go
connection := NewInstanceRabbitmq()
```

### **Declarando Exchanges e Queues**

Após inicializar uma conexão, o próximo passo é declarar as exchanges e as queues que serão utilizadas para o envio e recebimento de mensagens. Veja abaixo exemplos de como declarar uma exchange e uma queue.

#### **Declarando uma Exchange**

As exchanges são declaradas através do método `DeclareExchange`, onde você pode especificar propriedades como nome, tipo e durabilidade da exchange:

```go
exchange := rabbitmq.Exchange{
    Name:    "my_exchange",
    Type:    rabbitmq.Direct,
    Durable: true,
}

err := rabbitmq.DeclareExchange(exchange)
if err != nil {
    log.Fatalf("Failed to declare exchange: %s", err)
}
```

#### **Declarando uma Queue**

Da mesma forma, as queues são declaradas usando o método `DeclareQueue`, definindo propriedades como nome e durabilidade da queue:

```go
queue := rabbitmq.Queue{
    Name:    "my_queue",
    Key:     "my_routing_key",
    Durable: true,
}

q, err := rabbitmq.DeclareQueue(queue)
if err != nil {
    log.Fatalf("Failed to declare queue: %s", err)
}
```

### **Publicando Mensagens**

Uma vez que as exchanges e queues estejam devidamente configuradas, você poderá começar a publicar mensagens utilizando o método `Publish`. Este método requer que você especifique o nome da exchange, a routing key e o corpo da mensagem a ser enviada.

Aqui está um exemplo de como publicar uma mensagem:

```go
body := map[string]interface{}{"hello": "world"}
rabbitmq.Publish(exchange.Name, queue.Key, body)
```

### **Consumindo Mensagens**

Para consumir mensagens, você precisa criar um consumidor que implemente a interface `Consumer`. Este consumidor deve ter um método `Consume`, que será chamado cada vez que uma nova mensagem for recebida.

Aqui está um exemplo de como você pode criar um consumidor e iniciar a consumir mensagens:

```go
type MyConsumer struct{}

func (c *MyConsumer) Consume(body map[string]interface{}, headers map[string]interface{}) {
    log.Printf("Received message with body: %v and headers: %v", body, headers)
}

// ...

rabbitmq.Consume(exchange, queue, &MyConsumer{})
```

A função `Consume` será chamada para cada mensagem recebida, onde `body` e `headers` representam o corpo e os cabeçalhos da mensagem, respectivamente.

