package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	AWSRegion           string
	RDSConnectionURL    string
	S3ResultsBucketName string
	SNSPrefix           string
	APIBaseURL          string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	rdsConnectionURL, err := getEnv("DB_CONN_URL")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	s3ResultsBucketName, err := getEnv("S3_RESULTS_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	snsPrefix, err := getEnv("SNS_PREFIX")
	if err != nil {
		return nil, err
	}

	apiBaseURL, err := getEnv("API_BASE_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion:           awsRegion,
		RDSConnectionURL:    rdsConnectionURL,
		S3ResultsBucketName: s3ResultsBucketName,
		SNSPrefix:           snsPrefix,
		APIBaseURL:          apiBaseURL,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
