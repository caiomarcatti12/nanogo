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
package rsql

import (
	"fmt"
	"strings"
)

type Condition struct {
	Field    string
	Operator string
	Value    string
}

func Parse(rsql string) ([]Condition, error) {
	var conditions []Condition

	// Dividir a consulta RSQL em condições usando o delimitador ';'
	for _, part := range strings.Split(rsql, ";") {
		tokens := strings.Split(part, "==")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("formato RSQL inválido")
		}
		conditions = append(conditions, Condition{Field: tokens[0], Operator: "==", Value: tokens[1]})
	}

	return conditions, nil
}
