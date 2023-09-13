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
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func Transform(input interface{}, output any) error {
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
