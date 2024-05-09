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
package webserver_route

import (
	"net/http"

	webserver_types "github.com/caiomarcatti12/nanogo/v2/config/webserver/types"
)

type IHealthCheckController interface {
	Handler() (interface{}, error)
}

type HealthCheckController struct {
}

func NewHealthCheckController() IHealthCheckController {
	return &HealthCheckController{}
}

type Teste struct {
	Teste string
}

func (hc *HealthCheckController) Handler() (interface{}, error) {
	return webserver_types.APIResponse{
		Data:       "OK",
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
	}, nil
}
