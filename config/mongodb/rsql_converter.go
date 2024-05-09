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
package mongodb

import (
	"strings"

	"github.com/caiomarcatti12/nanogo/v2/config/rsql"
	"go.mongodb.org/mongo-driver/bson"
)

func RsqlConvertToBson(qp rsql.QueryFilter) (bson.M, bson.M, int64, int64, error) {
	query := bson.M{}

	// Verificar se Filter foi fornecido
	if qp.Filter != "" {
		// Convertendo o filtro RSQL para bson.M
		conditions, err := rsql.Parse(qp.Filter)
		if err != nil {
			return nil, nil, 0, 0, err
		}

		for _, condition := range conditions {
			switch condition.Operator {
			case "==":
				query[condition.Field] = condition.Value
			case "!=":
				query[condition.Field] = bson.M{"$ne": condition.Value}
			case "<>":
				query[condition.Field] = bson.M{"$ne": condition.Value}
			case "=gt=":
				query[condition.Field] = bson.M{"$gt": condition.Value}
			case "=gte=":
				query[condition.Field] = bson.M{"$gte": condition.Value}
			case "=lt=":
				query[condition.Field] = bson.M{"$lt": condition.Value}
			case "=lte=":
				query[condition.Field] = bson.M{"$lte": condition.Value}
			case "=like=":
				query[condition.Field] = bson.M{"$regex": condition.Value, "$options": "i"}
			}
		}
	}

	sort := bson.M{}
	// Verificar se SortParams foi fornecido
	if qp.SortParams != "" {
		// Convertendo o parâmetro de ordenação para bson.M
		sortFields := strings.Split(qp.SortParams, ":")
		if len(sortFields) == 2 {
			field := sortFields[0]
			direction := 0
			switch sortFields[1] {
			case "ASC":
				direction = 1
			case "DESC":
				direction = -1
			}
			if direction != 0 {
				sort[field] = direction
			}
		}
	}

	size := 15
	// Verificar se Size foi fornecido
	if qp.Size != 15 {
		size = int(qp.Size)
	}

	skip := 0
	// Verificar se Skip foi fornecido
	if qp.Skip != 0 {
		skip = int(qp.Skip)
	}

	return query, sort, int64(size), int64(skip), nil
}
