package apischema

import "github.com/jponc/rank-analyse/internal/types"

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type RunCrawlRequest struct {
	Keyword string `json:"keyword"`
	Email   string `json:"email"`
}

type RunCrawlResponse struct {
	Status string `json:"status"`
}

type GetCrawlResponse struct {
	*types.Crawl
}

type GetCrawlsResponse struct {
	Data *[]types.Crawl `json:"data"`
}
