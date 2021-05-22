package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type ExtractLink struct {
	ID        uuid.UUID `db:"id"`
	ResultID  uuid.UUID `db:"result_id"`
	Text      string    `db:"text"`
	LinkURL   string    `db:"link_url"`
	CreatedAt time.Time `db:"created_at"`
}
