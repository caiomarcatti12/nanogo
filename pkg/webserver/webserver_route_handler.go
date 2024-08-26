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
	"net/http"
	"reflect"

	"github.com/caiomarcatti12/nanogo/v3/pkg/errors"
	"github.com/caiomarcatti12/nanogo/v3/pkg/mapper"
	"github.com/caiomarcatti12/nanogo/v3/pkg/types"
	"github.com/caiomarcatti12/nanogo/v3/pkg/validator"
	webserver_types "github.com/caiomarcatti12/nanogo/v3/pkg/webserver/types"
	"github.com/gorilla/websocket"
	"github.com/mozillazg/go-httpheader"
)

func (ws *WebServer) Handler(w http.ResponseWriter, r *http.Request, route webserver_types.Route) {
	ws.logger.Trace(ws.i18n.Get("webserver.execute_handler", map[string]interface{}{"method": r.Method, "path": r.URL.Path}))
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

	data, err := ws.callHandler(w, r, route, payload, r.Header)

	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			ws.sendJSONError(w, customErr.Message, customErr.Code)
			return
		} else {
			ws.sendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if apiResponse, ok := data.(types.Response); ok {
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
	} else if !ws.isWebSocket(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
	}
}
func (ws *WebServer) callHandler(w http.ResponseWriter, r *http.Request, route webserver_types.Route, contextPayload map[string]interface{}, contextHeaders http.Header) (response interface{}, err error) {
	handler, err := ws.di.GetByFactory(route.IHandler)

	if err != nil {
		return nil, err
	}

	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	var structName string
	if handlerType.Kind() == reflect.Ptr {
		structName = handlerType.Elem().Name()
	} else {
		structName = handlerType.Name()
	}

	span := ws.telemetry.StartChildSpan(structName + "::" + route.HandlerFunc)
	defer (func() { ws.telemetry.EndSpan(span, err) })()

	method := handlerValue.MethodByName(route.HandlerFunc)

	if !method.IsValid() {
		return nil, errors.InternalServerError(ws.i18n.Get("webserver.method_not_found", map[string]interface{}{"method": route.HandlerFunc, "path": route.Path}))
	}

	methodType := method.Type()
	numArgs := methodType.NumIn()
	args := make([]reflect.Value, numArgs)

	for i := 0; i < numArgs; i++ {
		paramType := methodType.In(i)

		if paramType == reflect.TypeOf((*http.ResponseWriter)(nil)).Elem() {
			args[i] = reflect.ValueOf(w)
		} else if paramType == reflect.TypeOf((*http.Request)(nil)) {
			args[i] = reflect.ValueOf(r)
		} else {
			ptrToStruct := reflect.New(paramType)
			err := mapper.InjectData(contextPayload, ptrToStruct.Interface())

			if err != nil {
				return nil, errors.InternalServerError(ws.i18n.Get("webserver.error_injecting_data", map[string]interface{}{"error": err}))
			}

			err = httpheader.Decode(contextHeaders, ptrToStruct.Interface())
			if err != nil {
				return nil, errors.InternalServerError(ws.i18n.Get("webserver.error_decoding_headers", map[string]interface{}{"error": err}))
			}

			errorValidateStruct := validator.ValidateStruct(ptrToStruct.Interface())

			if errorValidateStruct != nil {
				return nil, errorValidateStruct
			}

			args[i] = ptrToStruct.Elem()
		}
	}

	results := method.Call(args)

	if len(results) == 1 {
		if results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Chan {
			if !results[0].IsNil() {
				err = results[0].Interface().(error)
			}

			if err != nil {
				return nil, err
			}

			return nil, nil
		} else {
			err = results[0].Interface().(error)

			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}

	var result interface{}
	if results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Chan {
		if !results[0].IsNil() {
			result = results[0].Interface()
		}
	} else {
		result = results[0].Interface()
	}

	if results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Chan {
		if !results[1].IsNil() {
			err = results[1].Interface().(error)
		}
	} else {
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

	ws.logger.Trace(string(json))
}

func (ws *WebServer) isWebSocket(r *http.Request) bool {
	return websocket.IsWebSocketUpgrade(r)
}
