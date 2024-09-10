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
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// i18n é a estrutura para o Singleton de internacionalização.
type I18NInstance struct {
	translations    map[string]map[string]interface{}
	defaultLanguage string
	selectedLang    string
}

var (
	instance *I18NInstance
	once     sync.Once
)

// GetInstance retorna a instância singleton de i18n.
func newI18n() I18N {
	once.Do(func() {
		instance = &I18NInstance{
			translations:    make(map[string]map[string]interface{}),
			defaultLanguage: "pt-br", // idioma padrão
			selectedLang:    "pt-br", // idioma padrão
		}

		instance.loadDefaultTranslations()
	})
	return instance
}

// SetLanguage define o idioma atual.
func (i *I18NInstance) SetLanguage(lang string) {
	if _, exists := i.translations[lang]; !exists {
		lang = i.defaultLanguage
	}

	i.selectedLang = strings.ToLower(lang)
}

// GetLanguage retorna o idioma atual.
func (i *I18NInstance) GetLanguage() string {
	return i.selectedLang
}

// GetLanguage retorna o idioma padrão.
func (i *I18NInstance) GetDefaultLanguage() string {
	return i.defaultLanguage
}

// LoadTranslations carrega os arquivos de tradução de um diretório.
func (i *I18NInstance) LoadTranslations(directory string) error {
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

			var newTranslations map[string]interface{}
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

// Get retorna uma tradução para uma chave específica usando notação de ponto.
func (i *I18NInstance) Get(key string, vars ...map[string]interface{}) string {
	// Primeiro, tenta no idioma selecionado
	if val := i.getValueByKey(i.translations[i.selectedLang], key); val != "" {
		return i.replaceVars(val, vars...)
	}

	// Se não encontrar, tenta no idioma padrão
	if val := i.getValueByKey(i.translations[i.defaultLanguage], key); val != "" {
		return i.replaceVars(val, vars...)
	}

	// Retorna o valor padrão ou a própria chave
	return key
}

// getValueByKey retorna o valor para uma chave usando notação de ponto.
func (i *I18NInstance) getValueByKey(translations map[string]interface{}, key string) string {
	keys := strings.Split(key, ".")
	var val interface{} = translations
	for _, k := range keys {
		if m, ok := val.(map[string]interface{}); ok {
			val = m[k]
		} else if m, ok := val.(map[interface{}]interface{}); ok {
			val = m[k]
		}
	}

	if valString, ok := val.(string); ok {
		return valString
	}

	return ""
}

// replaceVars substitui as variáveis no texto da tradução.
func (i *I18NInstance) replaceVars(value string, vars ...map[string]interface{}) string {
	for _, v := range vars {
		for k, val := range v {
			value = strings.ReplaceAll(value, "{{"+k+"}}", fmt.Sprintf("%v", val))
		}
	}
	return value
}

func (i *I18NInstance) loadDefaultTranslations() {
	currentDir, err := getCurrentDir()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}

	translationsPath := filepath.Join(currentDir, "translations")

	err = instance.LoadTranslations(translationsPath)

	if err != nil {
		fmt.Println("error loading translations:", err)
	}
}

func getCurrentDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get current directory")
	}
	return filepath.Dir(filename), nil
}
