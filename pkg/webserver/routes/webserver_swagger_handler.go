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

// import (
// 	"github.com/caiomarcatti12/nanogo/pkg/env"
// 	"io/ioutil"
// 	"net/http"
// )

// func SwaggerHandler(ctx *HandlerContext[any]) (interface{}, error) {
// 	docsPath := env.GetEnv("SWAGGER_DOCS_PATH", ".docs/swagger/swaggeer.yaml")

// 	content, err := ioutil.ReadFile(docsPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &APIResponse{
// 		Data:       content,
// 		StatusCode: http.StatusOK,
// 		Headers:    map[string]string{"Content-Type": "application/x-yaml"},
// 	}, nil
// }
