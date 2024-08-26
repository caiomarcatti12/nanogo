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
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Memory struct {
}

func NewOpenMemory() *Memory {
	return &Memory{}
}

func (t *Memory) initOpenMemory() {
}

// CreateRootSpan creates a new root span and stores it internally.
func (t *Memory) CreateRootSpan(name string, optionalAttrs ...interface{}) oteltrace.Span {
	return nil
}

// StartChildSpan starts a new child span using the stored root span context.
func (t *Memory) StartChildSpan(name string, optionalAttrs ...interface{}) oteltrace.Span {
	return nil
}

// EndSpan ends the given span and removes the latest context from the stack.
func (t *Memory) EndSpan(span oteltrace.Span, err error) {

}

// Shutdown shuts down the TracerProvider and ends all spans.
func (t *Memory) Shutdown() {

}

func (t *Memory) GetTraceID(span oteltrace.Span) oteltrace.TraceID {
	return span.SpanContext().TraceID()
}
