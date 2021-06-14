package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type ExtractInfo struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ResultID    uuid.UUID `db:"result_id" json:"result_id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	CleanedText string    `db:"cleaned_text" json:"cleaned_text"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
