package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/resultrankings"
	pkgHttp "github.com/jponc/rank-analyse/pkg/http"
	"github.com/jponc/rank-analyse/pkg/postgres"
	"github.com/jponc/rank-analyse/pkg/sns"
	"github.com/jponc/rank-analyse/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	pgClient, err := postgres.NewClient(config.RDSConnectionURL)
	if err != nil {
		log.Fatalf("cannot initialise pg client: %v", err)
	}

	dbRepository, err := dbrepository.NewRepository(pgClient)
	if err != nil {
		log.Fatalf("cannot initialise repository: %v", err)
	}

	httpClient := pkgHttp.DefaultHTTPClient(time.Duration(1 * time.Minute))
	zenserpClient, err := zenserp.NewClient(config.ZenserpApiKey, httpClient)
	if err != nil {
		log.Fatalf("cannot initialise zenserp client %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := resultrankings.NewService(zenserpClient, dbRepository, snsClient)

	lambda.Start(service.ProcessKeyword)
}
