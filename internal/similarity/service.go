package similarity

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/apischema"
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
	locations string,
	country string,
) *Service {
	s := &Service{
		zenserpClient: zenserpClient,
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
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(s.locations) * len(keywords))

	// run queries and add to map
	keyword1Output := map[string]*zenserp.QueryResult{}
	keyword2Output := map[string]*zenserp.QueryResult{}

	for _, location := range s.locations {
		for _, keyword := range keywords {
			go func() {
				defer wg.Done()
				res, err := s.zenserpClient.SearchWithLocation(
					ctx,
					keyword,
					"google.com",
					"desktop",
					s.country,
					location,
					100,
				)

				if err != nil {
					log.Errorf("failed to run zenserp request with location: %v", err)
				}

				key := locationResultKey{
					country:  s.country,
					location: location,
					keyword:  keyword,
				}

				m.Lock()
				output[key] = res
				m.Unlock()
			}()
		}
	}

	wg.Wait()

	return lambdaresponses.Respond200(res)
}
