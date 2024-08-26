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
	"fmt"
	"net/http"
	"sync"

	"github.com/caiomarcatti12/nanogo/v1/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v1/pkg/di"
	"github.com/caiomarcatti12/nanogo/v1/pkg/env"
	"github.com/caiomarcatti12/nanogo/v1/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
	"github.com/caiomarcatti12/nanogo/v1/pkg/telemetry"
	webserver_middleware "github.com/caiomarcatti12/nanogo/v1/pkg/webserver/middleware"
	webserver_route "github.com/caiomarcatti12/nanogo/v1/pkg/webserver/routes"
	webserver_types "github.com/caiomarcatti12/nanogo/v1/pkg/webserver/types"
	"github.com/gorilla/mux"
)

type WebServer struct {
	host           string
	port           string
	crt            string
	key            string
	logInput       bool
	logger         log.ILog
	i18n           i18n.I18N
	di             di.IContainer
	telemetry      telemetry.ITelemetry
	contextManager context_manager.ISafeContextManager
	tLSConfig      *tls.Config
	router         *mux.Router
}

var (
	once     sync.Once
	instance *WebServer
)

func newWebServer(
	env env.IEnv,
	logger log.ILog,
	i18n i18n.I18N,
	diContainer di.IContainer,
	telemetry telemetry.ITelemetry,
	contextManager context_manager.ISafeContextManager,
) IWebServer {
	once.Do(func() {
		instance = &WebServer{
			host:           env.GetEnv("WEB_SERVER_HOST", "localhost"),
			port:           env.GetEnv("WEB_SERVER_PORT", "8080"),
			crt:            env.GetEnv("WEB_SERVER_CERTIFICATE", ""),
			key:            env.GetEnv("WEB_SERVER_KEY", ""),
			logInput:       env.GetEnvBool("WEBSERVER_ACCESS_LOG", "false"),
			logger:         logger,
			i18n:           i18n,
			di:             diContainer,
			telemetry:      telemetry,
			contextManager: contextManager,
			router:         mux.NewRouter(),
			tLSConfig: &tls.Config{
				ClientAuth: tls.RequestClientCert,
			},
		}

		instance.AddMidleware(webserver_middleware.NewCorsMiddleware(env, logger, i18n))
		instance.AddMidleware(webserver_middleware.NewPayloadExtractorMiddleware(env, logger, i18n))
		instance.AddMidleware(webserver_middleware.NewCorrelationIdMiddleware(logger, i18n))
		instance.AddMidleware(webserver_middleware.NewTelemetryMiddleware(env, logger, i18n, telemetry, contextManager))

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
	ws.logger.Trace(ws.i18n.Get("webserver.add_middleware", map[string]interface{}{"middleware": middleware.GetName()}))
	ws.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware.Process(w, r, next)
		})
	})
}

func (ws *WebServer) AddRoute(route webserver_types.Route) {
	ws.logger.Trace(ws.i18n.Get("webserver.add_route", map[string]interface{}{"method": route.Method, "path": route.Path}))

	ws.di.Register(route.IHandler)

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
	ws.logger.Info(ws.i18n.Get("webserver.server_https_started", map[string]interface{}{"host": ws.host, "port": ws.port}))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", ws.host, ws.port),
		Handler: ws.router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
	}
	server.ListenAndServeTLS(ws.crt, ws.key)
}

func (ws *WebServer) startWebserverHttp() {
	ws.logger.Info(ws.i18n.Get("webserver.server_http_started", map[string]interface{}{"host": ws.host, "port": ws.port}))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", ws.host, ws.port),
		Handler: ws.router,
	}
	server.ListenAndServe()
}
