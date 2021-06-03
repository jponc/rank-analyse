package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type ExtractLink struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ResultID  uuid.UUID `db:"result_id" json:"result_id"`
	Text      string    `db:"text" json:"text"`
	LinkURL   string    `db:"link_url" json:"link_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
