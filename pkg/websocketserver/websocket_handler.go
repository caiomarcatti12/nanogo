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
package websocketserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/caiomarcatti12/nanogo/v3/pkg/mapper"
	"github.com/caiomarcatti12/nanogo/v3/pkg/validator"
	"github.com/gorilla/websocket"
)

func (wss *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// 	tokenString := r.URL.Query().Get("token")
	// 	if tokenString == "" {
	// 		http.Error(w, "Forbidden", http.StatusForbidden)
	// 		return
	// 	}

	// 	claims := &Claims{}
	// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	// 		return jwtKey, nil
	// 	})

	// 	if err != nil || !token.Valid {
	// 		http.Error(w, "Forbidden", http.StatusForbidden)
	// 		return
	// 	}

	clientConnection, err := wss.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}

	defer clientConnection.Close()

	for {
		_, msg, err := clientConnection.ReadMessage()

		if err != nil {
			// wss.sendJSONError(clientConnection, err, 500)
			// continue
			break
		}

		payload, err := wss.parseMessage(msg)

		if err != nil {
			wss.sendJSONError(clientConnection, err, 500)
			continue
		}

		wss.debugInput(payload)

		route, err := wss.getRoute(payload)

		if err != nil {
			wss.sendJSONError(clientConnection, err, 500)
			continue
		}

		response, err := wss.callHandler(clientConnection, route, payload)

		if err != nil {
			wss.sendJSONError(clientConnection, err, 500)
			continue
		}

		err = wss.sendJSONResponse(clientConnection, response)

		if err != nil {
			wss.sendJSONError(clientConnection, err, 500)
			continue
		}
	}

	return nil, nil
}

func (wss *WebSocketServer) parseMessage(msg []byte) (Message, error) {
	var payload Message

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		return Message{}, err
	}

	return payload, nil
}

func (wss *WebSocketServer) getRoute(msg Message) (Route, error) {
	if route, exists := wss.routes[msg.Path]; exists {
		return route, nil
	}

	return Route{}, errors.New(wss.i18n.Get("websocketserver.route_not_found", map[string]interface{}{"path": msg.Path}))
}

func (wss *WebSocketServer) callHandler(clientConnection *websocket.Conn, route Route, msg Message) (interface{}, error) {
	handler, err := wss.di.GetByFactory(route.IHandler)

	if err != nil {
		return nil, err
	}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName(route.HandlerFunc)

	if !method.IsValid() {
		return nil, errors.New(wss.i18n.Get("websocketserver.method_not_found", map[string]interface{}{"path": route.Path, "method": route.HandlerFunc}))
	}

	methodType := method.Type()
	numArgs := methodType.NumIn()
	args := make([]reflect.Value, numArgs)

	for i := 0; i < numArgs; i++ {
		paramType := methodType.In(i)

		if paramType == reflect.TypeOf((*websocket.Conn)(nil)).Elem() {
			args[i] = reflect.ValueOf(clientConnection)
		} else {
			ptrToStruct := reflect.New(paramType)
			err := mapper.InjectData(msg.payload, ptrToStruct.Interface())

			if err != nil {
				return nil, errors.New(wss.i18n.Get("websocketserver.error_injecting_data", map[string]interface{}{"error": err}))
			}

			errorValidateStruct := validator.ValidateStruct(ptrToStruct.Interface())

			if errorValidateStruct != nil {
				return nil, errors.New(errorValidateStruct.Error())
			}

			args[i] = ptrToStruct.Elem()
		}
	}

	results := method.Call(args)

	if len(results) != 2 {
		return nil, nil
	}

	var result interface{}
	if !results[0].IsNil() {
		result = results[0].Interface()
	}

	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return result, err
}

func (wss *WebSocketServer) sendJSONError(clientConnection *websocket.Conn, err error, statusCode int) {
	if statusCode == 500 {
		wss.logger.Error(err.Error())
	} else {
		wss.logger.Warning(err.Error())
	}

	wss.sendJSONResponse(clientConnection, map[string]interface{}{"error": err.Error(), "status": statusCode})
}

func (wss *WebSocketServer) sendJSONResponse(clientConnection *websocket.Conn, response interface{}) error {
	responseBytes, err := json.Marshal(response)

	if err != nil {
		return err
	}

	err = clientConnection.WriteMessage(websocket.TextMessage, responseBytes)

	if err != nil {
		return err
	}

	return nil
}

func (wss *WebSocketServer) debugInput(payload Message) {
	if !wss.logInput {
		return
	}

	logData := make(map[string]interface{})

	logData["path"] = payload.Path
	logData["payload"] = payload.payload

	json, _ := json.MarshalIndent(logData, "", "  ")

	wss.logger.Trace(string(json))
}
