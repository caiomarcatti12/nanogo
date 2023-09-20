package util

type ResultPaginated struct {
	Rows  []interface{} `json:"rows"`
	Total int64         `json:"total"`
}
