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
	"strings"

	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

type CorsMiddleware struct {
	allowedOrigins []string
	allowedHeaders []string
	allowedMethods []string
	log            log.ILog
	i18n           i18n.I18N
}

func NewCorsMiddleware(env env.IEnv, log log.ILog, i18n i18n.I18N) IMiddleware {
	return &CorsMiddleware{
		allowedOrigins: strings.Split(env.GetEnv("WEBSERVER_ORIGINS", "*"), ","),
		allowedHeaders: strings.Split(env.GetEnv("WEBSERVER_HEADERS", "Content-Type"), ","),
		allowedMethods: strings.Split(env.GetEnv("WEBSERVER_METHODS", "GET,POST,PUT,DELETE"), ","),
		log:            log,
		i18n:           i18n,
	}
}

func (m *CorsMiddleware) GetName() string {
	return "CorsMiddleware"
}

// Valida as origens, cabeçalhos e métodos permitidos nas requisições.

func (m *CorsMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	m.log.Trace(m.i18n.Get("webserver.middleware.resolving_cors"))
	origin := r.Header.Get("Origin")

	// Verifica se a origem da requisição está na lista de origens permitidas.
	if m.originAllowed(origin, m.allowedOrigins) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.allowedMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.allowedHeaders, ","))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	next.ServeHTTP(w, r)
}

// Retorna true se a origem for permitida, false caso contrário.
func (m *CorsMiddleware) originAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*") {
		return true
	}

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}
