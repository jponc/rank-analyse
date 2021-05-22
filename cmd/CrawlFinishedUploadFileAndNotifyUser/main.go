package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/repository/s3repository"
	"github.com/jponc/rank-analyse/internal/uploadnotify"
	"github.com/jponc/rank-analyse/pkg/postgres"
	"github.com/jponc/rank-analyse/pkg/s3"
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

	s3Client, err := s3.NewClient(config.AWSRegion)
	if err != nil {
		log.Fatalf("cannot initialise s3Client: %v", err)
	}

	s3Repository, err := s3repository.NewClient(s3Client, config.S3ResultsBucketName)
	if err != nil {
		log.Fatalf("cannot initialise s3Repository: %v", err)
	}

	service := uploadnotify.NewService(dbRepository, s3Repository)

	lambda.Start(service.CrawlFinishedUploadFileAndNotifyUser)
}
