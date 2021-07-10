package zenserp

import (
	"context"
	"fmt"
)

func (c *Client) Batch(ctx context.Context, name, webhookUrl string, jobs []Job) (*BatchResult, error) {

	batchRequest := &BatchRequest{
		WebhookURL: webhookUrl,
		Name:       name,
		Jobs:       jobs,
	}

	var batchResult BatchResult

	err := c.postJSON(ctx, batchPath, batchRequest, &batchResult)
	if err != nil {
		return nil, fmt.Errorf("failed to post batch: %w", err)
	}

	return &batchResult, nil
}
