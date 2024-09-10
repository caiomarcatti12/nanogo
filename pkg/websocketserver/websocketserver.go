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
package websocketserver

import (
	"net/http"
	"sync"

	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/webserver"
	webserver_types "github.com/caiomarcatti12/nanogo/pkg/webserver/types"
	"github.com/gorilla/websocket"
)

var (
	once     sync.Once
	instance IWebSocketServer
)

type WebSocketServer struct {
	upgrader    *websocket.Upgrader
	routes      map[string]Route
	connections map[*websocket.Conn]string

	webserver webserver.IWebServer
	logger    log.ILog
	i18n      i18n.I18N
	di        di.IContainer

	logInput bool
}

func newWebSocketServer(
	env env.IEnv,
	logger log.ILog,
	i18n i18n.I18N,
	ws webserver.IWebServer,
	di di.IContainer,
) IWebSocketServer {
	once.Do(func() {
		instance = &WebSocketServer{
			upgrader: &websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			routes:      make(map[string]Route, 0),
			connections: make(map[*websocket.Conn]string),
			webserver:   ws,
			logger:      logger,
			i18n:        i18n,
			di:          di,
			logInput:    env.GetEnvBool("WEBSOCKET_SERVER_LOG_INPUT", "false"),
		}

		instance.AddRoute(Route{
			Path:        "/ping",
			IHandler:    NewPingController,
			HandlerFunc: "Handler",
		})

	})

	return instance
}

func (ws *WebSocketServer) Start() {
	ws.webserver.AddRoute(webserver_types.Route{
		Method:      "GET",
		Path:        "/ws",
		IHandler:    newWebSocketServer,
		HandlerFunc: "HandleConnections",
	})
	ws.webserver.Start()
}

func (wss *WebSocketServer) AddRoute(route Route) {
	wss.logger.Trace(wss.i18n.Get("websocketserver.add_route", map[string]interface{}{"path": route.Path}))

	wss.di.Register(route.IHandler)

	wss.routes[route.Path] = route
}
