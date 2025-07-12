/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package queue

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/bytedance/sonic"

	"github.com/caiomarcatti12/nanogo/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/mapper"
	"github.com/caiomarcatti12/nanogo/pkg/metric"
	"github.com/caiomarcatti12/nanogo/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var (
	once     sync.Once
	instance Rabbitmq
)

type Rabbitmq struct {
	DataConnection DataConnection
	Connection     *amqp.Connection
	Channel        *amqp.Channel
	logger         log.ILog
	metricMonitor  metric.IMetric
	telemetry      telemetry.ITelemetry
	exchanges      map[string]RabbitmqExchange
	queues         map[string]RabbitmqQueue
}

type DataConnection struct {
	Protocol string
	User     string
	Password string
	Host     string
	Port     string
	Vhost    string
}

type ExchangeType string

const (
	Direct  ExchangeType = "direct"
	Fanout  ExchangeType = "fanout"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

type RabbitmqExchange struct {
	Name       string
	Durable    bool
	Type       ExchangeType
	AutoDel    bool
	Internal   bool
	NoWait     bool
	Parameters amqp.Table
}

type RabbitmqQueue struct {
	Name        string
	RoutingKey  string
	ConsumerTag string
	Durable     bool
	AutoDel     bool
	Exclusive   bool
	NoWait      bool
	NoLocal     bool
	Parameters  amqp.Table
}

// BindingConfig descreve *qualquer* tipo de bind:
// Exchange -> Exchange, ou Exchange -> Queue, etc.
type BindingConfig struct {
	Source      string // nome do exchange de origem
	Destination string // nome do exchange ou queue de destino
	RoutingKey  string
	NoWait      bool
	Args        map[string]interface{}
}

func (r *RabbitmqQueue) GetName() string {
	return r.Name
}

func NewInstanceRabbitmq(env env.IEnv, logger log.ILog, metricMonitor metric.IMetric, telemetry telemetry.ITelemetry) (IQueue, error) {
	logger.Info("Creating instance of RabbitMQ...")

	once.Do(func() {
		instance = Rabbitmq{
			logger:        logger,
			metricMonitor: metricMonitor,
			telemetry:     telemetry,
			DataConnection: DataConnection{
				Protocol: env.GetEnv("RABBITMQ_PROTOCOL"),
				User:     env.GetEnv("RABBITMQ_USER"),
				Password: env.GetEnv("RABBITMQ_PASSWORD"),
				Host:     env.GetEnv("RABBITMQ_HOST"),
				Port:     env.GetEnv("RABBITMQ_PORT"),
				Vhost:    env.GetEnv("RABBITMQ_VHOST"),
			},
			exchanges: make(map[string]RabbitmqExchange),
			queues:    make(map[string]RabbitmqQueue),
		}

		instance.Connect()
	})

	return &instance, nil
}

func (r *Rabbitmq) Connect() error {

	r.logger.Info("Connecting to RabbitMQ...")

	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		r.DataConnection.Protocol,
		r.DataConnection.User,
		r.DataConnection.Password,
		r.DataConnection.Host,
		r.DataConnection.Port,
		r.DataConnection.Vhost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	r.logger.Info("Creating channel...")
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	r.logger.Info("Connected to RabbitMQ!")
	r.metricMonitor.SetGauge(QueueManagerConnected.String(), 1, map[string]string{"provider": "rabbitmq"})

	r.Connection = conn
	r.Channel = ch

	return nil
}

func (r *Rabbitmq) SetPrefetch(prefetchCount int) error {
	r.logger.Infof("Setting prefetch count to %d", prefetchCount)
	err := r.Channel.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	return nil
}

func (r *Rabbitmq) Configure(args ...interface{}) error {
	r.logger.Info("Declaring exchanges, queues, and bindings...")

	for _, arg := range args {
		switch v := arg.(type) {
		case RabbitmqExchange:
			if err := r.declareExchange(v); err != nil {
				return err
			}
			r.exchanges[v.Name] = v

		case RabbitmqQueue:
			if err := r.declareQueue(v); err != nil {
				return err
			}
			r.queues[v.Name] = v

		case BindingConfig:
			exchangeSourceConfig, ok := r.exchanges[v.Source]
			if !ok {
				err := fmt.Errorf("Source exchange '%s' not found in configuration", v.Source)
				r.logger.Error(err.Error())
				return err
			}

			if queueConfig, ok := r.queues[v.Destination]; ok {
				r.logger.Infof("Binding queue '%s' to exchange '%s' with routing key '%s'", queueConfig.Name, exchangeSourceConfig.Name, v.RoutingKey)
				if err := r.Channel.QueueBind(
					queueConfig.Name,
					v.RoutingKey,
					exchangeSourceConfig.Name,
					v.NoWait,
					amqp.Table(v.Args),
				); err != nil {
					r.logger.Errorf("Failed to bind queue '%s' to exchange '%s': %v", queueConfig.Name, exchangeSourceConfig.Name, err)
					return err
				}
				r.metricMonitor.SetGauge(QueueBindExchange.String(), 1, map[string]string{"exchange": exchangeSourceConfig.Name, "queue": queueConfig.Name})
				r.logger.Infof("Successfully bound queue '%s' to exchange '%s'", queueConfig.Name, exchangeSourceConfig.Name)
			} else if exchangeDestConfig, ok := r.exchanges[v.Destination]; ok {
				r.logger.Infof("Binding exchange '%s' to exchange '%s' with routing key '%s'", exchangeDestConfig.Name, exchangeSourceConfig.Name, v.RoutingKey)
				if err := r.Channel.ExchangeBind(
					exchangeDestConfig.Name,
					v.RoutingKey,
					exchangeSourceConfig.Name,
					v.NoWait,
					amqp.Table(v.Args),
				); err != nil {
					r.logger.Errorf("Failed to bind exchange '%s' to exchange '%s': %v", exchangeDestConfig.Name, exchangeSourceConfig.Name, err)
					return err
				}
				r.metricMonitor.SetGauge(QueueBindExchange.String(), 1, map[string]string{"exchange": exchangeSourceConfig.Name, "queue": exchangeDestConfig.Name})
				r.logger.Infof("Successfully bound exchange '%s' to exchange '%s'", exchangeDestConfig.Name, exchangeSourceConfig.Name)
			} else {
				err := fmt.Errorf("Destination '%s' not found in queue or exchange configuration", v.Destination)
				r.logger.Error(err.Error())
				return err
			}
		}
	}

	return nil
}

func (r *Rabbitmq) Publish(exchange string, routingKey string, body interface{}) (err error) {
	span := r.telemetry.StartChildSpan("Publish message to RabbitMQ", map[string]interface{}{"exchange": exchange, "routing_key": routingKey, "args": body})
	defer func() { r.telemetry.EndSpan(span, err) }()

	r.logger.Info("Publishing message to RabbitMQ...")

	r.metricMonitor.IncrementCounter(QueueMessagePublish.String(), map[string]string{"queue": exchange, "routing_key": routingKey})

	bodyBytes, err := sonic.Marshal(body)
	if err != nil {
		r.logger.Error(err.Error())
		return
	}

	fcm := context_manager.NewSafeContextManager()
	correlationID, ok := fcm.GetValue("x-correlation-id")

	if !ok {
		correlationID = uuid.New().String()
	}

	// Add correlation ID to headers
	headers := amqp.Table{
		"x-correlation-id": correlationID,
	}

	err = r.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			Headers:     headers,
		},
	)

	if err != nil {
		r.logger.Error(err.Error())
		return
	}

	return nil
}

// Estrutura para exemplo
type Example struct {
	Id int `json:"id"`
}

func (r *Rabbitmq) Consume(queue Queue, consumerHandler interface{}) error {
	// Validação dos parâmetros de entrada
	if queue == nil {
		return fmt.Errorf("queue cannot be nil")
	}
	if consumerHandler == nil {
		return fmt.Errorf("consumerHandler cannot be nil")
	}

	r.logger.Infof("Consuming messages from %s", queue.GetName())

	// Verificar se a fila existe na configuração
	rabbitmQueueConfig, ok := r.queues[queue.GetName()]
	if !ok {
		return fmt.Errorf("queue %s not found in configuration", queue.GetName())
	}

	msgs, err := r.Channel.Consume(
		rabbitmQueueConfig.Name,        // queue
		rabbitmQueueConfig.ConsumerTag, // consumer
		false,                          // auto-ack
		rabbitmQueueConfig.Exclusive,   // exclusive
		rabbitmQueueConfig.NoLocal,     // no-local
		rabbitmQueueConfig.NoWait,      // no-wait
		nil,                            // args
	)

	if err != nil {
		return err
	}

	// Criar a instância do consumidor uma vez
	consumerInstance, err := r.createConsumerInstance(consumerHandler)
	if err != nil {
		return err
	}

	r.metricMonitor.SetGauge(QueueConsummerConnected.String(), 1, map[string]string{"queue": rabbitmQueueConfig.Name})

	for d := range msgs {
		go r.processMessages(d, queue, consumerInstance)
	}

	return nil
}

func (r *Rabbitmq) createConsumerInstance(consumerHandler interface{}) (interface{}, error) {
	di.GetInstance().Register(consumerHandler)

	consumer, err := di.GetInstance().GetByFactory(consumerHandler)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// Função privada para processar mensagens
func (r *Rabbitmq) processMessages(d amqp.Delivery, queue Queue, consumer interface{}) {
	r.logger.Trace("Processing message from %s", queue.GetName())

	fcm := context_manager.NewSafeContextManager()

	correlationID, ok := d.Headers["x-correlation-id"].(string)
	if !ok || correlationID == "" {
		correlationID = uuid.New().String()
	}

	contextValues := fcm.CreateValue("x-correlation-id", correlationID)
	fcm.SetValues(contextValues, func() {
		rootSpan := r.telemetry.CreateRootSpan(fmt.Sprintf("Process message queue %s", queue.GetName()), map[string]interface{}{"correlationID": correlationID})
		defer func() { r.telemetry.EndSpan(rootSpan, nil) }()

		headers := make(map[string]interface{})
		for k, v := range d.Headers {
			headers[k] = v
		}

		err := r.callConsumerHandler(consumer, d.Body, headers)
		if err != nil {
			r.logger.Error(err.Error())
			d.Nack(false, false)
			r.metricMonitor.IncrementCounter(QueueMessageNack.String(), map[string]string{"queue": queue.GetName()})
		} else {
			d.Ack(false)
			r.metricMonitor.IncrementCounter(QueueMessageAck.String(), map[string]string{"queue": queue.GetName()})
		}
	})
}

func (r *Rabbitmq) Disconnect() error {
	r.logger.Info("Disconnecting from RabbitMQ...")
	return r.Connection.Close()
}

func (r *Rabbitmq) declareExchange(exchange RabbitmqExchange) error {
	r.logger.Infof("Declaring exchange %s on RabbitMQ...", exchange.Name)
	err := r.Channel.ExchangeDeclare(
		exchange.Name,
		string(exchange.Type),
		exchange.Durable,
		exchange.AutoDel,
		exchange.Internal,
		exchange.NoWait,
		nil,
	)

	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	r.exchanges[exchange.Name] = exchange

	r.metricMonitor.SetGauge(QueueExchangeCreated.String(), 1, map[string]string{"exchange": exchange.Name})

	return nil
}

func (r *Rabbitmq) declareQueue(queue RabbitmqQueue) error {
	r.logger.Infof("Declaring queue %s on RabbitMQ...", queue.Name)

	_, err := r.Channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDel,
		queue.Exclusive,
		queue.NoWait,
		queue.Parameters,
	)

	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	r.queues[queue.Name] = queue

	r.metricMonitor.SetGauge(QueueCreated.String(), 1, map[string]string{"queue": queue.Name})

	return nil
}

// bindExchange realiza o bind entre duas exchanges
func (r *Rabbitmq) bindExchange(exchangeSource RabbitmqExchange, exchangeDest RabbitmqExchange) error {
	r.logger.Infof("Binding exchange %s to exchange %s with routing key '%s'", exchangeSource.Name, exchangeSource.Name, "*")

	err := r.Channel.ExchangeBind(
		exchangeDest.Name,   // exchange destino
		"*",                 // routing key
		exchangeSource.Name, // exchange origem
		true,                // noWait
		nil,                 // argumentos adicionais
	)

	if err != nil {
		r.logger.Errorf("Failed to bind exchange %s to exchange %s: %v", exchangeSource.Name, exchangeSource.Name, err)
		return err
	}

	r.metricMonitor.SetGauge(QueueBindExchange.String(), 1, map[string]string{"exchange": exchangeSource.Name, "queue": exchangeSource.Name})

	r.logger.Infof("Successfully bound exchange %s to exchange %s", exchangeSource.Name, exchangeSource.Name)

	return nil
}

func (r *Rabbitmq) bindQueue(exchange RabbitmqExchange, queue RabbitmqQueue) error {
	r.logger.Infof("Binding queue %s to exchange %s on RabbitMQ...", queue.Name, exchange.Name)

	err := r.Channel.QueueBind(
		queue.Name,
		queue.RoutingKey,
		exchange.Name,
		true,
		nil,
	)

	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	r.metricMonitor.SetGauge(QueueBindExchange.String(), 1, map[string]string{"exchange": exchange.Name, "queue": queue.Name})

	return nil
}

// Função para verificar se consumerHandler implementa IConsumer ignorando o tipo
func (r *Rabbitmq) callConsumerHandler(consumer interface{}, body []byte, headers map[string]interface{}) (err error) {
	span := r.telemetry.StartChildSpan("consumerHandler")
	defer r.telemetry.EndSpan(span, err)

	handlerType := reflect.TypeOf(consumer)
	handlerValue := reflect.ValueOf(consumer)

	// Verificar se o método Handler existe
	method, exists := handlerType.MethodByName("Handler")
	if !exists {
		err = errors.New("Handler method not found")
		return
	}

	// Obter o primeiro parâmetro do método Handler
	if method.Type.NumIn() != 3 {
		err = errors.New("Handler method has an incorrect number of parameters")
		return
	}

	bodyType := method.Type.In(1)

	// Criar uma nova instância do tipo do primeiro parâmetro
	bodyValue := reflect.New(bodyType).Interface()

	// Declarar uma variável do tipo []map[string]interface{}
	var result map[string]interface{}

	// Fazer o parse do JSON
	err = sonic.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal:", err)
		return
	}

	// Desserializar msg para a instância criada
	if err = mapper.Deserialize(result, &bodyValue); err != nil {
		err = fmt.Errorf("failed to unmarshal message: %v", err)
		return
	}

	// Preparar os parâmetros para a chamada do método Handler
	params := []reflect.Value{
		handlerValue,
		reflect.ValueOf(bodyValue).Elem(),
		reflect.ValueOf(headers),
	}

	// Chamar o método Handler usando reflexão
	results := method.Func.Call(params)

	// Verificar se ocorreu algum erro
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}
