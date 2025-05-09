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

type I18N interface {
	SetLanguage(lang string)
	GetLanguage() string
	GetDefaultLanguage() string
	LoadTranslations(path string) error
	Get(key string, vars ...map[string]interface{}) string
}

type TranslationLoader interface {
	Load(path string) (map[string]map[string]string, error)
}

type TranslationResolver interface {
	Resolve(locale, key string) (string, bool)
}

type VariableReplacer interface {
	Replace(text string, vars map[string]interface{}) string
}
