package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/caiomarcatti12/nanogo/v2/config/errors"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"os"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/gorilla/mux"
)

type WebServer struct {
	router *mux.Router
	port   string
}

var (
	once sync.Once
	ws   *WebServer
)

func getWebServerInstance() *WebServer {
	once.Do(func() {
		router := mux.NewRouter()
		port := getPortWebServer()

		ws = &WebServer{
			router: router,
			port:   port,
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

func AddRouter[T any](method string, path string, f func(ctx *HandlerContext) (interface{}, error), decoderType ...T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		contextPayload := r.Context().Value("payload")

		// type assertion aqui
		if len(decoderType) > 0 && decoderType[0] != nil {
			err := mapstructure.Decode(contextPayload, decoderType[0])

			if err != nil {
				http.Error(w, "Invalid payload format", http.StatusBadRequest)
				return
			}

			contextPayload = decoderType[0]
		}

		data, err := f(&HandlerContext{
			Payload:  contextPayload,
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
			if apiResponse.Data != nil {
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
	port := getPortWebServer()

	fmt.Printf("Servidor iniciado em localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, ws.router))
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
