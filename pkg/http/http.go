package http

import (
	"net/http"
	"time"
)

// DefaultHTTPClient returns http client with default settings
func DefaultHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
