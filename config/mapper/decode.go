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
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			decodeTimeFromString,       // Adiciona o hook personalizado para uuid
			decodeUUIDFromString,       // Adiciona o hook personalizado para uuid
			decodeBasicTypesFromString, // Adiciona o hook personalizado para uuid

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

// decodeUUIDFromString decodifica UUIDs a partir de strings.
func decodeUUIDFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String && (t == reflect.TypeOf(uuid.UUID{}) || t == reflect.TypeOf(&uuid.UUID{})) {
		if data.(string) == "" {
			return uuid.Nil, nil
		}
		uuidVal, err := uuid.Parse(data.(string))
		if err != nil {
			return nil, err
		}
		if t == reflect.TypeOf(&uuid.UUID{}) {
			return &uuidVal, nil
		}
		return uuidVal, nil
	}
	return data, nil
}

// decodeTimeFromString decodifica time.Time a partir de strings.
func decodeTimeFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String && t == reflect.TypeOf(time.Time{}) {
		if data.(string) == "" {
			return nil, nil
		}
		timeVal, err := time.Parse(time.RFC3339, data.(string))
		if err != nil {
			return nil, err
		}
		return timeVal, nil
	}
	return data, nil
}

// decodeBasicTypesFromString decodifica tipos básicos como int, float64 e bool a partir de strings.
func decodeBasicTypesFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String {
		switch t.Kind() {
		case reflect.Int:
			if data.(string) == "" {
				return 0, nil
			}
			intVal, err := strconv.Atoi(data.(string))
			if err != nil {
				return nil, err
			}
			return intVal, nil
		case reflect.Float64:
			if data.(string) == "" {
				return 0.0, nil
			}
			floatVal, err := strconv.ParseFloat(data.(string), 64)
			if err != nil {
				return nil, err
			}
			return floatVal, nil
		case reflect.Bool:
			if data.(string) == "" {
				return false, nil
			}
			boolVal, err := strconv.ParseBool(data.(string))
			if err != nil {
				return nil, err
			}
			return boolVal, nil
		}
	}
	return data, nil
}
