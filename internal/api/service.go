package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	repository *dbrepository.Repository
}

func NewService(repository *dbrepository.Repository) *Service {
	s := &Service{
		repository: repository,
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
	if s.repository == nil {
		return lambdaresponses.Respond500()
	}

	req := &apischema.RunCrawlRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil || req.Keyword == "" {
		log.Errorf("Failed to Unmarshal")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	return lambdaresponses.Respond200(apischema.RunCrawlResponse{Status: "OK"})
}
