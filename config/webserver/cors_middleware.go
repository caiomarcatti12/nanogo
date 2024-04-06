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
package webserver

import (
	"net/http"
	"strings"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
)

// CorsMiddleware cria um middleware para tratar requisições CORS.
// Valida as origens, cabeçalhos e métodos permitidos nas requisições.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recupera a lista de origens permitidas do ambiente.
		allowedOrigins := strings.Split(env.GetEnv("WEBSERVER_ORIGINS", "*"), ",")
		allowedHeaders := strings.Split(env.GetEnv("WEBSERVER_HEADERS", "Content-Type"), ",")
		allowedMethods := strings.Split(env.GetEnv("WEBSERVER_METHODS", "GET,POST,PUT,DELETE"), ",")

		origin := r.Header.Get("Origin")

		// Verifica se a origem da requisição está na lista de origens permitidas.
		if originAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Se for uma requisição preflight, retorna apenas os cabeçalhos necessários.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// originAllowed verifica se a origem está na lista de origens permitidas.
// Retorna true se a origem for permitida, false caso contrário.
func originAllowed(origin string, allowedOrigins []string) bool {
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