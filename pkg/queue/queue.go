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
	"fmt"

	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	"github.com/caiomarcatti12/nanogo/v3/pkg/metric"
	"github.com/caiomarcatti12/nanogo/v3/pkg/telemetry"
)

type IConsumer[T any] interface {
	Handler(body T, headers map[string]interface{}) error
}

type IConfig interface {
	Get() interface{}
}

type Queue interface {
	GetName() string
}

type IQueue interface {
	Connect() error
	Configure(args ...interface{}) error
	Consume(queue Queue, consumerHandler interface{}) error
	Publish(exchange string, routingKey string, body interface{}) error
	Disconnect() error
}

type QueueMetrics string

const (
	QueueManagerConnected   QueueMetrics = "queue_manager_connected"
	QueueExchangeCreated    QueueMetrics = "queue_exchange_created"
	QueueCreated            QueueMetrics = "queue_created"
	QueueBindExchange       QueueMetrics = "queue_binded"
	QueueConsummerConnected QueueMetrics = "queue_consummer_connected"
	QueueMessageAck         QueueMetrics = "queue_messages_ack"
	QueueMessageNack        QueueMetrics = "queue_messages_nack"
	QueueMessagePublish     QueueMetrics = "queue_messages_publish"
)

func (qm QueueMetrics) String() string {
	return string(qm)
}

func Factory(env env.IEnv, logger log.ILog, metricMonitor metric.IMetric, telemetry telemetry.ITelemetry) IQueue {
	logger.Info("Creating queue provider...")

	logger.Info("Creating metric monitor queues...")
	metricMonitor.CreateMetric(metric.Gauge, QueueManagerConnected.String(), "Indicates if queue manager is connected", []string{"provider"})
	metricMonitor.CreateMetric(metric.Gauge, QueueExchangeCreated.String(), "Indicates if exchange is created", []string{"exchange"})
	metricMonitor.CreateMetric(metric.Gauge, QueueCreated.String(), "Indicates if queue is connected", []string{"queue"})
	metricMonitor.CreateMetric(metric.Gauge, QueueBindExchange.String(), "Indicates if queue is binded", []string{"exchange", "queue"})
	metricMonitor.CreateMetric(metric.Gauge, QueueConsummerConnected.String(), "Indicates if queue consumer is connected", []string{"queue"})
	metricMonitor.CreateMetric(metric.Counter, QueueMessageAck.String(), "Indicates quantity messages consumed successfully", []string{"queue"})
	metricMonitor.CreateMetric(metric.Counter, QueueMessageNack.String(), "Indicates quantity messages consumed with error", []string{"queue"})
	metricMonitor.CreateMetric(metric.Counter, QueueMessagePublish.String(), "Indicates quantity messages published", []string{"queue", "routing_key"})

	provider := env.GetEnv("QUEUE_PROVIDER", "RABBITMQ")

	switch provider {
	case "RABBITMQ":
		instance, err := NewInstanceRabbitmq(env, logger, metricMonitor, telemetry)

		if err != nil {
			panic(err)
		}

		return instance
	default:
		panic(fmt.Errorf("queue provider %s not found", provider))
	}
}
