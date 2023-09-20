package rsql

type ResultPaginated[T any] struct {
	Rows  *T    `json:"rows"`
	Total int64 `json:"total"`
}
