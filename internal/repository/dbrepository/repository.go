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

func (r *Repository) Connect() error {
	return r.dbClient.Connect()
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

func (r *Repository) GetCrawls(ctx context.Context) (*[]types.Crawl, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	crawls := []types.Crawl{}

	err := r.dbClient.SelectContext(
		ctx,
		&crawls,
		`SELECT * FROM crawl ORDER BY created_at DESC`,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get crawl: %w", err)
	}

	return &crawls, nil
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

	if len(links) == 0 {
		return nil
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
		return fmt.Errorf("failed to insert extract links: %v", err)
	}

	return nil
}

func (r *Repository) CreateAnalyzeTopics(ctx context.Context, resultID uuid.UUID, topics types.AnalyzeTopicArray) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	if len(topics) == 0 {
		return nil
	}

	queryInsert := `INSERT INTO analyze_topics (result_id, label, score) VALUES `
	insertParams := []interface{}{}

	for i, topic := range topics {
		p := i * 3 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d,$%d,$%d),", p+1, p+2, p+3)
		insertParams = append(insertParams, resultID, topic.Label, topic.Score)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","

	_, err := r.dbClient.Exec(
		ctx,
		queryInsert,
		insertParams...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert analyze topics: %v", err)
	}

	return nil
}

func (r *Repository) CreateAnalyzeEntities(ctx context.Context, resultID uuid.UUID, entities types.AnalyzeEntityArray) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	if len(entities) == 0 {
		return nil
	}

	queryInsert := `INSERT INTO analyze_entities (result_id, confidence_score, relevance_score, entity, matched_text) VALUES `
	insertParams := []interface{}{}

	for i, entity := range entities {
		p := i * 5 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", p+1, p+2, p+3, p+4, p+5)
		insertParams = append(insertParams, resultID, entity.ConfidenceScore, entity.RelevanceScore, entity.Entity, entity.MatchedText)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","

	_, err := r.dbClient.Exec(
		ctx,
		queryInsert,
		insertParams...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert analyze topics: %v", err)
	}

	return nil
}

func (r *Repository) MarkResultAsDone(ctx context.Context, resultID uuid.UUID, isError bool) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.Exec(
		ctx,
		`UPDATE result SET done = true, is_error = $1 WHERE id = $2`,
		isError,
		resultID,
	)

	if err != nil {
		return fmt.Errorf("failed to update result done: %v", err)
	}

	return nil
}

func (r *Repository) MarkCrawlAsDone(ctx context.Context, crawlID uuid.UUID) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.Exec(
		ctx,
		`UPDATE crawl SET done = true WHERE id = $1`,
		crawlID,
	)

	if err != nil {
		return fmt.Errorf("failed to update crawl done: %v", err)
	}

	return nil
}

func (r *Repository) IsAllCrawlResultsDone(ctx context.Context, crawlID uuid.UUID) (bool, error) {
	if r.dbClient == nil {
		return false, fmt.Errorf("dbClient not initialised")
	}

	var isDone bool

	err := r.dbClient.GetContext(
		ctx,
		&isDone,
		`
			SELECT COUNT(id) = 0
			FROM result
			WHERE crawl_id = $1 AND done = false
		`,
		crawlID,
	)

	if err != nil {
		return false, fmt.Errorf("failed to get not done crawl results: %v", err)
	}

	return isDone, nil
}

func (r *Repository) GetCrawlResults(ctx context.Context, crawlID uuid.UUID) (*[]types.Result, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var results []types.Result

	err := r.dbClient.SelectContext(
		ctx,
		&results,
		`
			SELECT *
			FROM result
			WHERE crawl_id = $1 AND done = true AND is_error = false
			ORDER BY position ASC
		`,
		crawlID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get crawl results: %v", err)
	}

	return &results, nil
}

func (r *Repository) GetExtractInfo(ctx context.Context, resultID uuid.UUID) (*types.ExtractInfo, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var extractInfo types.ExtractInfo

	err := r.dbClient.GetContext(
		ctx,
		&extractInfo,
		`
			SELECT *
			FROM extract_info
			WHERE result_id = $1
		`,
		resultID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get extract info: %v", err)
	}

	return &extractInfo, nil
}

func (r *Repository) GetExtractLinks(ctx context.Context, resultID uuid.UUID) (*[]types.ExtractLink, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var extractLinks []types.ExtractLink

	err := r.dbClient.SelectContext(
		ctx,
		&extractLinks,
		`
			SELECT *
			FROM extract_links
			WHERE result_id = $1
		`,
		resultID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to select extract links: %v", err)
	}

	return &extractLinks, nil
}

func (r *Repository) GetTopics(ctx context.Context, resultID uuid.UUID) (*[]types.AnalyzeTopic, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var topics []types.AnalyzeTopic

	err := r.dbClient.SelectContext(
		ctx,
		&topics,
		`
			SELECT *
			FROM analyze_topics
			WHERE result_id = $1
		`,
		resultID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to select topics: %v", err)
	}

	if topics == nil {
		return &[]types.AnalyzeTopic{}, nil
	}

	return &topics, nil
}

func (r *Repository) GetEntities(ctx context.Context, resultID uuid.UUID) (*[]types.AnalyzeEntity, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var entities []types.AnalyzeEntity

	err := r.dbClient.SelectContext(
		ctx,
		&entities,
		`
			SELECT *
			FROM analyze_entities
			WHERE result_id = $1
		`,
		resultID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to select entities: %v", err)
	}

	if entities == nil {
		return &[]types.AnalyzeEntity{}, nil
	}

	return &entities, nil
}

func (r *Repository) SaveCleanedText(ctx context.Context, resultID uuid.UUID, cleanedText string) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.Exec(
		ctx,
		`UPDATE extract_info SET cleaned_text = $1 WHERE result_id = $2`,
		cleanedText,
		resultID,
	)

	if err != nil {
		return fmt.Errorf("failed to update extract info cleaned text: %v", err)
	}

	return nil
}

func (r *Repository) Close() error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	return r.dbClient.Close()
}
