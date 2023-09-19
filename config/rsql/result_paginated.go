package rsql

type ResultPaginated[T interface{}] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
}
