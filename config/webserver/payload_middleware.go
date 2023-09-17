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
	"context"
	"encoding/json"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PayloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload := make(map[string]interface{})

		if r.Method == http.MethodGet {
			parseGetPayload(r, payload)
		}

		parseRoutePayload(r, payload)

		if isMultiPart(r) {
			if err := parseMultiPartPayload(r, w, payload); err != nil {
				return
			}
		} else {
			if err := parseJSONPayload(r, w, payload); err != nil {
				return
			}
		}

		if len(payload) > 0 {
			newCtx := context.WithValue(ctx, "payload", payload)
			next.ServeHTTP(w, r.WithContext(newCtx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func parseGetPayload(r *http.Request, payload map[string]interface{}) {
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			payload[key] = parseValue(values[0])
		}
	}
}

func parseRoutePayload(r *http.Request, payload map[string]interface{}) {
	vars := mux.Vars(r)
	for key, value := range vars {
		payload[key] = parseValue(value)
	}
}

func parseMultiPartPayload(r *http.Request, w http.ResponseWriter, payload map[string]interface{}) error {
	if err := checkMaxUploadSize(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			payload[key] = parseValue(values[0])
		}
	}

	for key := range r.MultipartForm.File {
		fileUpload, err := parseUpload(r, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		payload[key] = fileUpload
	}
	return nil
}

func parseJSONPayload(r *http.Request, w http.ResponseWriter, payload map[string]interface{}) error {
	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

func parseValue(value string) interface{} {
	if id, err := uuid.Parse(value); err == nil {
		return id
	} else if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	} else if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	} else if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	} else if timeValue, err := time.Parse(time.RFC3339, value); err == nil {
		return timeValue
	}
	return value
}

func isMultiPart(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

func checkMaxUploadSize(r *http.Request) error {
	maxUploadSizeStr := env.GetEnv("SERVER_MAX_UPLOAD_SIZE", "5")
	maxUploadSize, err := strconv.ParseInt(maxUploadSizeStr, 10, 64)
	if err != nil || maxUploadSize <= 0 {
		maxUploadSize = 5 << 20
	}

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return err
	}
	return nil
}

func parseUpload(r *http.Request, name string) (*FileUpload, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &FileUpload{
		Filename: handler.Filename,
		Size:     handler.Size,
		Content:  fileBytes,
	}, nil
}
