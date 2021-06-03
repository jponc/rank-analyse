package apischema

import "github.com/jponc/rank-analyse/internal/types"

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type RunCrawlRequest struct {
	Keyword string `json:"keyword"`
}

type RunCrawlResponse struct {
	Status string `json:"status"`
}

type GetCrawlResponse struct {
	Data *types.Crawl `json:"data"`
}

type GetResultResponse struct {
	Data *types.Result `json:"data"`
}

type GetResultInfoResponse struct {
	Data *types.ExtractInfo `json:"data"`
}

type GetResultLinksResponse struct {
	Data *[]types.ExtractLink `json:"data"`
}

type GetCrawlsResponse struct {
	Data *[]types.Crawl `json:"data"`
}

type GetResultsResponse struct {
	Data *[]types.Result `json:"data"`
}
