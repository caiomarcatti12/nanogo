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
package env

import (
	"errors"
	"os"

	"github.com/caiomarcatti12/nanogo/v1/pkg/i18n"
)

func Loader(i18n i18n.I18N) error {
	provider := os.Getenv("ENV_PROVIDER")

	if provider == "" {
		provider = "ENV_FILE"
	}

	// Retrieve the provider function from the map
	providerFunc, exists := providers[provider]
	if !exists {
		return errors.New(i18n.Get("env.provider_not_found", map[string]interface{}{"provider": provider}))
	}

	// Assert that the provider is a function
	providerFuncTyped, ok := providerFunc.(func() error)
	if !ok {
		return errors.New(i18n.Get("env.provider_not_valid_function", map[string]interface{}{"provider": provider}))
	}

	return providerFuncTyped()
}
