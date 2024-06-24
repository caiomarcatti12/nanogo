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
	"github.com/caiomarcatti12/nanogo/v3/src/cache"
	"github.com/caiomarcatti12/nanogo/v3/src/cli"
	"github.com/caiomarcatti12/nanogo/v3/src/db"
	"github.com/caiomarcatti12/nanogo/v3/src/di"
	"github.com/caiomarcatti12/nanogo/v3/src/env"
	"github.com/caiomarcatti12/nanogo/v3/src/event"
	"github.com/caiomarcatti12/nanogo/v3/src/log"
	"github.com/caiomarcatti12/nanogo/v3/src/metric"
	"github.com/caiomarcatti12/nanogo/v3/src/queue"
	"github.com/caiomarcatti12/nanogo/v3/src/telemetry"
	"github.com/caiomarcatti12/nanogo/v3/src/webserver"
)

func Bootstrap() {
	container := di.GetContainer()
	container.Register(env.Factory)
	container.Register(log.Factory)
	container.Register(telemetry.Factory)
	container.Register(metric.Factory)
	container.Register(cache.Factory)
	container.Register(event.Factory)
	container.Register(webserver.FactoryWebServer)
	container.Register(db.Factory)
	container.Register(queue.Factory)
	container.Register(db.NewMongoORM[any])
	container.Register(cli.Factory)
}

func Register(factoryFunc interface{}) {
	container := di.GetContainer()
	nameFunc, err := container.GetNameReturn(factoryFunc)

	if err != nil {
		panic(err)
	}

	logger, err := di.Get[log.ILog]()

	if err != nil {
		panic(err)
	}

	logger.Debugf("Register function DI: %s", nameFunc)
	err = container.Register(factoryFunc)

	if err != nil {
		panic(err)
	}
}

// func StartMongoDB() {
// 	container := di.GetContainer()

// 	instanceInterface, err := container.GetByFunctionConstructor(mongodb.NewMongoConnector)
// 	if err != nil {
// 		panic(err)
// 	}

// 	instance, ok := instanceInterface.(mongodb.IMongoConnector)
// 	if !ok {
// 		panic("could not assert type to IMongoConnector")
// 	}

// 	err = instance.ConnectMongoDB()

// 	if err != nil {
// 		panic(err)
// 	}
// }

// func WebserverAddRoutes(routes []webserver_types.Route) {
// 	ws := getWebServer()

// 	for _, route := range routes {
// 		ws.AddRoute(route)
// 	}
// }
// func WebserverStart() {
// 	ws := getWebServer()
// 	ws.Start()
// }
