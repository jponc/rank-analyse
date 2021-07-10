package main

import (
	webhooks "github.com/jponc/rank-analyse/internal/webhook"
	"github.com/jponc/rank-analyse/pkg/pusher"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	pusherClient, err := pusher.NewClient(config.PusherAppID, config.PusherKey, config.PusherSecret, config.PusherCluster)
	if err != nil {
		log.Fatalf("cannot initialise pusher client %v", err)
	}

	service := webhooks.NewService(pusherClient)
	lambda.Start(service.ZenserpBatchWebhook)
}
