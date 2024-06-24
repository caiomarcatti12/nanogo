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
package env

import (
	"fmt"
	"os"
	"strconv"
)

type Env struct {
}

func NewEnv() IEnv {
	return &Env{}
}

func (e *Env) GetEnv(variable string, default_ ...string) string {
	value := os.Getenv(variable)

	if value == "" {
		if len(default_) > 0 {
			return default_[0]
		}
		panic(fmt.Sprintf("A variável de ambiente %s não foi definida.", variable))
	}

	return value
}

func (e *Env) GetEnvBool(variable string, default_ ...string) bool {
	value := e.GetEnv(variable, default_...)

	b, err := strconv.ParseBool(value)

	if err != nil {
		return false
	}

	return b
}
