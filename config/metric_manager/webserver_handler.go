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

package metric_manager

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MetricsHandler(ctx *webserver.HandlerContext[interface{}]) (interface{}, error) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Response, ctx.Request)

	return &webserver.APIResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func MetricsHandlerAuthenticated(ctx *webserver.HandlerContext[interface{}]) (interface{}, error) {
	tokenSecret := env.GetEnv("PROMETHEUS_TOKEN", "")

	if tokenSecret != "" {
		providedToken := ctx.Request.Header.Get("Authorization")

		if providedToken != "Bearer "+tokenSecret {
			return &webserver.APIResponse{
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
	}

	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Response, ctx.Request)

	return &webserver.APIResponse{
		StatusCode: http.StatusOK,
	}, nil
}
