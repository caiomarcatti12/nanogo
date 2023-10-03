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
package log

import (
	"encoding/json"
	"github.com/caiomarcatti12/nanogo/v2/config/context_manager"
	"os"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/sirupsen/logrus"
)

var (
	logger         *logrus.Entry
	logInitialized bool
	fcm            = context_manager.NewSafeContextManager()
)

type Fields map[string]interface{}

// GetCorrelationID retrieves the correlationID from gls.
func GetCorrelationID() string {
	if correlationID, ok := fcm.GetValue("x-correlation-id"); ok {
		return correlationID.(string)
	}
	return ""
}

func InitializeLogger() {
	if logInitialized {
		return
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	if env.GetEnv("ENV", "dev") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logInitialized = true

	UpdateCorrelationID() // Call this to set the initial logger
}

func UpdateCorrelationID() {
	if logger == nil {
		InitializeLogger()
	}

	correlationID := GetCorrelationID()

	logger = logrus.WithFields(logrus.Fields{
		"app":              env.GetEnv("APP_NAME"),
		"env":              env.GetEnv("ENV"),
		"version":          env.GetEnv("VERSION"),
		"x-correlation-id": correlationID,
	})
}

func extractFields(args ...interface{}) (string, logrus.Fields) {
	if len(args) == 0 {
		return "", nil
	}

	// Obtendo o primeiro nível
	innerArgs, ok := args[0].([]interface{})
	if !ok || len(innerArgs) == 0 {
		return "", nil
	}

	// Extraindo a mensagem do primeiro nível
	msg, _ := innerArgs[0].(string)

	// Extraindo os campos, se existirem
	fields := logrus.Fields{}
	if len(innerArgs) > 1 {
		switch v := innerArgs[1].(type) {
		case Fields:
			fields = logrus.Fields(v)
		default: // Tratar como struct ou qualquer outro tipo
			data, err := json.Marshal(v)
			if err == nil {
				_ = json.Unmarshal(data, &fields)
			}
		}
	}

	return msg, fields
}

func Fatal(args ...interface{}) {
	UpdateCorrelationID()
	msg, fields := extractFields(args)
	logger.WithFields(fields).Fatal(msg)
}

func Debug(args ...interface{}) {
	UpdateCorrelationID()
	msg, fields := extractFields(args)
	logger.WithFields(fields).Debug(msg)
}

func Info(args ...interface{}) {
	UpdateCorrelationID()
	msg, fields := extractFields(args)
	logger.WithFields(fields).Info(msg)
}

func Error(args ...interface{}) {
	UpdateCorrelationID()
	msg, fields := extractFields(args)
	logger.WithFields(fields).Error(msg)
}

func Warning(args ...interface{}) {
	UpdateCorrelationID()
	msg, fields := extractFields(args)
	logger.WithFields(fields).Warning(msg)
}
