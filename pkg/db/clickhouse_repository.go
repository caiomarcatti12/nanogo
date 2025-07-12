package db

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

// IClickhouseRepository defines basic query operations for ClickHouse
// using the official driver.
type IClickhouseRepository[T any] interface {
	// Exec executes queries that don't return rows, e.g. INSERT/UPDATE.
	Exec(query string, args ...any) error
	// Select fills dest with the result of the query.
	Select(dest *[]T, query string, args ...any) error
}

// ClickhouseRepository is a simple repository around a clickhouse.Conn.
type ClickhouseRepository[T any] struct {
	conn   clickhouse.Conn
	logger log.ILog
}

// NewClickhouseRepository creates a ClickhouseRepository for the
// given IDatabase. The database must have been created with the
// clickhouse provider.
func NewClickhouseRepository[T any](db IDatabase, logger log.ILog) IClickhouseRepository[T] {
	if conn, ok := db.GetClient().(clickhouse.Conn); ok {
		return &ClickhouseRepository[T]{conn: conn, logger: logger}
	}
	panic("error constructing ClickhouseRepository: client is not clickhouse.Conn")
}

func (r *ClickhouseRepository[T]) Exec(query string, args ...any) error {
	r.logger.Trace("clickhouse exec: " + query)
	return r.conn.Exec(context.Background(), query, args...)
}

func (r *ClickhouseRepository[T]) Select(dest *[]T, query string, args ...any) error {
	r.logger.Trace("clickhouse select: " + query)
	if err := r.conn.Select(context.Background(), dest, query, args...); err != nil {
		r.logger.Error(err.Error())
		return err
	}
	if len(*dest) == 0 {
		return nil
	}
	return nil
}
