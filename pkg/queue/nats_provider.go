package queue

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/caiomarcatti12/nanogo/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/mapper"
	"github.com/caiomarcatti12/nanogo/pkg/metric"
	"github.com/caiomarcatti12/nanogo/pkg/telemetry"
	"github.com/google/uuid"
	nats "github.com/nats-io/nats.go"
)

// Nats represents a NATS connection manager.
type Nats struct {
	url           string
	Conn          *nats.Conn
	logger        log.ILog
	metricMonitor metric.IMetric
	telemetry     telemetry.ITelemetry
	queues        map[string]NatsQueue
}

// NatsQueue holds configuration of a subscription subject and queue group.
type NatsQueue struct {
	Name       string
	QueueGroup string
}

func (n *NatsQueue) GetName() string { return n.Name }

// NewInstanceNats creates a new NATS provider instance.
func NewInstanceNats(env env.IEnv, logger log.ILog, metricMonitor metric.IMetric, telemetry telemetry.ITelemetry) (IQueue, error) {
	instance := &Nats{
		url:           env.GetEnv("NATS_URL", nats.DefaultURL),
		logger:        logger,
		metricMonitor: metricMonitor,
		telemetry:     telemetry,
		queues:        make(map[string]NatsQueue),
	}

	if err := instance.Connect(); err != nil {
		return nil, err
	}
	return instance, nil
}

// Connect establishes the connection to NATS.
func (n *Nats) Connect() error {
	n.logger.Info("Connecting to NATS...")
	conn, err := nats.Connect(n.url)
	if err != nil {
		return err
	}
	n.Conn = conn
	n.metricMonitor.SetGauge(QueueManagerConnected.String(), 1, map[string]string{"provider": "nats"})
	n.logger.Info("Connected to NATS!")
	return nil
}

// Configure stores queue configurations. NATS does not require server side setup.
func (n *Nats) Configure(args ...interface{}) error {
	for _, arg := range args {
		switch v := arg.(type) {
		case *NatsQueue:
			n.queues[v.Name] = *v
			n.metricMonitor.SetGauge(QueueCreated.String(), 1, map[string]string{"queue": v.Name})
		}
	}
	return nil
}

// Publish sends a message to a NATS subject.
func (n *Nats) Publish(subject string, routingKey string, body interface{}) (err error) {
	span := n.telemetry.StartChildSpan("Publish message to NATS", map[string]interface{}{"subject": subject, "routing_key": routingKey, "args": body})
	defer func() { n.telemetry.EndSpan(span, err) }()

	finalSubject := subject
	if routingKey != "" {
		finalSubject = fmt.Sprintf("%s.%s", subject, routingKey)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	fcm := context_manager.NewSafeContextManager()
	correlationID, ok := fcm.GetValue("x-correlation-id")
	if !ok {
		correlationID = uuid.New().String()
	}

	msg := &nats.Msg{
		Subject: finalSubject,
		Data:    bodyBytes,
		Header:  nats.Header{"x-correlation-id": []string{fmt.Sprint(correlationID)}},
	}

	n.metricMonitor.IncrementCounter(QueueMessagePublish.String(), map[string]string{"queue": finalSubject, "routing_key": routingKey})
	return n.Conn.PublishMsg(msg)
}

// Consume subscribes to a subject and processes incoming messages.
func (n *Nats) Consume(queue Queue, consumerHandler interface{}) error {
	if queue == nil || consumerHandler == nil {
		return fmt.Errorf("queue and consumerHandler are required")
	}

	qCfg, ok := n.queues[queue.GetName()]
	if !ok {
		return fmt.Errorf("queue %s not found in configuration", queue.GetName())
	}

	consumerInstance, err := n.createConsumerInstance(consumerHandler)
	if err != nil {
		return err
	}

	_, err = n.Conn.QueueSubscribe(qCfg.Name, qCfg.QueueGroup, func(m *nats.Msg) {
		go n.processMessages(m, queue, consumerInstance)
	})
	if err != nil {
		return err
	}

	n.metricMonitor.SetGauge(QueueConsummerConnected.String(), 1, map[string]string{"queue": qCfg.Name})
	n.Conn.Flush()
	return nil
}

// Disconnect closes the connection to NATS.
func (n *Nats) Disconnect() error {
	n.logger.Info("Disconnecting from NATS...")
	if n.Conn != nil {
		n.Conn.Close()
	}
	return nil
}

// AddConsumer configures a queue and sets up a consumer in a single call
func (n *Nats) AddConsumer(consumer QueueConsumer) error {
	// Configure the queue
	if err := n.Configure(consumer.Queue); err != nil {
		return fmt.Errorf("failed to configure queue: %v", err)
	}

	// Register the handler in the DI container
	if err := di.GetInstance().Register(consumer.Handler); err != nil {
		n.logger.Warning("Handler already registered in DI container")
	}

	// Start consuming
	if err := n.Consume(consumer.Queue, consumer.Handler); err != nil {
		return fmt.Errorf("failed to start consuming: %v", err)
	}

	n.logger.Info("Consumer added successfully", "queue", consumer.Queue.GetName())
	return nil
}

func (n *Nats) createConsumerInstance(consumerHandler interface{}) (interface{}, error) {
	di.GetInstance().Register(consumerHandler)
	return di.GetInstance().GetByFactory(consumerHandler)
}

func (n *Nats) processMessages(m *nats.Msg, queue Queue, consumer interface{}) {
	fcm := context_manager.NewSafeContextManager()
	correlationID := m.Header.Get("x-correlation-id")
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	ctxVals := fcm.CreateValue("x-correlation-id", correlationID)
	fcm.SetValues(ctxVals, func() {
		root := n.telemetry.CreateRootSpan(fmt.Sprintf("Process message queue %s", queue.GetName()), map[string]interface{}{"correlationID": correlationID})
		defer func() { n.telemetry.EndSpan(root, nil) }()

		headers := map[string]interface{}{}
		for k, v := range m.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		err := n.callConsumerHandler(consumer, m.Data, headers)
		if err != nil {
			n.logger.Error(err.Error())
			n.metricMonitor.IncrementCounter(QueueMessageNack.String(), map[string]string{"queue": queue.GetName()})
		} else {
			n.metricMonitor.IncrementCounter(QueueMessageAck.String(), map[string]string{"queue": queue.GetName()})
		}
	})
}

func (n *Nats) callConsumerHandler(consumer interface{}, body []byte, headers map[string]interface{}) error {
	handlerType := reflect.TypeOf(consumer)
	method, exists := handlerType.MethodByName("Handler")
	if !exists {
		return fmt.Errorf("Handler method not found")
	}
	if method.Type.NumIn() != 3 {
		return fmt.Errorf("Handler method has an incorrect number of parameters")
	}

	bodyType := method.Type.In(1)
	bodyValue := reflect.New(bodyType).Interface()
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(body, &jsonMap); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}
	if err := mapper.Deserialize(jsonMap, bodyValue); err != nil {
		return fmt.Errorf("failed to deserialize message: %v", err)
	}

	params := []reflect.Value{
		reflect.ValueOf(consumer),
		reflect.ValueOf(bodyValue).Elem(),
		reflect.ValueOf(headers),
	}

	results := method.Func.Call(params)
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}
	return nil
}
