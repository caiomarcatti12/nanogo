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
package mapper

import (
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func Transform(input interface{}, output interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				if f.Kind() == reflect.String && t == reflect.TypeOf(&uuid.UUID{}) {
					uuidVal, err := uuid.Parse(data.(string))
					if err != nil {
						return nil, err
					}
					return &uuidVal, nil
				}
				if f.Kind() == reflect.String && t == reflect.TypeOf(uuid.UUID{}) {
					uuidVal, err := uuid.Parse(data.(string))
					if err != nil {
						return nil, err
					}
					return &uuidVal, nil
				}

				if f.Kind() == reflect.String && t == reflect.TypeOf(time.Time{}) {
					timeVal, err := time.Parse(time.RFC3339, data.(string))
					if err != nil {
						return nil, err
					}
					return timeVal, nil
				}

				if f.Kind() == reflect.String && t.Kind() == reflect.Int {
					intVal, err := strconv.Atoi(data.(string))
					if err != nil {
						return nil, err
					}
					return intVal, nil
				}

				if f.Kind() == reflect.String && t.Kind() == reflect.Float64 {
					floatVal, err := strconv.ParseFloat(data.(string), 64)
					if err != nil {
						return nil, err
					}
					return floatVal, nil
				}

				if f.Kind() == reflect.String && t.Kind() == reflect.Bool {
					boolVal, err := strconv.ParseBool(data.(string))
					if err != nil {
						return nil, err
					}
					return boolVal, nil
				}

				return data, nil
			},
		),
		Result: &output,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)

	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	if err != nil {
		return err
	}
	return nil
}
