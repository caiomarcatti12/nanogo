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
	"github.com/caiomarcatti12/nanogo/v2/config/di"
	"github.com/caiomarcatti12/nanogo/v2/config/eda"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
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
	container.Register(eda.FactoryEventDispatcher)
	container.Register(webserver.FactoryWebServer)
	container.Register(mongodb.NewMongoConnector)
	container.Register(mongodb.NewMongoORM[any])
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

func StartWebServer(routes []webserver_types.Route) {
	container := di.GetContainer()

	webserverInterface, err := container.GetByFunctionConstructor(webserver.FactoryWebServer)
	if err != nil {
		panic(err)
	}

	webserver, ok := webserverInterface.(webserver.IWebServer)
	if !ok {
		panic("could not assert type to IWebServer")
	}

	for _, route := range routes {
		webserver.AddRoute(route)
	}

	webserver.Start()
}
