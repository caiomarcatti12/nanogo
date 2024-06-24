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
package webserver_middleware

import (
	"net/http"

	"github.com/caiomarcatti12/nanogo/v2/src/telemetry"
)

type TelemetryMiddleware struct {
}

func NewTelemetryMiddleware() IMiddleware {
	return &TelemetryMiddleware{}
}

func (m *TelemetryMiddleware) GetName() string {
	return "TelemetryMiddleware"
}

func (m *TelemetryMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	telemetry := telemetry.NewOpenTelemetry()
	// fcm := context_manager.NewSafeContextManager()

	correlationID := r.Header.Get("X-Correlation-ID")
	span := telemetry.CreateRootSpan(r.Method+" "+r.URL.Path, map[string]interface{}{"correlationID": correlationID})
	// contextValues := fcm.CreateValue("correlationID", correlationID)

	// fcm.SetValues(contextValues, func() {
	next.ServeHTTP(w, r)

	telemetry.EndSpan(span, nil)
	telemetry.Shutdown()
	// })
}
