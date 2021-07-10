package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	ZenserpApiKey string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	zenserpApiKey, err := getEnv("ZENSERP_API_KEY")
	if err != nil {
		return nil, err
	}

	return &Config{
		ZenserpApiKey: zenserpApiKey,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
