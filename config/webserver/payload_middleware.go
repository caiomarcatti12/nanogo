package webserver

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func PayloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := make(map[string]interface{})

		// Se o método for GET, capture os parâmetros de consulta.
		if r.Method == http.MethodGet {
			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					val := values[0]

					// Se houver múltiplos valores para a mesma chave, armazene apenas o primeiro valor
					// Tente converter o valor para um inteiro
					if intValue, err := strconv.Atoi(val); err == nil {
						payload[key] = intValue
					} else {
						// Tente converter o valor para um float
						if floatValue, err := strconv.ParseFloat(val, 64); err == nil {
							payload[key] = floatValue
						} else {
							payload[key] = val
						}
					}
				}
			}
		}

		// Para todos os métodos, capture parâmetros da rota.
		vars := mux.Vars(r)
		for key, value := range vars {
			// Tente converter o valor para UUID
			if id, err := uuid.Parse(value); err == nil {
				payload[key] = id
			} else {
				payload[key] = value
			}
		}

		// Se o corpo da requisição estiver vazio, não tente decodificá-lo. Caso contrário, decodifique-o.
		if r.Body != http.NoBody {
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && err != io.EOF {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if payload != nil {
			newCtx := context.WithValue(ctx, "payload", payload)
			next.ServeHTTP(w, r.WithContext(newCtx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
