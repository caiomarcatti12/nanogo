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
package di

import (
	"errors"
	"reflect"
	"sync"

	"github.com/caiomarcatti12/nanogo/v1/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
)

type IContainer interface {
	Register(factoryFunc interface{}) error
	RegisterAll(factoryFunc []interface{}) error
	GetByFactory(factoryFunc interface{}) (interface{}, error)
	GetByName(interfaceName string) (interface{}, error)
}

// Container manages dependency injection.
type Container struct {
	constructors map[string]interface{}
	cached       map[string]interface{}
	i18n         i18n.I18N
	log          log.ILog
}

var (
	singletonInstance IContainer
	once              sync.Once
)

func Factory(i18n i18n.I18N, log log.ILog) IContainer {
	once.Do(func() {
		singletonInstance = &Container{
			constructors: make(map[string]interface{}),
			cached:       make(map[string]interface{}),
			i18n:         i18n,
			log:          log,
		}
	})
	return singletonInstance
}

func GetInstance() IContainer {
	return singletonInstance
}

// Register adds a factory function for creating instances of a type.
func (c *Container) Register(factoryFunc interface{}) error {
	cType, err := c.getNameInterface(factoryFunc)

	if err != nil {
		return err
	}

	c.log.Trace(c.i18n.Get("di.register_new_factory", map[string]interface{}{"factory": cType}))

	c.constructors[cType] = factoryFunc

	return nil
}

// Register adds a factory function for creating instances of a type.
func (c *Container) RegisterAll(factoryFunc []interface{}) error {
	for _, factory := range factoryFunc {
		err := c.Register(factory)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetByFactory retrieves an instance by factory func.
func (c *Container) GetByFactory(factoryFunc interface{}) (interface{}, error) {
	cType, err := c.getNameInterface(factoryFunc)

	if err != nil {
		c.log.Error(c.i18n.Get("di.get_by_name_return_error", map[string]interface{}{"error": err.Error()}))
		return nil, err
	}
	return c.GetByName(cType)
}

func (c *Container) GetByName(interfaceName string) (interface{}, error) {
	c.log.Trace(c.i18n.Get("di.get_by_factory", map[string]interface{}{"factory": interfaceName}))

	instanceFunc, exists := c.constructors[interfaceName]

	if !exists {
		return nil, errors.New(c.i18n.Get("di.no_service_registered_for_type", map[string]interface{}{"factory": interfaceName}))
	}

	if instanceCached, exists := c.cached[interfaceName]; exists {
		c.log.Trace(c.i18n.Get("di.get_by_factory_cached", map[string]interface{}{"factory": interfaceName}))
		return instanceCached, nil
	}

	newInstance, err := c.factoryInstance(instanceFunc, interfaceName)
	if err != nil {
		c.log.Error(c.i18n.Get("di.factory_error", map[string]interface{}{"factory": interfaceName, "error": err.Error()}))
		return nil, err
	}

	c.log.Trace(c.i18n.Get("di.set_factory_in_cache", map[string]interface{}{"factory": interfaceName}))

	c.cached[interfaceName] = newInstance

	return newInstance, nil
}

// Get retrieves an instance of the requested type.// GetNameReturn retorna o nome completo do tipo retornado pela função de fábrica.
// Validações realizadas:
// - Verifica se factoryFunc é uma função
// - Verifica se a função de fábrica retorna exatamente um resultado
func (c *Container) getNameInterface(factoryFunc interface{}) (string, error) {
	if factoryFunc == nil {
		return "", errors.New(c.i18n.Get("di.factory_func_is_nil"))
	}

	factoryType := reflect.TypeOf(factoryFunc)
	if factoryType.Kind() != reflect.Func {
		return "", errors.New("di.factory_func_is_not_func")
	}

	if factoryType.NumOut() < 1 || factoryType.NumOut() > 2 {
		return "", errors.New("di.factory_func_has_no_one_result")
	}

	name := factoryType.Out(0).PkgPath() + "/" + factoryType.Out(0).Name()

	return name, nil
}

func (c *Container) factoryInstance(instanceFunc interface{}, interfaceName string) (interface{}, error) {
	c.log.Trace(c.i18n.Get("di.factory_new_instance", map[string]interface{}{"factory": interfaceName}))

	factoryValue := reflect.ValueOf(instanceFunc)

	parameters, err := c.resolveParameters(factoryValue, interfaceName)

	if err != nil {
		return nil, err
	}

	c.log.Trace(c.i18n.Get("di.factory_call", map[string]interface{}{"factory": interfaceName}))

	results := factoryValue.Call(parameters)

	if len(results) == 0 {
		return nil, errors.New(c.i18n.Get("di.factory_function_returned_no_results", map[string]interface{}{"factory": interfaceName}))
	}

	if len(results) == 2 && results[1].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
	}

	if !results[0].IsValid() || results[0].IsNil() {
		return nil, errors.New(c.i18n.Get("di.factory_function_returned_nil", map[string]interface{}{"factory": interfaceName}))
	}

	return results[0].Interface(), nil
}

// resolveParameters resolves the function parameters from the container.
func (c *Container) resolveParameters(fnVal reflect.Value, interfaceName string) ([]reflect.Value, error) {
	c.log.Trace(c.i18n.Get("di.resolve_parameters", map[string]interface{}{"factory": interfaceName}))

	fnType := fnVal.Type()
	var in []reflect.Value

	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i).PkgPath() + "/" + fnType.In(i).Name()
		argInstance, err := c.GetByName(argType)

		if err != nil {
			return nil, err
		}

		in = append(in, reflect.ValueOf(argInstance))
	}

	return in, nil
}

// Get retrieves an instance of the requested type.
func Get[T any]() (T, error) {
	var result T

	resultType := reflect.TypeOf((*T)(nil)).Elem()
	iType := resultType.PkgPath() + "/" + resultType.Name()

	instance, err := singletonInstance.GetByName(iType)

	if err != nil {
		return result, err
	}

	result, ok := instance.(T)
	if !ok {
		return result, errors.New("could not assert type to the expected type")
	}

	return result, nil
}
