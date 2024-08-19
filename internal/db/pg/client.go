package pg

import (
	"fmt"
	"wireguard-api/internal/db"

	_ "github.com/jackc/pgx/v5/stdlib" // enable db driver.
	"github.com/jmoiron/sqlx"
)

const driver = "pgx"

type Client struct {
	masterDBC db.DB
}

func NewClient(dsn string) (db.Client, error) {
	conn, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("pg_client.new_client %w", err)
	}

	return &Client{
		masterDBC: NewDB(conn),
	}, nil
}

func (c *Client) DB() db.DB {
	return c.masterDBC
}

func (c *Client) Close() error {
	if c.masterDBC != nil {
		err := c.masterDBC.Close()
		if err != nil {
			return fmt.Errorf("pg_client.close %w", err)
		}

		return nil
	}

	return nil
}
