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
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/metric_manager"
)

func WebserverDefaultRouter() {
	AddRouter("GET", "/healthcheck", HealthcheckHandler)

	if env.GetEnvBool("ENABLE_SWAGGER", "false") {
		docsRoute := env.GetEnv("SWAGGER_DOCS_ROUTE", "/swagger")
		AddRouter("GET", docsRoute, SwaggerHandler)
	}

	if env.GetEnvBool("ENABLE_PROMETHEUS", "false") {
		prometheusRoute := env.GetEnv("PROMETHEUS_ROUTE", "/metrics")

		if env.GetEnv("PROMETHEUS_TOKEN", "") != "" {
			AddRouter("GET", prometheusRoute, metric_manager.MetricsHandlerAuthenticated)
		} else {
			AddRouter("GET", prometheusRoute, metric_manager.MetricsHandler)
		}
	}
}
