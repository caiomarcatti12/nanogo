package rsql

import "github.com/caiomarcatti12/nanogo/v2/config/repository"

type ResultPaginated[T repository.Model] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
}
