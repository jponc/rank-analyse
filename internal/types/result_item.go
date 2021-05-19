package types

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/pkg/zenserp"
)

type ResultItemArray []ResultItem

type ResultItem struct {
	ID            uuid.UUID `json:"id"`
	CrawlResultID uuid.UUID `json:"crawl_result_id"`
	Query         string    `json:"query"`
	SearchEngine  string    `json:"search_engine"`
	Device        string    `json:"device"`
	QueryURL      string    `json:"query_url"`
	Position      int       `json:"position"`
	Title         string    `json:"title"`
	ItemURL       string    `json:"item_url"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

func (t *ResultItemArray) Unmarshal(src interface{}) error {
	switch src.(type) {
	case *zenserp.QueryResult:
		return t.unmarshalZenserpQueryResult(*src.(*zenserp.QueryResult))
	case zenserp.QueryResult:
		return t.unmarshalZenserpQueryResult(src.(zenserp.QueryResult))
	}

	return fmt.Errorf("Failed to unmarshal types.ResultItemArrray: '%v'", src)
}

func (t *ResultItemArray) unmarshalZenserpQueryResult(zenserpQueryResult zenserp.QueryResult) error {
	items := ResultItemArray{}
	crawlResultID := uuid.Must(uuid.NewV4())
	timeNow := time.Now()

	for _, item := range zenserpQueryResult.ResulItems {
		items = append(items, ResultItem{
			ID:            uuid.Must(uuid.NewV4()),
			CrawlResultID: crawlResultID,
			Query:         zenserpQueryResult.Query.Query,
			SearchEngine:  zenserpQueryResult.Query.SearchEngine,
			Device:        zenserpQueryResult.Query.Device,
			QueryURL:      zenserpQueryResult.Query.URL,
			Position:      item.Position,
			Title:         item.Title,
			ItemURL:       item.URL,
			Description:   item.Description,
			CreatedAt:     timeNow,
		})
	}

	*t = items

	return nil
}
