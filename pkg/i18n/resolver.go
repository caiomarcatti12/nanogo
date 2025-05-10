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

type DefaultResolver struct {
	translations map[string]map[string]string
}

func NewResolver(translations map[string]map[string]string) *DefaultResolver {
	return &DefaultResolver{translations: translations}
}

func (r *DefaultResolver) Resolve(locale, key string) (string, bool) {
	// Obtenha o mapa de traduções para o locale
	current, exists := r.translations[locale]
	if !exists {
		return "", false // Retorna falso se o locale não existir
	}

	// Verifica se a chave 'key' existe no mapa 'current'
	val, exists := current[key]
	return val, exists
}
