package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"
	"github.com/jponc/rank-analyse/pkg/sns"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	repository *dbrepository.Repository
	snsClient  *sns.Client
}

func NewService(repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		repository: repository,
		snsClient:  snsClient,
	}

	return s
}

func (s *Service) Healthcheck(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lambdaresponses.Respond200(apischema.HealthcheckResponse{Status: "OK"})
}

func (s *Service) LambdaTest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return lambdaresponses.Respond500()
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lambdaresponses.Respond500()
	}
	//Convert the body to type string
	sb := string(body)

	return lambdaresponses.Respond200(apischema.LambdaTestResponse{Out: sb})
}

func (s *Service) RunCrawl(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil {
		log.Errorf("snsClient not defined")
		return lambdaresponses.Respond500()
	}

	req := &apischema.RunCrawlRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil || req.Keyword == "" {
		log.Errorf("failed to Unmarshal")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	msg := eventschema.ProcessKeywordMessage{
		Keyword:      req.Keyword,
		Device:       "desktop",
		SearchEngine: "google.com",
		Count:        100,
	}

	err = s.snsClient.Publish(ctx, eventschema.ProcessKeyword, msg)
	if err != nil {
		log.Errorf("failed to publish SNS")
		return lambdaresponses.Respond500()
	}

	log.Infof("successfully queued keyword %s for processing", msg.Keyword)

	return lambdaresponses.Respond200(apischema.RunCrawlResponse{Status: "OK"})
}

func (s *Service) GetCrawl(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.repository == nil {
		log.Errorf("repository not defined")
		return lambdaresponses.Respond500()
	}

	crawlId, err := uuid.FromString(request.PathParameters["id"])
	if err != nil {
		log.Errorf("crawlId missing from path parameters")
		return lambdaresponses.Respond500()
	}

	crawl, err := s.repository.GetCrawl(ctx, crawlId)
	if err != nil {
		log.Errorf("error when fetching crawl: %v", err)
		return lambdaresponses.Respond500()
	}

	return lambdaresponses.Respond200(apischema.GetCrawlResponse{Crawl: crawl})
}
