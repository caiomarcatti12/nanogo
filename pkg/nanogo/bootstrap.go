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
	"github.com/caiomarcatti12/nanogo/v3/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v3/pkg/db"
	"github.com/caiomarcatti12/nanogo/v3/pkg/di"
	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/event"
	"github.com/caiomarcatti12/nanogo/v3/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	"github.com/caiomarcatti12/nanogo/v3/pkg/telemetry"
	"github.com/caiomarcatti12/nanogo/v3/pkg/webserver"
	"github.com/caiomarcatti12/nanogo/v3/pkg/websocketserver"
)

func Bootstrap() {
	i18nAdapter := i18n.Factory()
	err := env.Loader(i18nAdapter)

	if err != nil {
		panic(err)
	}

	contextManager := context_manager.NewSafeContextManager()

	envAdapter := env.Factory(i18nAdapter)

	logAdapter := log.Factory(envAdapter, contextManager)

	container := di.Factory(i18nAdapter, logAdapter)

	if err := container.Register(i18n.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(env.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(log.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(di.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(webserver.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(websocketserver.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(db.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(db.NewMongoORM[any]); err != nil {
		panic(err)
	}

	if err := container.Register(context_manager.NewSafeContextManager); err != nil {
		panic(err)
	}

	if err := container.Register(telemetry.Factory); err != nil {
		panic(err)
	}

	if err := container.Register(event.Factory); err != nil {
		panic(err)
	}

	// container.Register(metric.Factory)
	// container.Register(cache.Factory)
	// container.Register(queue.Factory)
	// container.Register(cli.Factory)

}
