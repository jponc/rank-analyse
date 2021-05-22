package dbrepository

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/postgres"
	"github.com/jponc/rank-analyse/pkg/webscraper"
)

type Repository struct {
	dbClient *postgres.Client
}

func NewRepository(dbClient *postgres.Client) (*Repository, error) {
	if dbClient == nil {
		return nil, fmt.Errorf("failed to initialise repository: dbClient is nil")
	}

	r := &Repository{
		dbClient,
	}

	return r, nil
}

func (r *Repository) CreateCrawl(ctx context.Context, keyword, searchEngine, device string) (*types.Crawl, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var id uuid.UUID

	err := r.dbClient.GetContext(
		ctx,
		&id,
		`
			INSERT INTO crawl (keyword, search_engine, device)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
		keyword, searchEngine, device,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert: %v", err)
	}

	return r.GetCrawl(ctx, id)
}

func (r *Repository) CreateResult(ctx context.Context, crawlID uuid.UUID, link, title, description string, position int) (*types.Result, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var id uuid.UUID

	err := r.dbClient.GetContext(
		ctx,
		&id,
		`
			INSERT INTO result (crawl_id, link, position, title, description)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`,
		crawlID, link, position, title, description)

	if err != nil {
		return nil, fmt.Errorf("failed to insert Result: %v", err)
	}

	return r.GetResult(ctx, id)
}

func (r *Repository) GetResult(ctx context.Context, id uuid.UUID) (*types.Result, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	result := types.Result{}
	err := r.dbClient.GetContext(
		ctx,
		&result,
		`SELECT * FROM result WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	return &result, nil
}

func (r *Repository) GetCrawl(ctx context.Context, id uuid.UUID) (*types.Crawl, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	crawl := types.Crawl{}

	err := r.dbClient.GetContext(
		ctx,
		&crawl,
		`SELECT * FROM crawl WHERE id=$1`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get crawl: %w", err)
	}

	return &crawl, nil
}

func (r *Repository) CreateExtractInfo(ctx context.Context, resultID uuid.UUID, title, content string) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.Exec(
		ctx,
		`INSERT INTO extract_info (result_id, title, content) VALUES ($1, $2, $3)`,
		resultID, title, content,
	)

	if err != nil {
		return fmt.Errorf("failed to insert extract info: %v", err)
	}

	return nil
}

func (r *Repository) CreateExtractLinks(ctx context.Context, resultID uuid.UUID, links []webscraper.Link) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	// Batch insert links to external_links table

	queryInsert := `INSERT INTO extract_links (result_id, text, link_url) VALUES `
	insertParams := []interface{}{}

	for i, link := range links {
		p := i * 3 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d,$%d,$%d),", p+1, p+2, p+3)
		insertParams = append(insertParams, resultID, link.Text, link.LinkURL)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","

	_, err := r.dbClient.Exec(
		ctx,
		queryInsert,
		insertParams...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert extract info: %v", err)
	}

	return nil
}
