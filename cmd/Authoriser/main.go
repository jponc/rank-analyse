package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/jponc/rank-analyse/internal/auth"
	"github.com/jponc/rank-analyse/internal/authoriser"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	authClient, err := auth.NewClient(config.JWTSecret)
	if err != nil {
		log.Fatalf("cannot initialise auth client %v", err)
	}

	service := authoriser.NewService(authClient)
	lambda.Start(service.Authorise)
}
