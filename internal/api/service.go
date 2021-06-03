package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/repository/s3repository"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"
	"github.com/jponc/rank-analyse/pkg/sns"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	repository   *dbrepository.Repository
	snsClient    *sns.Client
	s3repository *s3repository.Repository
}

func NewService(repository *dbrepository.Repository, snsClient *sns.Client, s3repository *s3repository.Repository) *Service {
	s := &Service{
		repository:   repository,
		snsClient:    snsClient,
		s3repository: s3repository,
	}

	return s
}

func (s *Service) Healthcheck(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lambdaresponses.Respond200(apischema.HealthcheckResponse{Status: "OK"})
}

func (s *Service) RunCrawl(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil {
		log.Errorf("snsClient not defined")
		return lambdaresponses.Respond500()
	}

	req := &apischema.RunCrawlRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil || req.Keyword == "" {
		log.Errorf("failed to Unmarshal or error keyword")
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

func (s *Service) GetCrawls(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s.repository.Connect()
	defer s.repository.Close()

	if s.repository == nil {
		log.Errorf("repository not defined")
		return lambdaresponses.Respond500()
	}

	crawls, err := s.repository.GetCrawls(ctx)
	if err != nil {
		log.Errorf("error getting crawls: %v", err)
		return lambdaresponses.Respond500()
	}

	res := apischema.GetCrawlsResponse{Data: crawls}

	return lambdaresponses.Respond200(res)
}

func (s *Service) GetCrawl(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s.repository.Connect()
	defer s.repository.Close()

	if s.repository == nil {
		log.Errorf("repository not defined")
		return lambdaresponses.Respond500()
	}

	crawlID, err := uuid.FromString(request.PathParameters["id"])
	if err != nil {
		log.Errorf("crawlId missing from path parameters")
		return lambdaresponses.Respond500()
	}

	crawl, err := s.repository.GetCrawl(ctx, crawlID)
	if err != nil {
		log.Errorf("error getting crawl: %v", err)
		return lambdaresponses.Respond500()
	}

	res := apischema.GetCrawlResponse{Data: crawl}

	return lambdaresponses.Respond200(res)
}

func (s *Service) GetResults(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s.repository.Connect()
	defer s.repository.Close()

	if s.repository == nil {
		log.Errorf("repository not defined")
		return lambdaresponses.Respond500()
	}

	crawlID, err := uuid.FromString(request.QueryStringParameters["crawl_id"])
	if err != nil {
		log.Errorf("crawlId missing from path parameters")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	results, err := s.repository.GetCrawlResults(ctx, crawlID)
	if err != nil {
		log.Errorf("error getting crawls: %v", err)
		return lambdaresponses.Respond500()
	}

	res := apischema.GetResultsResponse{Data: results}

	return lambdaresponses.Respond200(res)
}
