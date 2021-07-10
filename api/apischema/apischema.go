package apischema

import (
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/zenserp"
)

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

type GetResultTopicsResponse struct {
	Data *[]types.AnalyzeTopic `json:"data"`
}

type GetResultEntiitesResponse struct {
	Data *[]types.AnalyzeEntity `json:"data"`
}

type GetCrawlsResponse struct {
	Data *[]types.Crawl `json:"data"`
}

type GetResultsResponse struct {
	Data *[]types.Result `json:"data"`
}

type SimilarityAnalysisRequest struct {
	Keyword1 string `json:"keyword1"`
	Keyword2 string `json:"keyword2"`
}

type SimilarityAnalysisResponse struct {
	Keyword1Similarity *types.SimilarityKeyword `json:"keyword1_similarity"`
	Keyword2Similarity *types.SimilarityKeyword `json:"keyword2_similarity"`
	Locations          []string                 `json:"locations"`
	Country            string                   `json:"country"`
}

type SimilarityAnalysisBatchRequest struct {
	Keyword1 string `json:"keyword1"`
	Keyword2 string `json:"keyword2"`
	ClientID string `json:"client_id"`
}

type SimilarityAnalysisBatchResponse struct {
	BatchID string `json:"batch_id"`
}

type ZenserpBatchWebhookRequest []zenserp.QueryResult
type ZenserpBatchWebhookResponse struct {
	Message string `json:"message"`
}

type SimilarityAnalysisBatchStatusRequest struct {
	BatchID string `json:"batch_id"`
}

type SimilarityAnalysisBatchStatusResponse struct {
	Keyword1Similarity *types.SimilarityKeyword `json:"keyword1_similarity"`
	Keyword2Similarity *types.SimilarityKeyword `json:"keyword2_similarity"`
}
