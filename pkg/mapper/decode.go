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

//func InjectData(input interface{}, output interface{}) error {
//	decoderConfig := &mapstructure.DecoderConfig{
//		DecodeHook: mapstructure.ComposeDecodeHookFunc(
//			decodeTimeFromString,
//			decodeBasicTypesFromString,
//			decodeUUIDFromString,
//		),
//		Result: &output,
//	}
//
//	decoder, err := mapstructure.NewDecoder(decoderConfig)
//
//	if err != nil {
//		return err
//	}
//
//	err = decoder.Decode(input)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//func Transform(input interface{}, output interface{}) error {
//	jsonBytes, err := json.Marshal(input)
//	if err != nil {
//		return err
//	}
//
//	// Usando reflexão para verificar se o input é uma slice
//	inputVal := reflect.ValueOf(input)
//	if inputVal.Kind() == reflect.Slice {
//		// Tratamento para slice
//		var sliceOfMaps []map[string]interface{}
//		err = json.Unmarshal(jsonBytes, &sliceOfMaps)
//		if err != nil {
//			return err
//		}
//		return InjectData(sliceOfMaps, output)
//	} else {
//		// Tratamento para uma única struct
//		var singleMap map[string]interface{}
//		err = json.Unmarshal(jsonBytes, &singleMap)
//		if err != nil {
//			return err
//		}
//
//		return InjectData(singleMap, output)
//	}
//}
//
//// decodeUUIDFromString decodifica UUIDs a partir de strings.
//func decodeUUIDFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
//	if f.Kind() == reflect.String && (t == reflect.TypeOf(uuid.UUID{}) || t == reflect.TypeOf(&uuid.UUID{})) {
//		if data.(string) == "" {
//			if t == reflect.TypeOf(uuid.UUID{}) {
//				return uuid.Nil, nil
//			}
//
//			return nil, nil
//		}
//		uuidVal, err := uuid.Parse(data.(string))
//
//		if err != nil {
//			return uuid.Nil, err
//		}
//
//		if t == reflect.TypeOf(&uuid.UUID{}) {
//			return &uuidVal, nil
//		}
//		return uuidVal, nil
//	}
//	return data, nil
//}
//
//// decodeTimeFromString decodifica time.Time a partir de strings.
//func decodeTimeFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
//	if f.Kind() == reflect.String && t == reflect.TypeOf(time.Time{}) {
//		strData := data.(string)
//		if strData == "" {
//			return time.Time{}, nil // Considerando time.Time zero se a string for vazia.
//		}
//		timeVal, err := time.Parse(time.RFC3339, strData)
//		if err != nil {
//			fmt.Printf("Erro ao decodificar data: %v, dado recebido: %s\n", err, strData)
//			return nil, err
//		}
//		return timeVal, nil
//	}
//	return data, nil
//}
//
//// decodeBasicTypesFromString decodifica tipos básicos como int, float64 e bool a partir de strings.
//func decodeBasicTypesFromString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
//	if f.Kind() == reflect.String {
//		switch t.Kind() {
//		case reflect.Int:
//			if data.(string) == "" {
//				return 0, nil
//			}
//			intVal, err := strconv.Atoi(data.(string))
//			if err != nil {
//				return nil, err
//			}
//			return intVal, nil
//		case reflect.Float64:
//			if data.(string) == "" {
//				return 0.0, nil
//			}
//			floatVal, err := strconv.ParseFloat(data.(string), 64)
//			if err != nil {
//				return nil, err
//			}
//			return floatVal, nil
//		case reflect.Bool:
//			if data.(string) == "" {
//				return false, nil
//			}
//			boolVal, err := strconv.ParseBool(data.(string))
//			if err != nil {
//				return nil, err
//			}
//			return boolVal, nil
//		}
//	}
//
//	return data, nil
//}
