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
package i18n

import (
	"fmt"
	"strings"
)

type DefaultReplacer struct{}

func NewReplacer() *DefaultReplacer {
	return &DefaultReplacer{}
}

func (r *DefaultReplacer) Replace(text string, vars map[string]interface{}) string {
	for k, v := range vars {
		text = strings.ReplaceAll(text, "{{"+k+"}}", fmt.Sprintf("%v", v))
	}
	return text
}
