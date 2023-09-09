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
	"os"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/sirupsen/logrus"
)

var (
	logger         *logrus.Entry
	logInitialized bool
)

func LoadLog(correlationID ...string) *logrus.Entry {
	var cid string
	if len(correlationID) > 0 {
		cid = correlationID[0]
	} else {
		cid = ""
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	if env.GetEnv("ENV", "dev") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logger = logrus.WithFields(logrus.Fields{
		"app":              env.GetEnv("APP_NAME"),
		"env":              env.GetEnv("ENV"),
		"version":          env.GetEnv("VERSION"),
		"x-correlation-id": cid,
	})

	logInitialized = true

	return logger
}

func Fatal(args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Fatal(args)
}

func Debug(args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Debug(args)
}

func Info(args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Info(args)
}

func Error(args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Error(args)
}

func Warning(args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Warning(args)
}

func Debugf(format string, args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Infof(format, args)
}

func Fatalf(format string, args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Fatalf(format, args)
}

func Errorf(format string, args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Errorf(format, args)
}

func Warningf(format string, args ...interface{}) {
	if !logInitialized {
		LoadLog()
	}

	logger.Warningf(format, args)
}
