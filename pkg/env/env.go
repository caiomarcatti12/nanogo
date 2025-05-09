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
	"os"
	"strconv"
	"strings"

	"github.com/caiomarcatti12/nanogo/pkg/i18n"
)

type Env struct {
	i18n i18n.I18N
}

func NewEnv(i18n i18n.I18N) IEnv {
	return &Env{
		i18n: i18n,
	}
}

// GetEnv obtém uma variável de ambiente do sistema.
// Valida explicitamente se o nome da variável é válido (não vazio ou espaço em branco).
// Caso não encontrada e sem valor padrão fornecido, dispara panic.
func (e *Env) GetEnv(variable string, defaultValue ...string) string {
	if len(strings.TrimSpace(variable)) == 0 {
		panic(e.i18n.Get("env.invalid_variable_name", map[string]interface{}{"variable": variable}))
	}

	value := os.Getenv(variable)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		panic(e.i18n.Get("env.not_found", map[string]interface{}{"variable": variable}))
	}

	return value
}

// GetEnvBool obtém uma variável booleana do ambiente do sistema.
// Valida explicitamente o nome da variável e realiza conversão segura para bool.
// Caso ocorra erro na conversão e nenhum valor padrão é fornecido, dispara panic.
func (e *Env) GetEnvBool(variable string, defaultValue ...string) bool {
	value := e.GetEnv(variable, defaultValue...)

	b, err := strconv.ParseBool(value)

	if err != nil {
		return false
	}

	return b
}
