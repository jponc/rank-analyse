package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	ZenserpBatchWebhookURL string
	ZenserpApiKey          string
	Locations              []string
	Country                string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	zenserpBatchWebhookURL, err := getEnv("ZENSERP_BATCH_WEBHOOK_URL")
	if err != nil {
		return nil, err
	}

	zenserpApiKey, err := getEnv("ZENSERP_API_KEY")
	if err != nil {
		return nil, err
	}

	country := "US"

	locations := []string{
		"Mather,California,United States",
		"Melstone,Montana,United States",
		"Austin County,Texas,United States",
		"Denton,North Carolina,United States",
		"Kingfield,Maine,United States",
	}

	return &Config{
		ZenserpBatchWebhookURL: zenserpBatchWebhookURL,
		ZenserpApiKey:          zenserpApiKey,
		Country:                country,
		Locations:              locations,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
