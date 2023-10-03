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

// GetCorrelationID retrieves the correlationID from gls.
func GetCorrelationID() string {
	if correlationID, ok := fcm.GetValue("x-correlation-id"); ok {
		return correlationID.(string)
	}
	return ""
}

func InitializeLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	if env.GetEnv("ENV", "dev") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	UpdateCorrelationID() // Call this to set the initial logger
}

func UpdateCorrelationID() {
	correlationID := GetCorrelationID()

	logger = logrus.WithFields(logrus.Fields{
		"app":              env.GetEnv("APP_NAME"),
		"env":              env.GetEnv("ENV"),
		"version":          env.GetEnv("VERSION"),
		"x-correlation-id": correlationID,
	})
}

func Fatal(args ...interface{}) {
	UpdateCorrelationID()
	logger.Fatal(args)
}

func Debug(args ...interface{}) {
	UpdateCorrelationID()
	logger.Debug(args)
}

func Info(args ...interface{}) {
	UpdateCorrelationID()
	logger.Info(args)
}

func Error(args ...interface{}) {
	UpdateCorrelationID()
	logger.Error(args)
}

func Warning(args ...interface{}) {
	UpdateCorrelationID()
	logger.Warning(args)
}

func Debugf(format string, args ...interface{}) {
	UpdateCorrelationID()
	logger.Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	UpdateCorrelationID()
	logger.Infof(format, args)
}

func Fatalf(format string, args ...interface{}) {
	UpdateCorrelationID()
	logger.Fatalf(format, args)
}

func Errorf(format string, args ...interface{}) {
	UpdateCorrelationID()
	logger.Errorf(format, args)
}

func Warningf(format string, args ...interface{}) {
	UpdateCorrelationID()
	logger.Warningf(format, args)
}
