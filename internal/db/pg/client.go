package pg

import (
	"wireguard-api/internal/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const driver = "pgx"

type Client struct {
	masterDBC db.DB
}

func NewClient(dsn string) (db.Client, error) {
	conn, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
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
		return c.masterDBC.Close()
	}

	return nil
}
