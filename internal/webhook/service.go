package webhooks

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/apischema"
	"github.com/jponc/rank-analyse/pkg/lambdaresponses"
	"github.com/jponc/rank-analyse/pkg/pusher"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	pusherClient *pusher.Client
}

// NewService instantiates a new service
func NewService(pusherClient *pusher.Client) *Service {
	return &Service{
		pusherClient: pusherClient,
	}
}

func (s *Service) ZenserpBatchWebhook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.pusherClient == nil {
		log.Errorf("pusherClient not defined")
		return lambdaresponses.Respond500()
	}

	log.Info("%v", request.Body)

	res := apischema.ZenserpBatchWebhookResponse{Message: "OK"}

	return lambdaresponses.Respond200(res)
}
