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
package context_manager

import (
	"sync"

	"github.com/jtolds/gls"
)

type ISafeContextManager interface {
	SetValues(values map[interface{}]interface{}, handler func())
	CreateValue(key string, value interface{}) map[interface{}]interface{}
	GetValue(key string) (interface{}, bool)
}
type SafeContextManager struct {
	mgr *gls.ContextManager
	mu  sync.Mutex
}

var instance ISafeContextManager
var once sync.Once

// NewSafeContextManager retorna a única instância do SafeContextManager.
func NewSafeContextManager() ISafeContextManager {
	once.Do(func() {
		instance = &SafeContextManager{
			mgr: gls.NewContextManager(),
		}
	})
	return instance
}

// SetValues define vários valores no contexto e executa a função handler.
func (scm *SafeContextManager) SetValues(values map[interface{}]interface{}, handler func()) {
	scm.mu.Lock()
	valuesToSet := gls.Values(values)
	scm.mu.Unlock()

	scm.mgr.SetValues(valuesToSet, handler)
}

// CreateValue cria e retorna um valor no formato aceito pelo gls.Values a partir de uma única chave-valor.
func (scm *SafeContextManager) CreateValue(key string, value interface{}) map[interface{}]interface{} {
	return gls.Values{key: value}
}

// GetValue recupera um valor específico do contexto.
func (scm *SafeContextManager) GetValue(key string) (interface{}, bool) {
	value, ok := scm.mgr.GetValue(key)
	return value, ok
}
