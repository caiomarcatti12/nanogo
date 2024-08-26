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
package webserver_middleware

import (
	"net/http"
	"strings"

	"github.com/caiomarcatti12/nanogo/v1/pkg/context_manager"
	"github.com/caiomarcatti12/nanogo/v1/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
)

type SelectLanguageMiddleware struct {
	log  log.ILog
	i18n i18n.I18N
}

func NewSelectLanguageMiddleware(log log.ILog, i18n i18n.I18N) IMiddleware {
	return &SelectLanguageMiddleware{
		log:  log,
		i18n: i18n,
	}
}

func (m *SelectLanguageMiddleware) GetName() string {
	return "SelectLanguageMiddleware"
}

// Valida as origens, cabeçalhos e métodos permitidos nas requisições.

func (m *SelectLanguageMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	m.log.Trace(m.i18n.Get("webserver.middleware.resolving_SelectLanguage"))
	language := m.getPrimaryLanguage(r.Header.Get("Accept-Language"))

	m.i18n.SetLanguage(language)

	fcm := context_manager.NewSafeContextManager()

	contextValues := fcm.CreateValue("x-accept-language", language)

	fcm.SetValues(contextValues, func() {
		next.ServeHTTP(w, r)
	})

	next.ServeHTTP(w, r)
}

func (m *SelectLanguageMiddleware) getPrimaryLanguage(acceptLanguageHeader string) string {
	// Divide a string por vírgula para obter cada idioma e seu peso.
	languages := strings.Split(acceptLanguageHeader, ",")
	if len(languages) > 0 {
		// Extrai apenas a parte do idioma do primeiro elemento.
		// Por exemplo, de "pt-BR;q=0.9" para "pt-BR".
		primaryLang := strings.Split(languages[0], ";")[0]
		return strings.TrimSpace(primaryLang)
	}

	return m.i18n.GetDefaultLanguage()
}
