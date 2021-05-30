package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSes "github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
)

const CharSet = "UTF-8"

type Client struct {
	sesClient *awsSes.SES
}

// NewClient instantiates a SNS client
func NewClient(awsRegion string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	sesClient := awsSes.New(sess)
	xray.AWS(sesClient.Client)

	c := &Client{
		sesClient: sesClient,
	}

	return c, nil
}

func (c *Client) SendEmail(ctx context.Context, toAddresses, ccAddresses []string, sender, subject, htmlBody, textBody string) error {
	awsToAddresses := []*string{}
	for _, toAddress := range toAddresses {
		awsToAddresses = append(awsToAddresses, aws.String(toAddress))
	}

	awsCcAddresses := []*string{}
	for _, ccAddress := range ccAddresses {
		awsCcAddresses = append(awsCcAddresses, aws.String(ccAddress))
	}

	input := &awsSes.SendEmailInput{
		Destination: &awsSes.Destination{
			CcAddresses: awsCcAddresses,
			ToAddresses: awsToAddresses,
		},
		Message: &awsSes.Message{
			Body: &awsSes.Body{
				Html: &awsSes.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &awsSes.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &awsSes.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	result, err := c.sesClient.SendEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Infof("Email successfully sent: %v", result)
	return nil
}
