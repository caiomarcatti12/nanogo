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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// i18n é a estrutura para o Singleton de internacionalização.
type i18n struct {
	translations map[string]map[string]string
	selectedLang string
}

var (
	instance *i18n
	once     sync.Once
)

// GetInstance retorna a instância singleton de i18n.
func GetInstance() *i18n {
	once.Do(func() {
		instance = &i18n{
			translations: make(map[string]map[string]string),
			selectedLang: "en-us", // idioma padrão
		}

	})
	return instance
}

// SetLanguage define o idioma atual.
func (i *i18n) SetLanguage(lang string) {
	i.selectedLang = strings.ToLower(lang)
}

// GetLanguage retorna o idioma atual.
func (i *i18n) GetLanguage() string {
	return i.selectedLang
}

// LoadTranslations carrega os arquivos de tradução de um diretório.
func (i *i18n) LoadTranslations(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
			content, err := ioutil.ReadFile(filepath.Join(directory, f.Name()))
			if err != nil {
				return err
			}

			var newTranslations map[string]string
			err = yaml.Unmarshal(content, &newTranslations)
			if err != nil {
				return err
			}

			locale := strings.Split(f.Name(), ".")[0]

			// Mescla as novas traduções com as existentes
			if existingTranslations, exists := i.translations[locale]; exists {
				for key, value := range newTranslations {
					existingTranslations[key] = value
				}
			} else {
				i.translations[locale] = newTranslations
			}
		}
	}

	return nil
}

// Get retorna uma tradução para uma chave específica.
//
//	text := i18n.Get("hello"),
//	text := i18n.Get("hello", "Hello, world!"),
//	text := i18n.Get("greeting", "", map[string]interface{}{"name": "Alice"})
func Get(key string, vars ...map[string]interface{}) string {
	i := GetInstance()
	lang := i.GetLanguage()

	// Primeiro, tenta no idioma selecionado
	if val, ok := i.translations[lang][key]; ok {
		return i.replaceVars(val, vars...)
	}

	// Se não encontrar, tenta no idioma padrão
	if val, ok := i.translations["en-us"][key]; ok {
		return i.replaceVars(val, vars...)
	}

	// Retorna o valor padrão ou a própria chave
	return i.replaceVars(key, vars...)
}

// replaceVars substitui as variáveis no texto da tradução.
func (i *i18n) replaceVars(value string, vars ...map[string]interface{}) string {
	for _, v := range vars {
		for k, val := range v {
			value = strings.ReplaceAll(value, "{{"+k+"}}", fmt.Sprintf("%v", val))
		}
	}
	return value
}
