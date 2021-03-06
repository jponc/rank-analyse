package similarity

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"
	"github.com/jponc/rank-analyse/pkg/zenserp"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	zenserpClient          *zenserp.Client
	locations              []string
	country                string
	zenserpBatchWebhookURL string
}

func NewService(
	zenserpClient *zenserp.Client,
	locations []string,
	country string,
	zenserpBatchWebhookURL string,
) *Service {
	s := &Service{
		zenserpClient:          zenserpClient,
		locations:              locations,
		country:                country,
		zenserpBatchWebhookURL: zenserpBatchWebhookURL,
	}

	return s
}

func (s *Service) SimilarityAnalysis(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.zenserpClient == nil {
		log.Errorf("zenserpClient not defined")
		return lambdaresponses.Respond500()
	}

	var req apischema.SimilarityAnalysisRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil || req.Keyword1 == "" || req.Keyword2 == "" {
		log.Errorf("failed to Unmarshal or error keywords")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	keywords := []string{req.Keyword1, req.Keyword2}

	// Create waitgroups
	wg := sync.WaitGroup{}
	count := len(s.locations) * len(keywords)
	wg.Add(count)

	// run queries and add to map
	keyword1Result := map[string]*zenserp.QueryResult{}
	keyword2Result := map[string]*zenserp.QueryResult{}

	for _, location := range s.locations {
		for _, keyword := range keywords {
			go func(l, k string) {
				defer wg.Done()
				res, err := s.zenserpClient.SearchWithLocation(
					ctx,
					k,
					"google.com",
					"desktop",
					s.country,
					l,
					100,
				)

				if err != nil {
					log.Errorf("failed to run zenserp request with location: %v", err)
					return
				}

				if k == req.Keyword1 {
					keyword1Result[l] = res
				} else {
					keyword2Result[l] = res
				}
			}(location, keyword)
		}
	}

	// wait for all to finish
	wg.Wait()

	keyword1SimilarityKeyword := s.buildSimilarityKeyword(req.Keyword1, keyword1Result)
	keyword2SimilarityKeyword := s.buildSimilarityKeyword(req.Keyword2, keyword2Result)

	res := apischema.SimilarityAnalysisResponse{
		Keyword1Similarity: &keyword1SimilarityKeyword,
		Keyword2Similarity: &keyword2SimilarityKeyword,
		Locations:          s.locations,
		Country:            s.country,
	}

	return lambdaresponses.Respond200(res)
}

func (s *Service) SimilarityAnalysisBatch(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.zenserpClient == nil {
		log.Errorf("zenserpClient not defined")
		return lambdaresponses.Respond500()
	}

	if s.zenserpBatchWebhookURL == "" {
		log.Errorf("zenserpBatchWebhookURL not defined")
		return lambdaresponses.Respond500()
	}

	var req apischema.SimilarityAnalysisBatchRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil || req.Keyword1 == "" || req.Keyword2 == "" || req.ClientID == "" {
		log.Errorf("failed to Unmarshal or error keywords")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	jobs := []zenserp.Job{}
	keywords := []string{req.Keyword1, req.Keyword2}

	for _, location := range s.locations {
		for _, keyword := range keywords {

			job := zenserp.Job{
				Query:        keyword,
				Num:          "100",
				SearchEngine: "google.com",
				Device:       "desktop",
				Country:      s.country,
				Location:     location,
			}

			jobs = append(jobs, job)
		}
	}

	batchRes, err := s.zenserpClient.Batch(ctx, req.ClientID, s.zenserpBatchWebhookURL, jobs)
	if err != nil {
		log.Errorf("failed to request zenserp batch: %w", err)
		return lambdaresponses.Respond500()
	}

	res := apischema.SimilarityAnalysisBatchResponse{
		BatchID: batchRes.BatchID,
	}

	return lambdaresponses.Respond200(res)
}

func (s *Service) SimilarityAnalysisBatchStatus(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.zenserpClient == nil {
		log.Errorf("zenserpClient not defined")
		return lambdaresponses.Respond500()
	}

	var req apischema.SimilarityAnalysisBatchStatusRequest

	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil || req.BatchID == "" {
		log.Errorf("failed to Unmarshal or batchID missing: %w", err)
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	b, err := s.zenserpClient.GetBatch(ctx, req.BatchID)
	if err != nil || b.State != "notified" {
		log.Errorf("failed to get batch: %w", err)
		return lambdaresponses.Respond404(fmt.Errorf("failed to get batch: %s", req.BatchID))
	}

	similarityKeywords := s.buildSimilarityKeywordFromBatch(*b)
	if len(similarityKeywords) != 2 {
		log.Errorf("similarityKeywords are not 2: %d", len(similarityKeywords))
		return lambdaresponses.Respond500()
	}

	res := apischema.SimilarityAnalysisBatchStatusResponse{
		Keyword1Similarity: &similarityKeywords[0],
		Keyword2Similarity: &similarityKeywords[1],
	}

	return lambdaresponses.Respond200(res)
}

func (s *Service) buildSimilarityKeywordFromBatch(batch zenserp.Batch) []types.SimilarityKeyword {
	similarityKeywords := []types.SimilarityKeyword{}

	type KeywordItem struct {
		Title     string
		URL       string
		SeenCount int
		Positions []int
	}

	keywordMap := map[string]map[string]KeywordItem{}

	for _, result := range batch.Results {
		keyword := result.Query.Query

		if _, found := keywordMap[keyword]; !found {
			keywordMap[keyword] = map[string]KeywordItem{}
		}

		for _, item := range result.ResulItems {
			if itemRes, found := keywordMap[keyword][item.URL]; found {
				itemRes.Positions = append(itemRes.Positions, item.Position)
				itemRes.SeenCount++

				keywordMap[keyword][item.URL] = itemRes
			} else {
				keywordMap[keyword][item.URL] = KeywordItem{
					Title:     item.Title,
					URL:       item.URL,
					SeenCount: 1,
					Positions: []int{item.Position},
				}
			}
		}
	}

	for k, r := range keywordMap {
		results := []types.SimilarityResult{}

		for _, i := range r {
			if i.Title == "" {
				continue
			}

			// Compute for ave position
			sum := 0
			totalCount := len(i.Positions)
			for _, p := range i.Positions {
				sum = sum + p
			}

			results = append(results, types.SimilarityResult{
				SeenCount:   i.SeenCount,
				Title:       i.Title,
				Link:        i.URL,
				AvePosition: float32(sum) / float32(totalCount),
			})
		}

		sort.SliceStable(results, func(i, j int) bool {
			return results[i].AvePosition < results[j].AvePosition
		})

		similarityKeywords = append(similarityKeywords, types.SimilarityKeyword{
			Keyword: k,
			Results: results,
		})
	}

	return similarityKeywords
}

func (s *Service) buildSimilarityKeyword(keyword string, locationResult map[string]*zenserp.QueryResult) types.SimilarityKeyword {
	// title to similarity result map
	resultsMap := map[string]*types.SimilarityResult{}

	// title to positions array
	positionsMap := map[string][]int{}

	for _, result := range locationResult {
		for _, item := range result.ResulItems {
			if item.Title == "" {
				continue
			}

			if res, found := resultsMap[item.Title]; found {
				positionsMap[item.Title] = append(positionsMap[item.Title], item.Position)
				res.SeenCount++
			} else {
				resultsMap[item.Title] = &types.SimilarityResult{
					SeenCount: 1,
					Title:     item.Title,
					Link:      item.URL,
				}

				positionsMap[item.Title] = []int{item.Position}
			}
		}
	}

	results := []types.SimilarityResult{}
	for _, r := range resultsMap {
		sum := 0

		for _, p := range positionsMap[r.Title] {
			sum = sum + p
		}

		r.AvePosition = float32(sum) / float32(len(positionsMap[r.Title]))
		results = append(results, *r)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].AvePosition < results[j].AvePosition
	})

	similarityKeyword := types.SimilarityKeyword{
		Keyword: keyword,
		Results: results,
	}

	return similarityKeyword
}
