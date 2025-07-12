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
	"fmt"
	"os"
	"strings"

	"github.com/bytedance/sonic"

	"github.com/caiomarcatti12/nanogo/pkg/context_manager"
	"github.com/google/uuid"

	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Logger struct {
	logger        *logrus.Entry
	fcm           context_manager.ISafeContextManager
	defaultFields logrus.Fields
}

func newLogger(env env.IEnv, contextManager context_manager.ISafeContextManager) ILog {
	l := &Logger{
		fcm: contextManager,
	}

	logrus.SetOutput(os.Stdout)

	// Set log level from environment variable
	logLevel := env.GetEnv("LOG_LEVEL", "DEBUG") // Default to "trace" if not specified
	l.setLogLevel(logLevel)

	if env.GetEnv("ENV", "dev") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	l.defaultFields = logrus.Fields{
		"app":              env.GetEnv("APP_NAME"),
		"env":              env.GetEnv("ENV"),
		"version":          env.GetEnv("VERSION"),
		"x-correlation-id": "",
	}

	l.updateCorrelationID()

	return l
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Fatal(message)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Debug(message)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Info(message)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Error(message)
}

func (l *Logger) Warning(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Warning(message)
}
func (l *Logger) Trace(message string, args ...interface{}) {
	l.updateCorrelationID()
	fields := l.extractFields(args)
	l.logger.WithFields(fields).Trace(message)
}

func (l *Logger) setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.TraceLevel) // Default level if none is provided
	}
}

func (l *Logger) updateCorrelationID() {
	correlationID, _ := l.fcm.GetValue("x-correlation-id")

	if correlationID == "" {
		correlationID = uuid.New().String()
	}
	fieldsLogger := l.defaultFields

	fieldsLogger["x-correlation-id"] = correlationID

	l.logger = logrus.WithFields(fieldsLogger)
}

func (l *Logger) extractFields(args ...interface{}) logrus.Fields {
	if len(args) == 0 {
		return nil
	}

	innerArgs, ok := args[0].([]interface{})
	if !ok || len(innerArgs) == 0 {
		return nil
	}

	fields := logrus.Fields{}
	if len(innerArgs) > 1 {
		switch v := innerArgs[1].(type) {
		case Fields:
			fields = logrus.Fields(v)
		default:
			data, err := sonic.Marshal(v)
			if err == nil {
				_ = sonic.Unmarshal(data, &fields)
			}
		}
	}

	return fields
}

func (l *Logger) Fatalf(message string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(message, args...))
}

func (l *Logger) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

func (l *Logger) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

func (l *Logger) Errorf(message string, args ...interface{}) {
	l.Error(fmt.Sprintf(message, args...))
}

func (l *Logger) Warningf(message string, args ...interface{}) {
	l.Warning(fmt.Sprintf(message, args...))
}

func (l *Logger) Tracef(message string, args ...interface{}) {
	l.Trace(fmt.Sprintf(message, args...))
}
