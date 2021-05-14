package api

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"
)

type Service struct {
}

func NewService() *Service {
	s := &Service{}

	return s
}

func (s *Service) Healthcheck(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lambdaresponses.Respond200(apischema.HealthcheckResponse{Status: "OK"})
}
