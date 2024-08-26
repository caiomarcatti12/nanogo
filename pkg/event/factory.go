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
package event

import (
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

func Factory(env env.IEnv, log log.ILog, i18n i18n.I18N) IEventDispatcher {
	eventDispatcher := env.GetEnv("EVENT_DISPATCHER", "IN_MEMORY")

	switch eventDispatcher {
	case "IN_MEMORY":
		return NewInMemoryBroker(log, i18n)
	default:
		panic(i18n.Get("event.provider_not_found", map[string]interface{}{"provider": eventDispatcher}))
	}
}
