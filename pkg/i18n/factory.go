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
	"github.com/caiomarcatti12/nanogo/pkg/util"
	"github.com/caiomarcatti12/nanogo/pkg/yaml"
)

func Factory() (I18N, error) {
	translationsPath := util.GetAbsolutePath("pkg/i18n/translations")

	yamlLoader := yaml.NewYAMLLoader()
	translations, err := yamlLoader.Load(translationsPath)
	if err != nil {
		return nil, err
	}

	resolver := NewResolver(translations)
	replacer := NewReplacer()
	service := NewService(yamlLoader, resolver, replacer)

	return service, nil
}
