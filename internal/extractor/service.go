package extractor

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/pkg/sns"
	"github.com/jponc/rank-analyse/pkg/webscraper"
)

type Service struct {
	repository    *dbrepository.Repository
	snsClient     *sns.Client
	scraperClient *webscraper.Client
}

func NewService(repository *dbrepository.Repository, snsClient *sns.Client, scraperClient *webscraper.Client) *Service {
	s := &Service{
		repository:    repository,
		snsClient:     snsClient,
		scraperClient: scraperClient,
	}

	return s
}

func (s *Service) ResultCreatedExtractPageInfo(ctx context.Context, snsEvent events.SNSEvent) {
	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.snsClient == nil {
		log.Fatalf("snsClient not defined")
	}

	if s.scraperClient == nil {
		log.Fatalf("scraperClient not defined")
	}

	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.ResultCreatedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	resultId, err := uuid.FromString(msg.ResultID)
	if err != nil {
		log.Fatalf("failed to get result UUID: %v", err)
	}

	result, err := s.repository.GetResult(ctx, resultId)
	if err != nil {
		log.Fatalf("failed to get result: %v", err)
	}

	if result.Link == "" {
		log.Errorf("can't extract for result %s because there's no link specified", result.ID)
		return
	}

	// Run scraping
	scrapeResult, err := s.scraperClient.Scrape(ctx, result.Link)
	if err != nil {
		log.Fatalf("error scraping: %v", err)
	}

	// Store ExtractInfo
	err = s.repository.CreateExtractInfo(ctx, result.ID, scrapeResult.Title, scrapeResult.Body)
	if err != nil {
		log.Fatalf("error creating extract info: %v", err)
	}

	// Store ExtractLinks
	err = s.repository.CreateExtractLinks(ctx, result.ID, scrapeResult.Links)
	if err != nil {
		log.Fatalf("error creating extract links: %v", err)
	}

	log.Infof("finished extracting %s (%s)", result.ID.String(), result.Link)
}
