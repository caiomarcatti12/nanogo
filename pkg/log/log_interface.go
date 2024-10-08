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

type ILog interface {
	Fatal(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Trace(message string, args ...interface{})

	
	Fatalf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Tracef(message string, args ...interface{})
}
