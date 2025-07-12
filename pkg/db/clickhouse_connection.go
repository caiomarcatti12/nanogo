package db

import (
	"context"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

var clickhouseOnce sync.Once

// ClickhouseCredential holds the connection configuration.
type ClickhouseCredential struct {
	Addr     string
	Username string
	Password string
	Database string
}

// ClickhouseClient implements IDatabase for ClickHouse.
type ClickhouseClient struct {
	conn       clickhouse.Conn
	logger     log.ILog
	credential ClickhouseCredential
}

// NewInstanceClickhouse creates a ClickhouseClient instance.
func NewInstanceClickhouse(cred ClickhouseCredential, logger log.ILog) IDatabase {
	return &ClickhouseClient{
		logger:     logger,
		credential: cred,
	}
}

func (c *ClickhouseClient) Connect() error {
	var err error
	clickhouseOnce.Do(func() {
		c.logger.Trace("Connecting to ClickHouse...")
		c.conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{c.credential.Addr},
			Auth: clickhouse.Auth{
				Database: c.credential.Database,
				Username: c.credential.Username,
				Password: c.credential.Password,
			},
		})
		if err != nil {
			return
		}
		if pingErr := c.conn.Ping(context.Background()); pingErr != nil {
			err = pingErr
			return
		}
		c.logger.Trace("Connection to ClickHouse established!")
	})
	if err != nil {
		c.logger.Error(err.Error())
	}
	return err
}

func (c *ClickhouseClient) GetClient() interface{} { return c.conn }

func (c *ClickhouseClient) Disconnect() {
	if c.conn != nil {
		c.logger.Info("Disconnecting from ClickHouse...")
		_ = c.conn.Close()
	}
}
