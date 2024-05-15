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

	"github.com/caiomarcatti12/nanogo/v2/config/context_manager"
	"github.com/google/uuid"
)

type CorrelationIdMiddleware struct {
}

func NewCorrelationIdMiddleware() IMiddleware {
	return &CorrelationIdMiddleware{}
}

func (m *CorrelationIdMiddleware) GetName() string {
	return "CorrelationIdMiddleware"
}

func (m *CorrelationIdMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	correlationID := r.Header.Get("X-Correlation-ID")
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	w.Header().Set("X-Correlation-ID", correlationID)

	fcm := context_manager.NewSafeContextManager()

	contextValues := fcm.CreateValue("x-correlation-id", correlationID)

	fcm.SetValues(contextValues, func() {
		next.ServeHTTP(w, r)
	})
}
