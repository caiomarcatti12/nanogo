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
	"fmt"

	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TelemetryInterface defines the methods that the Telemetry struct must implement.
type ITelemetry interface {
	CreateRootSpan(name string, optionalAttrs ...interface{}) oteltrace.Span
	StartChildSpan(name string, optionalAttrs ...interface{}) oteltrace.Span
	EndSpan(span oteltrace.Span, err error)
	GetTraceID(span oteltrace.Span) oteltrace.TraceID
	Shutdown()
}

func Factory(env env.IEnv, log log.ILog) ITelemetry {
	enable := env.GetEnvBool("TELEMETRY_ENABLE", "false")
	telemetryDispatcher := "MEMORY"
	if enable {
		telemetryDispatcher = env.GetEnv("TELEMETRY_DISPATCHER", "OPEN_TELEMETRY")
	}

	switch telemetryDispatcher {
	case "OPEN_TELEMETRY":
		return NewOpenTelemetry(env, log)
	case "MEMORY":
		return NewOpenMemory()
	default:
		panic(fmt.Sprintf("invalid telemetry dispatcher: %s", telemetryDispatcher))
	}
}
