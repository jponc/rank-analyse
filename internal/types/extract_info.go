package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type ExtractInfo struct {
	ID        uuid.UUID `db:"id"`
	ResultID  uuid.UUID `db:"result_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}