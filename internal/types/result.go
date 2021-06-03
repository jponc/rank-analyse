package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type Result struct {
	ID          uuid.UUID `db:"id" json:"id"`
	CrawlID     uuid.UUID `db:"crawl_id" json:"crawl_id"`
	Link        string    `db:"link" json:"link"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    int       `db:"position" json:"position"`
	Done        bool      `db:"done" json:"done"`
	IsError     bool      `db:"is_error" json:"is_error"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
