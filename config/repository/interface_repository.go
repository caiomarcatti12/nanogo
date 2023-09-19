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
package repository

import (
	"github.com/caiomarcatti12/nanogo/v2/config/rsql"
	"github.com/google/uuid"
)

type Repository[T Model] interface {
	Insert(document *T) (*T, error)
	Update(document *T) (*T, error)
	Delete(document *T) (bool, error)
	FindById(id uuid.UUID) (*T, error)
	DeleteById(id uuid.UUID) (bool, error)
	FindAll() ([]*T, error)
	RawQueryRsqlFiltered(filter rsql.QueryFilter) (*rsql.ResultPaginated[T], error)
}
