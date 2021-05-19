package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	db     *sqlx.DB
	TestDB *sqlx.DB
}

func NewClient(connectionURL string) (*Client, error) {
	db, err := sqlx.Connect("postgres", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	p := &Client{
		db:     db,
		TestDB: db,
	}

	return p, nil
}

func (c *Client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *Client) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}
