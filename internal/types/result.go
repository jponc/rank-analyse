package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type Result struct {
	ID          uuid.UUID `db:"id"`
	CrawlID     uuid.UUID `db:"crawl_id"`
	Link        string    `db:"link"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Position    int       `db:"position"`
	Done        bool      `db:"done"`
	IsError     bool      `db:"is_error"`
	CreatedAt   time.Time `db:"created_at"`
}
