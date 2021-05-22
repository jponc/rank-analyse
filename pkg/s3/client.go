package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
	s3Manager *s3manager.Uploader
	s3Client  *s3.S3
}

// NewClient instantiates an S3 client
func NewClient(awsRegion string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	s3Manager := s3manager.NewUploader(sess)
	s3Client := s3.New(sess)

	c := &Client{
		s3Client:  s3Client,
		s3Manager: s3Manager,
	}

	return c, nil
}

// Upload uploads content to s3 using concurrent upload
func (c *Client) Upload(input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	return c.s3Manager.Upload(input)
}

// GetObjectRequest creates get object request, used to create presigned urls
func (c *Client) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	return c.s3Client.GetObjectRequest(input)
}
