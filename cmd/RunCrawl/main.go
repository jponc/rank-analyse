package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-analyse/internal/api"
)

func main() {
	service := api.NewService(nil)
	lambda.Start(service.RunCrawl)
}
