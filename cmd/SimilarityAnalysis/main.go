package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-analyse/internal/api"
	pkgHttp "github.com/jponc/rank-analyse/pkg/http"
	"github.com/jponc/rank-analyse/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	httpClient := pkgHttp.DefaultHTTPClient(time.Duration(1 * time.Minute))
	zenserpClient, err := zenserp.NewClient(config.ZenserpApiKey, httpClient)
	if err != nil {
		log.Fatalf("cannot initialise zenserp client %v", err)
	}

	service := api.NewService(nil, nil, nil, zenserpClient)

	lambda.Start(service.SimilarityAnalysis)
}
