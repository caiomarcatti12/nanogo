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

// Service é a implementação principal de I18N
type Service struct {
	loader          TranslationLoader
	resolver        TranslationResolver
	replacer        VariableReplacer
	translations    map[string]map[string]string
	defaultLanguage string
	selectedLang    string
}

// Construtor do serviço de tradução
func NewService(loader TranslationLoader, resolver TranslationResolver, replacer VariableReplacer) *Service {
	return &Service{
		loader:          loader,
		resolver:        resolver,
		replacer:        replacer,
		defaultLanguage: "pt-br",
		selectedLang:    "pt-br",
	}
}

func (s *Service) SetLanguage(lang string) {
	if _, exists := s.translations[lang]; exists {
		s.selectedLang = strings.ToLower(lang)
	} else {
		s.selectedLang = s.defaultLanguage
	}
}

func (s *Service) GetLanguage() string {
	return s.selectedLang
}

func (s *Service) GetDefaultLanguage() string {
	return s.defaultLanguage
}

func (s *Service) LoadTranslations(path string) error {
	translations, err := s.loader.Load(path)
	if err != nil {
		return err
	}

	s.translations = translations
	return nil
}

func (s *Service) Get(key string, vars ...map[string]interface{}) string {
	var variables map[string]interface{}

	// verifica se vars foi fornecido
	if len(vars) > 0 {
		variables = vars[0]
	} else {
		variables = make(map[string]interface{})
	}

	if val, ok := s.resolver.Resolve(s.selectedLang, key); ok {
		return s.replacer.Replace(val, variables)
	}

	if val, ok := s.resolver.Resolve(s.defaultLanguage, key); ok {
		return s.replacer.Replace(val, variables)
	}

	return key
}
