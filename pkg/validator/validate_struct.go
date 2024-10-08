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
package validator

import (
	"fmt"
	"sync"

	"github.com/caiomarcatti12/nanogo/pkg/errors"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

var instance *Validator
var once sync.Once

func getValidatorInstance() *Validator {
	once.Do(func() {
		instance = &Validator{
			validator: validator.New(),
		}
	})
	return instance
}

func ValidateStruct(s interface{}) *errors.CustomError {
	getValidatorInstance()

	err := instance.validator.Struct(s)
	if err != nil {
		// Verifica se o erro é do tipo ValidationErrors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				field := err.Field()
				tag := err.Tag()

				switch tag {
				case "required":
					return InvalidStructException(fmt.Sprintf("field %s is required", field))
				case "email":
					return InvalidStructException(fmt.Sprintf("field %s is not a valid email", field))
				case "gte":
					return InvalidStructException(fmt.Sprintf("field %s should be greater than or equal to %s", field, err.Param()))
				case "lte":
					return InvalidStructException(fmt.Sprintf("field %s should be less than or equal to %s", field, err.Param()))
				default:
					return InvalidStructException(fmt.Sprintf("field %s has invalid value", field))
				}
			}
		} /* else if invalidValidationError, ok := err.(*validator.InvalidValidationError); ok {
			return InvalidStructException(invalidValidationError.Error())
		} else {
			// Lidar com outros tipos de erros aqui
		}*/

	}
	return nil
}
