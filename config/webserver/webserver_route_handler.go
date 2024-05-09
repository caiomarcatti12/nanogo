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
package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/caiomarcatti12/nanogo/v2/config/di"
	"github.com/caiomarcatti12/nanogo/v2/config/errors"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/caiomarcatti12/nanogo/v2/config/mapper"
	webserver_types "github.com/caiomarcatti12/nanogo/v2/config/webserver/types"
)

func (ws *WebServer) Handler(w http.ResponseWriter, r *http.Request, route webserver_types.Route) {
	payload := make(map[string]interface{})

	if _, ok := r.Context().Value("payload").(map[string]interface{}); ok {
		payload = r.Context().Value("payload").(map[string]interface{})
	}

	if payload == nil {
		payload = make(map[string]interface{})
	}

	if ws.logInput {
		ws.debugInput(w, r, payload)
	}

	data, err := ws.executeHandler(route, payload)

	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			ws.sendJSONError(w, customErr.Message, customErr.Code)
			return
		} else {
			ws.sendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if apiResponse, ok := data.(webserver_types.APIResponse); ok {
		for key, value := range apiResponse.Headers {
			w.Header().Set(key, value)
		}

		w.WriteHeader(apiResponse.StatusCode)

		contentType := w.Header().Get("Content-Type")
		if contentType != "" && contentType != "application/json" && apiResponse.Data != nil {
			switch v := apiResponse.Data.(type) {
			case []byte:
				w.Write(v)
			case string:
				w.Write([]byte(v))
			default:
				ws.sendJSONError(w, "Unsupported data type", http.StatusInternalServerError)
				return
			}
		} else if apiResponse.Data != nil {
			json.NewEncoder(w).Encode(apiResponse.Data)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
func (ws *WebServer) executeHandler(route webserver_types.Route, contextPayload map[string]interface{}) (interface{}, error) {
	handler, err := di.GetContainer().GetByFunctionConstructor(route.IHandler)

	if err != nil {
		return nil, err
	}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName(route.HandlerFunc)

	if !method.IsValid() {
		return nil, fmt.Errorf("handlerFunc not found")
	}

	methodType := method.Type()
	numArgs := methodType.NumIn()

	// Criando slice de reflect.Value para os argumentos
	args := make([]reflect.Value, numArgs)

	if numArgs > 0 {
		for i := 0; i < numArgs; i++ {
			paramType := methodType.In(0)
			instanceParam := reflect.New(paramType).Interface()

			err := mapper.Transform(contextPayload, instanceParam)

			if err != nil {
				args[i] = reflect.Zero(paramType)
			} else {
				args[i] = reflect.ValueOf(instanceParam).Elem()
			}
		}
	}

	results := method.Call(args)

	if len(results) != 2 {
		return nil, fmt.Errorf("handler must return two values")
	}

	// Processando o resultado
	var result interface{}
	if !results[0].IsNil() {
		result = results[0].Interface()
	}

	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return result, err
}

func (ws *WebServer) sendJSONError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode == 500 {
		ws.logger.Error(errorMessage)
	} else {
		ws.logger.Warning(errorMessage)
	}

	// Encode and send the error message
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errorMessage,
	})
}

func (ws *WebServer) debugInput(w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	logData := make(map[string]interface{})

	logData["method"] = r.Method
	logData["path"] = r.URL.Path
	logData["headers"] = r.Header
	logData["payload"] = payload

	json, _ := json.MarshalIndent(logData, "", "  ")

	instanceInterface, err := di.GetContainer().GetByFunctionConstructor(log.FactoryLogger)

	if err != nil {
		return
	}

	instance, ok := instanceInterface.(log.ILog)
	if !ok {
		panic("could not assert type to ILog")
	}

	instance.Debug(string(json))
}
