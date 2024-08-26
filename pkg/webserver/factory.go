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
	"github.com/caiomarcatti12/nanogo/v1/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v1/pkg/di"
	"github.com/caiomarcatti12/nanogo/v1/pkg/env"
	"github.com/caiomarcatti12/nanogo/v1/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
	"github.com/caiomarcatti12/nanogo/v1/pkg/telemetry"
)

func Factory(env env.IEnv,
	logger log.ILog,
	i18n i18n.I18N,
	diContainer di.IContainer,
	telemetry telemetry.ITelemetry,
	contextManager context_manager.ISafeContextManager) IWebServer {
	return newWebServer(env, logger, i18n, diContainer, telemetry, contextManager)
}
