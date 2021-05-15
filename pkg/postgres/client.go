package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(connectionURL string) (*Client, error) {
	db, err := sqlx.Connect("postgres", connectionURL)
	log.Infof(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	p := &Client{
		db,
	}

	return p, nil
}
