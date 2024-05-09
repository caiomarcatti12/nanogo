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
)

// Container manages dependency injection.
type Container struct {
	constructors map[string]interface{}
	lock         sync.RWMutex
}

var singletonInstance *Container
var once sync.Once

// GetContainer returns the singleton instance of Container.
func GetContainer() *Container {
	once.Do(func() {
		singletonInstance = &Container{
			constructors: make(map[string]interface{}),
		}
	})
	return singletonInstance
}

// Register adds a factory function for creating instances of a type.
func (c *Container) Register(factoryFunc interface{}) error {
	cType, err := c.GetNameReturn(factoryFunc)

	if err != nil {
		return err
	}

	c.constructors[cType] = factoryFunc

	return nil
}

// Get retrieves an instance of the requested type.
func (c *Container) GetNameReturn(factoryFunc interface{}) (string, error) {
	factoryType := reflect.TypeOf(factoryFunc)
	if factoryType.Kind() != reflect.Func {
		return "", errors.New("register requires a function")
	}

	if factoryType.NumOut() != 1 {
		return "", errors.New("factory function must return exactly one result")
	}

	name := factoryType.Out(0).Name()
	return name, nil
}

// getByName retrieves an instance by type name.
func (c *Container) GetByFunctionConstructor(factoryFunc interface{}) (interface{}, error) {
	cType, err := c.GetNameReturn(factoryFunc)

	if err != nil {
		return nil, err
	}

	return c.GetByName(cType)
}

// Get retrieves an instance of the requested type.
func (c *Container) GetByInterface(interfaceType interface{}) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	iType := reflect.TypeOf(interfaceType).Elem().Name()
	instance, exists := c.constructors[iType]
	if !exists {
		return nil, errors.New("no service registered for type")
	}

	factoryValue := reflect.ValueOf(instance)
	parameters, err := c.resolveParameters(instance)
	if err != nil {
		return nil, err
	}

	results := factoryValue.Call(parameters)
	if len(results) == 0 {
		return nil, errors.New("factory function returned no results")
	}

	if !results[0].IsValid() || results[0].IsNil() {
		return nil, errors.New("factory function returned nil")
	}

	return results[0].Interface(), nil
}

// Invoke calls a function with dependencies resolved from the container.
func (c *Container) Invoke(fn interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()
	if fnType.Kind() != reflect.Func {
		return errors.New("invoke target is not a function")
	}

	parameters, err := c.resolveParameters(fn)
	if err != nil {
		return err
	}

	fnVal.Call(parameters)
	return nil
}

// resolveParameters resolves the function parameters from the container.
func (c *Container) resolveParameters(fn interface{}) ([]reflect.Value, error) {
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()
	var in []reflect.Value

	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i).Name()
		argInstance, err := c.GetByName(argType)
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(argInstance))
	}

	return in, nil
}

func (c *Container) GetByName(interfaceName string) (interface{}, error) {
	instance, exists := c.constructors[interfaceName]
	if !exists {
		return nil, errors.New("no service registered for type " + interfaceName)
	}

	factoryValue := reflect.ValueOf(instance)
	parameters, err := c.resolveParameters(instance)
	if err != nil {
		return nil, err
	}

	results := factoryValue.Call(parameters)

	if len(results) == 0 {
		return nil, errors.New("factory function returned no results")
	}

	if !results[0].IsValid() || results[0].IsNil() {
		return nil, errors.New("factory function returned nil")
	}

	return results[0].Interface(), nil
}
