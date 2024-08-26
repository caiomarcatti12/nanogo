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

	"github.com/caiomarcatti12/nanogo/v3/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	"github.com/caiomarcatti12/nanogo/v3/pkg/telemetry"
	"github.com/gorilla/mux"
)

type TelemetryMiddleware struct {
	enable         bool
	log            log.ILog
	i18n           i18n.I18N
	telemetry      telemetry.ITelemetry
	contextManager context_manager.ISafeContextManager
}

func NewTelemetryMiddleware(
	env env.IEnv,
	log log.ILog,
	i18n i18n.I18N,
	telemetry telemetry.ITelemetry,
	contextManager context_manager.ISafeContextManager,
) IMiddleware {
	return &TelemetryMiddleware{
		enable:         env.GetEnvBool("TELEMETRY_ENABLE", "false"),
		log:            log,
		i18n:           i18n,
		telemetry:      telemetry,
		contextManager: contextManager,
	}
}

func (m *TelemetryMiddleware) GetName() string {
	return "TelemetryMiddleware"
}

func (m *TelemetryMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	correlationID, _ := m.contextManager.GetValue("x-correlation-id")

	route := mux.CurrentRoute(r)
	pathTemplate, err := route.GetPathTemplate()

	if err != nil {
		http.Error(w, "Rota não encontrada", http.StatusNotFound)
		return
	}

	payload := make(map[string]interface{})

	if _, ok := r.Context().Value("payload").(map[string]interface{}); ok {
		payload = r.Context().Value("payload").(map[string]interface{})
	}

	span := m.telemetry.CreateRootSpan(r.Method+" "+pathTemplate, map[string]interface{}{"payload": payload})
	contextValues := m.contextManager.CreateValue("correlationID", correlationID)

	m.contextManager.SetValues(contextValues, func() {
		next.ServeHTTP(w, r)

		m.telemetry.EndSpan(span, nil)
		m.telemetry.Shutdown()
	})
}
