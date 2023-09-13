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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/errors"
	"github.com/caiomarcatti12/nanogo/v2/config/mapper"
	"net/http"
	"os"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/gorilla/mux"
)

type WebServer struct {
	router    *mux.Router
	port      string
	TLSConfig *tls.Config
	crt       string
	key       string
}

var (
	once sync.Once
	ws   *WebServer
)

func getWebServerInstance() *WebServer {
	once.Do(func() {
		router := mux.NewRouter()
		port := getPortWebServer()
		crt := env.GetEnv("SERVER_CERTIFICATE", "")
		key := env.GetEnv("SERVER_KEY", "")

		if crt != "" && key != "" {
			ws = &WebServer{
				router: router,
				port:   port,
				crt:    crt,
				key:    key,
				TLSConfig: &tls.Config{
					ClientAuth: tls.RequestClientCert,
				},
			}
		} else {
			ws = &WebServer{
				router: router,
				port:   port,
			}
		}

		ws.router.Use(CorrelationIDMiddleware)
		ws.router.Use(PayloadMiddleware)
	})

	return ws
}

func NewWebServer() *WebServer {
	ws := getWebServerInstance()

	WebserverDefaultRouter()

	return ws
}

func AddRouter[T any](method string, path string, f func(ctx *HandlerContext[T]) (interface{}, error), decoderType ...T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		contextPayload := r.Context().Value("payload")

		// type assertion aqui
		var typedPayload T
		if len(decoderType) > 0 && !isNil(decoderType[0]) {
			err := mapper.Transform(contextPayload, &typedPayload)
			if err != nil {
				http.Error(w, "Invalid payload format", http.StatusBadRequest)
				return
			}
		} else if contextPayload != nil {
			if tp, ok := contextPayload.(T); ok {
				typedPayload = tp
			} else {
				http.Error(w, "Mismatched payload type", http.StatusBadRequest)
				return
			}
		}

		data, err := f(&HandlerContext[T]{
			Payload:  typedPayload,
			RawQuery: r.URL.RawQuery,
			Headers:  r.Header,
			Request:  r,
			Response: w,
		})

		if err != nil {
			if customErr, ok := err.(*errors.CustomError); ok {
				sendJSONError(w, customErr.Message, customErr.Code)
				return
			} else {
				sendJSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Verifique se a resposta é uma APIResponse
		if apiResponse, ok := data.(*APIResponse); ok {
			// Se for uma APIResponse, envie diretamente
			for key, value := range apiResponse.Headers {
				w.Header().Set(key, value)
			}

			w.WriteHeader(apiResponse.StatusCode)

			// Verifique o cabeçalho de Content-Type
			contentType := w.Header().Get("Content-Type")
			if contentType != "" && contentType != "application/json" && apiResponse.Data != nil {
				switch v := apiResponse.Data.(type) {
				case []byte:
					w.Write(v)
				case string:
					w.Write([]byte(v))
				default:
					http.Error(w, "Unsupported data type", http.StatusInternalServerError)
					return
				}
			} else if apiResponse.Data != nil {
				json.NewEncoder(w).Encode(apiResponse.Data)
			}
		} else {
			// Caso contrário, envolva os dados em uma resposta padrão
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
		}
	}

	getWebServerInstance().router.HandleFunc(path, handlerFunc).Methods(method)
}

func (ws *WebServer) Start() {
	fmt.Printf("Servidor iniciado em localhost:%s\n", ws.port)

	if ws.crt != "" && ws.key != "" {
		server := &http.Server{
			Addr:    ":" + ws.port,
			Handler: ws.router,
			TLSConfig: &tls.Config{
				ClientAuth: tls.RequestClientCert,
			},
		}
		log.Fatal(server.ListenAndServeTLS(ws.crt, ws.key))
	} else {
		server := &http.Server{
			Addr:    ":" + ws.port,
			Handler: ws.router,
		}
		log.Fatal(server.ListenAndServe())
	}
}

func getPortWebServer() string {
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		log.Fatal("A porta do servidor não foi definida no arquivo .env")
	}

	return port
}

func sendJSONError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode and send the error message
	json.NewEncoder(w).Encode(map[string]string{
		"error": errorMessage,
	})
}

func isNil(i interface{}) bool {
	return i == nil
}
