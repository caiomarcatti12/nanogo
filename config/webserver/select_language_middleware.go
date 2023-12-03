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
	"github.com/caiomarcatti12/nanogo/v2/config/context_manager"
	"github.com/caiomarcatti12/nanogo/v2/config/i18n"
	"net/http"
	"strings"
)

func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter o idioma do cabeçalho da requisição
		language := r.Header.Get("Accept-Language")
		language = getPrimaryLanguage(language)

		if language == "" {
			language = "en-us" // Define para inglês dos EUA caso não seja especificado
		}
		// Obter a instância do gerenciador i18n
		i18nInstace := i18n.GetInstance()
		i18nInstace.SetLanguage(language)

		fcm := context_manager.NewSafeContextManager()

		contextValues := fcm.CreateValue("x-accept-Language", language)

		fcm.SetValues(contextValues, func() {
			next.ServeHTTP(w, r)
		})
	})
}

func getPrimaryLanguage(acceptLanguageHeader string) string {
	// Divide a string por vírgula para obter cada idioma e seu peso.
	languages := strings.Split(acceptLanguageHeader, ",")
	if len(languages) > 0 {
		// Extrai apenas a parte do idioma do primeiro elemento.
		// Por exemplo, de "pt-BR;q=0.9" para "pt-BR".
		primaryLang := strings.Split(languages[0], ";")[0]
		return strings.TrimSpace(primaryLang)
	}
	return ""
}
