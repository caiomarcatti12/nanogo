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
package nanogo

import (
	"fmt"

	"github.com/caiomarcatti12/nanogo/v2/config/di"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/event"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/caiomarcatti12/nanogo/v2/config/mongodb"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
	webserver_types "github.com/caiomarcatti12/nanogo/v2/config/webserver/types"
)

func Bootstrap() {
	env.LoadProvider()

	container := di.GetContainer()
	container.Register(env.FactoryEnv)
	container.Register(log.FactoryLogger)
	container.Register(event.NewInMemoryBroker)
	container.Register(event.FactoryEventDispatcher)
	container.Register(webserver.FactoryWebServer)
	container.Register(mongodb.NewMongoConnector)
	container.Register(mongodb.NewMongoORM[any])
}

func Register(factoryFunc interface{}) {

	container := di.GetContainer()
	nameFunc, err := container.GetNameReturn(factoryFunc)

	if err != nil {
		panic(err)
	}

	logger := getLogger()
	logger.Debug(fmt.Sprintf("Registro função no DI: %s", nameFunc))

	err = container.Register(factoryFunc)

	if err != nil {
		panic(err)
	}
}

func StartMongoDB() {
	container := di.GetContainer()

	instanceInterface, err := container.GetByFunctionConstructor(mongodb.NewMongoConnector)
	if err != nil {
		panic(err)
	}

	instance, ok := instanceInterface.(mongodb.IMongoConnector)
	if !ok {
		panic("could not assert type to IMongoConnector")
	}

	err = instance.ConnectMongoDB()

	if err != nil {
		panic(err)
	}
}

func WebserverAddRoutes(routes []webserver_types.Route) {
	ws := getWebServer()

	for _, route := range routes {
		ws.AddRoute(route)
	}
}
func WebserverStart() {
	ws := getWebServer()
	ws.Start()
}

func getWebServer() webserver.IWebServer {
	container := di.GetContainer()

	webserverInterface, err := container.GetByFunctionConstructor(webserver.FactoryWebServer)
	if err != nil {
		panic(err)
	}

	webserver, ok := webserverInterface.(webserver.IWebServer)
	if !ok {
		panic("could not assert type to IWebServer")
	}

	return webserver
}

func getLogger() log.ILog {
	container := di.GetContainer()

	logInterface, err := container.GetByFunctionConstructor(log.FactoryLogger)
	if err != nil {
		panic(err)
	}

	logger, ok := logInterface.(log.ILog)
	if !ok {
		panic("could not assert type to ILog")
	}

	return logger
}
