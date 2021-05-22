package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-analyse/internal/api"
	"github.com/jponc/rank-analyse/pkg/sns"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := api.NewService(nil, snsClient, nil)
	lambda.Start(service.RunCrawl)
}
