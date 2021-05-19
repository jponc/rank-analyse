package dbrepository

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/postgres"
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
	id := uuid.Must(uuid.NewV4())

	_, err := r.dbClient.Exec(
		ctx,
		`INSERT INTO crawl (id, keyword, search_engine, device) VALUES ($1, $2, $3, $4)`,
		id, keyword, searchEngine, device,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert: %v", err)
	}

	return r.GetCrawl(ctx, id)
}

func (r *Repository) CreateResult(ctx context.Context, crawlID uuid.UUID, link, title, description string, position int) (*types.Result, error) {
	id := uuid.Must(uuid.NewV4())

	_, err := r.dbClient.Exec(
		ctx,
		`
			INSERT INTO result (id, crawl_id, link, position, title, description)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
		id, crawlID, link, position, title, description)

	if err != nil {
		return nil, fmt.Errorf("failed to insert Result: %v", err)
	}

	return r.GetResult(ctx, id)
}

func (r *Repository) GetResult(ctx context.Context, id uuid.UUID) (*types.Result, error) {
	result := types.Result{}
	err := r.dbClient.GetContext(ctx, &result, "SELECT * FROM result WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	return &result, nil
}

func (r *Repository) GetCrawl(ctx context.Context, id uuid.UUID) (*types.Crawl, error) {
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
