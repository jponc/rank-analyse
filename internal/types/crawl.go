package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type CrawlArr []Crawl

type Crawl struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Keyword      string    `db:"keyword" json:"keyword"`
	SearchEngine string    `db:"search_engine" json:"search_engine"`
	Device       string    `db:"device" json:"device"`
	Done         bool      `db:"done" json:"done"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
