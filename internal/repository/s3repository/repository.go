package s3repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	awsS3Manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/pkg/s3"
)

type Repository struct {
	s3Client   *s3.Client
	bucketName string
}

// NewClient instantiates a Repository
func NewClient(s3Client *s3.Client, bucketName string) (*Repository, error) {
	r := &Repository{
		s3Client:   s3Client,
		bucketName: bucketName,
	}

	return r, nil
}

func (r *Repository) UploadCrawlResults(ctx context.Context, crawlID uuid.UUID, body string) error {
	key := fmt.Sprintf("crawls/%s.json", crawlID)

	input := &awsS3Manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(body),
	}

	_, err := r.s3Client.Upload(input)
	if err != nil {
		return fmt.Errorf("failed to upload crawl results JSON to s3: %v", err)
	}

	return nil
}

func (r *Repository) GetCrawlResultsURL(ctx context.Context, crawlID uuid.UUID) (string, error) {
	key := fmt.Sprintf("crawls/%s.json", crawlID)

	input := &awsS3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}

	req, _ := r.s3Client.GetObjectRequest(input)
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		return "", fmt.Errorf("failed to get crawl results JSON url: %v", err)
	}

	return urlStr, nil
}
