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
package webserver_middleware

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caiomarcatti12/nanogo/v3/pkg/env"
	"github.com/caiomarcatti12/nanogo/v3/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/v3/pkg/log"
	webserver_types "github.com/caiomarcatti12/nanogo/v3/pkg/webserver/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PayloadExtractorMiddleware struct {
	maxUploadSize string
	log           log.ILog
	i18n          i18n.I18N
}

func NewPayloadExtractorMiddleware(env env.IEnv, log log.ILog, i18n i18n.I18N) IMiddleware {
	return &PayloadExtractorMiddleware{
		maxUploadSize: env.GetEnv("WEB_SERVER_MAX_UPLOAD_SIZE", "5"),
		log:           log,
		i18n:          i18n,
	}
}

func (m *PayloadExtractorMiddleware) GetName() string {
	return "PayloadExtractorMiddleware"
}

func (m *PayloadExtractorMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.Handler) {
	m.log.Trace(m.i18n.Get("webserver.middleware.extracting_payload"))

	payload := make(map[string]interface{})

	if r.Method == http.MethodGet {
		m.parseGetPayload(r, payload)
	}

	m.parseRoutePayload(r, payload)

	if m.isMultiPart(r) {
		if err := m.parseMultiPartPayload(r, w, payload); err != nil {
			return
		}
	} else {
		if err := m.parseJSONPayload(r, w, payload); err != nil {
			return
		}
	}

	newCtx := context.WithValue(r.Context(), "payload", payload)
	next.ServeHTTP(w, r.WithContext(newCtx))
}

func (m *PayloadExtractorMiddleware) parseGetPayload(r *http.Request, payload map[string]interface{}) {
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			payload[key] = m.parseValue(values[0])
		}
	}
}

func (m *PayloadExtractorMiddleware) parseRoutePayload(r *http.Request, payload map[string]interface{}) {
	vars := mux.Vars(r)
	for key, value := range vars {
		payload[key] = m.parseValue(value)
	}
}

func (m *PayloadExtractorMiddleware) parseMultiPartPayload(r *http.Request, w http.ResponseWriter, payload map[string]interface{}) error {
	if err := m.checkMaxUploadSize(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			payload[key] = m.parseValue(values[0])
		}
	}

	for key := range r.MultipartForm.File {
		fileUpload, err := m.parseUpload(r, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		payload[key] = fileUpload
	}
	return nil
}

func (m *PayloadExtractorMiddleware) parseJSONPayload(r *http.Request, w http.ResponseWriter, payload map[string]interface{}) error {
	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

func (m *PayloadExtractorMiddleware) parseValue(value string) interface{} {
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

func (m *PayloadExtractorMiddleware) isMultiPart(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

func (m *PayloadExtractorMiddleware) checkMaxUploadSize(r *http.Request) error {
	maxUploadSize, err := strconv.ParseInt(m.maxUploadSize, 10, 64)
	if err != nil || maxUploadSize <= 0 {
		maxUploadSize = 5 << 20
	}

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return err
	}
	return nil
}

func (m *PayloadExtractorMiddleware) parseUpload(r *http.Request, name string) (*webserver_types.FileUpload, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &webserver_types.FileUpload{
		Filename: handler.Filename,
		Size:     handler.Size,
		Content:  fileBytes,
	}, nil
}
