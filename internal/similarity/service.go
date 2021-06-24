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
	zenserpClient *zenserp.Client
	locations     []string
	country       string
}

func NewService(
	zenserpClient *zenserp.Client,
	locations []string,
	country string,
) *Service {
	s := &Service{
		zenserpClient: zenserpClient,
		locations:     locations,
		country:       country,
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
					20,
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

func (s *Service) buildSimilarityKeyword(keyword string, locationResult map[string]*zenserp.QueryResult) types.SimilarityKeyword {
	// title to similarity result map
	resultsMap := map[string]*types.SimilarityResult{}

	for _, result := range locationResult {
		for _, item := range result.ResulItems {
			if item.Title == "" {
				continue
			}

			if res, found := resultsMap[item.Title]; found {
				res.Positions = append(res.Positions, item.Position)
				res.SeenCount++
			} else {
				resultsMap[item.Title] = &types.SimilarityResult{
					Positions: []int{item.Position},
					SeenCount: 1,
					Title:     item.Title,
				}
			}
		}
	}

	results := []types.SimilarityResult{}
	for _, r := range resultsMap {
		results = append(results, *r)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].SeenCount > results[j].SeenCount
	})

	similarityKeyword := types.SimilarityKeyword{
		Keyword: keyword,
		Results: results,
	}

	return similarityKeyword
}
