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
package telemetry

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/caiomarcatti12/nanogo/v1/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v1/pkg/env"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	tracer            oteltrace.Tracer
	tp                *trace.TracerProvider
	mu                sync.Mutex
	contextStack      []context.Context
	serviceName       string
	endpointTelemetry string
	logger            log.ILog
}

var instance *Telemetry
var once sync.Once

func NewOpenTelemetry(env env.IEnv, logger log.ILog) *Telemetry {
	once.Do(func() {
		instance = &Telemetry{
			serviceName:       env.GetEnv("APP_NAME", "UNDEFINED"),
			endpointTelemetry: env.GetEnv("TELEMETRY_ENDPOINT", "localhost:4317"),
			logger:            logger,
		}
		instance.initOpenTelemetry()
	})
	return instance
}

func (t *Telemetry) initOpenTelemetry() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(t.endpointTelemetry), otlptracegrpc.WithInsecure())
	if err != nil {
		t.logger.Fatalf("Failed to create exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(t.serviceName),
		),
	)
	if err != nil {
		t.logger.Fatalf("Failed to create resource: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	t.tp = tp
	t.tracer = otel.Tracer("nanogo")
}

// CreateRootSpan creates a new root span and stores it internally.
func (t *Telemetry) CreateRootSpan(name string, optionalAttrs ...interface{}) oteltrace.Span {
	t.mu.Lock()
	defer t.mu.Unlock()

	attributes := t.extractAttributes(optionalAttrs...)
	newCtx, span := t.tracer.Start(context.Background(), name, oteltrace.WithAttributes(attributes...))
	t.contextStack = append(t.contextStack, newCtx)

	return span
}

// StartChildSpan starts a new child span using the stored root span context.
func (t *Telemetry) StartChildSpan(name string, optionalAttrs ...interface{}) oteltrace.Span {
	t.mu.Lock()
	defer t.mu.Unlock()

	attributes := t.extractAttributes(optionalAttrs...)
	parentCtx := t.contextStack[len(t.contextStack)-1]
	newCtx, span := t.tracer.Start(parentCtx, name, oteltrace.WithAttributes(attributes...))
	t.contextStack = append(t.contextStack, newCtx)

	return span
}

// EndSpan ends the given span and removes the latest context from the stack.
func (t *Telemetry) EndSpan(span oteltrace.Span, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	span.End()

	if len(t.contextStack) > 0 {
		t.contextStack = t.contextStack[:len(t.contextStack)-1]
	}
}

// Shutdown shuts down the TracerProvider and ends all spans.
func (t *Telemetry) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := t.tp.Shutdown(ctx); err != nil {
		t.logger.Fatalf("Failed to shutdown TracerProvider: %v", err)
	}
}

func (t *Telemetry) extractAttributes(optionalAttrs ...interface{}) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	if len(optionalAttrs) > 0 {
		if attrs, ok := optionalAttrs[0].(map[string]interface{}); ok {
			attributes = t.extractAttributesRecursively(attrs, "")
		}
	}

	fcm := context_manager.NewSafeContextManager()
	correlationID, ok := fcm.GetValue("x-correlation-id")

	if !ok {
		return attributes
	}

	if correlationIDString, ok := correlationID.(string); ok {
		attributes = append(attributes, attribute.String("x-correlation-id", correlationIDString))
	}

	return attributes
}

// extractAttributesRecursively extracts attributes from a map recursively, adding a prefix to each key.
func (t *Telemetry) extractAttributesRecursively(attrs map[string]interface{}, prefix string) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	for k, v := range attrs {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		switch v.(type) {
		case string:
			attributes = append(attributes, attribute.String(fullKey, v.(string)))
		case int:
			attributes = append(attributes, attribute.Int(fullKey, v.(int)))
		case float64:
			attributes = append(attributes, attribute.Float64(fullKey, v.(float64)))
		case map[string]interface{}:
			attributes = append(attributes, t.extractAttributesRecursively(v.(map[string]interface{}), fullKey)...)
		case uuid.UUID:
			attributes = append(attributes, attribute.String(fullKey, v.(uuid.UUID).String()))
		default:
			// Convert to JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				// Handle error if needed
				continue
			}
			// Convert JSON back to map[string]interface{}
			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			if err != nil {
				// Handle error if needed
				continue
			}
			attributes = append(attributes, t.extractAttributesRecursively(jsonMap, fullKey)...)

		}
	}

	return attributes
}

func (t *Telemetry) GetTraceID(span oteltrace.Span) oteltrace.TraceID {
	return span.SpanContext().TraceID()
}
