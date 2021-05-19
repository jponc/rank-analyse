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
)

type Service struct {
	repository *dbrepository.Repository
	snsClient  *sns.Client
}

func NewService(repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		repository: repository,
		snsClient:  snsClient,
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

	log.Infof("extracting result: %s, link: %s", result.ID.String(), result.Link)
}
