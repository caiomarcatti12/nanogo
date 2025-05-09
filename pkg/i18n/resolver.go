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

import "strings"

type DefaultResolver struct {
	translations map[string]map[string]string
}

func NewResolver(translations map[string]map[string]string) *DefaultResolver {
	return &DefaultResolver{translations: translations}
}

func (r *DefaultResolver) Resolve(locale, key string) (string, bool) {
	keys := strings.Split(key, ".")
	current := r.translations[locale]

	for i, k := range keys {
		if i == len(keys)-1 {
			val, exists := current[k]
			return val, exists
		}
		// Para simplificar, assumindo apenas uma profundidade
		return "", false
	}
	return "", false
}
