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
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/src/di"
	"github.com/caiomarcatti12/nanogo/v2/src/env"
	"github.com/caiomarcatti12/nanogo/v2/src/log"
	webserver_middleware "github.com/caiomarcatti12/nanogo/v2/src/webserver/middleware"
	webserver_route "github.com/caiomarcatti12/nanogo/v2/src/webserver/routes"
	webserver_types "github.com/caiomarcatti12/nanogo/v2/src/webserver/types"
	"github.com/gorilla/mux"
)

type WebServer struct {
	port      string
	crt       string
	key       string
	logInput  bool
	logger    log.ILog
	tLSConfig *tls.Config
	router    *mux.Router
}

var (
	once     sync.Once
	instance *WebServer
)

func NewWebServer(env env.IEnv, logger log.ILog) IWebServer {
	once.Do(func() {
		instance = &WebServer{
			port:     env.GetEnv("SERVER_PORT", "8080"),
			crt:      env.GetEnv("SERVER_CERTIFICATE", ""),
			key:      env.GetEnv("SERVER_KEY", ""),
			logInput: env.GetEnvBool("SERVER_LOG_INPUT", "false"),
			logger:   logger,
			router:   mux.NewRouter(),
			tLSConfig: &tls.Config{
				ClientAuth: tls.RequestClientCert,
			},
		}

		instance.AddMidleware(webserver_middleware.NewCorsMiddleware(env))
		instance.AddMidleware(webserver_middleware.NewPayloadExtractorMiddleware(env))
		instance.AddMidleware(webserver_middleware.NewCorrelationIdMiddleware())
		instance.AddMidleware(webserver_middleware.NewTelemetryMiddleware())

		instance.AddRoute(webserver_types.Route{
			Path:        "/healthcheck",
			Method:      http.MethodGet,
			IHandler:    webserver_route.NewHealthCheckController,
			HandlerFunc: "Handler",
		})
	})

	return instance
}

func (ws *WebServer) AddMidleware(middleware webserver_middleware.IMiddleware) {
	ws.logger.Debug("Adicionando middleware " + middleware.GetName())
	ws.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware.Process(w, r, next)
		})
	})
}

func (ws *WebServer) AddRoute(route webserver_types.Route) {
	ws.logger.Debug("Adicionando rota " + route.Method + " " + route.Path)

	di.GetContainer().Register(route.IHandler)

	ws.router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		ws.Handler(w, r, route)
	}).Methods(route.Method)

	// Adiciona automaticamente suporte para método OPTIONS para cada rota.
	ws.router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
}

func (ws *WebServer) Start() {
	if ws.crt != "" && ws.key != "" {
		ws.startWebserverHttps()
	} else {
		ws.startWebserverHttp()
	}
}

func (ws *WebServer) startWebserverHttps() {
	ws.logger.Info("Servidor (HTTPS) iniciado em localhost:" + ws.port)

	server := &http.Server{
		Addr:    ":" + ws.port,
		Handler: ws.router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
	}
	server.ListenAndServeTLS(ws.crt, ws.key)
}

func (ws *WebServer) startWebserverHttp() {
	ws.logger.Info("Servidor (HTTP) iniciado em localhost:" + ws.port)
	server := &http.Server{
		Addr:    ":" + ws.port,
		Handler: ws.router,
	}
	server.ListenAndServe()
}
