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
	"github.com/caiomarcatti12/nanogo/v2/config/i18n"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"net/http"
)

func HealthcheckHandler(ctx *HandlerContext[any]) (interface{}, error) {
	log.Debug("Healthcheck request received")

	return &APIResponse{
		Data:       i18n.Get("healthcheck"), // ou simplesmente nil se você não quiser enviar uma mensagem
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
	}, nil
}
