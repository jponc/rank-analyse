package analysis

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/textrazor"
)

type Service struct {
	repository      *dbrepository.Repository
	textrazorClient *textrazor.Client
}

func NewService(repository *dbrepository.Repository, textrazorClient *textrazor.Client) *Service {
	s := &Service{
		repository:      repository,
		textrazorClient: textrazorClient,
	}

	return s
}

func (s *Service) ResultCreatedRunAnalysis(ctx context.Context, snsEvent events.SNSEvent) {
	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.textrazorClient == nil {
		log.Fatalf("textrazorClient not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}
	defer s.repository.Close()

	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.ResultCreatedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	resultID, err := uuid.FromString(msg.ResultID)
	if err != nil {
		log.Fatalf("failed to get result UUID: %v", err)
	}

	result, err := s.repository.GetResult(ctx, resultID)
	if err != nil {
		log.Fatalf("failed to get result: %v", err)
	}

	log.Infof("Processing Result with ID (%s)", resultID.String())

	if result.Link == "" {
		log.Fatalf("can't process analysis for empty link")
	}

	extractors := []textrazor.Extractor{textrazor.Entities, textrazor.Topics}

	analyzeResponse, err := s.textrazorClient.Analyze(ctx, result.Link, extractors)
	if err != nil {
		log.Fatalf("failed to run textrazor analyze: %v", err)
	}

	var topicArr types.AnalyzeTopicArray
	err = topicArr.Unmarshal(&analyzeResponse.Topics)
	if err != nil {
		log.Fatalf("failed to unmarshal topics: %v", err)
	}

	var entitiesArr types.AnalyzeEntityArray
	err = entitiesArr.Unmarshal(&analyzeResponse.Entities)
	if err != nil {
		log.Fatalf("failed to unmarshal entities: %v", err)
	}

	err = s.repository.CreateAnalyzeTopics(ctx, resultID, topicArr)
	if err != nil {
		log.Fatalf("failed to save topics: %v", err)
	}

	err = s.repository.CreateAnalyzeEntities(ctx, resultID, entitiesArr)
	if err != nil {
		log.Fatalf("failed to save entities: %v", err)
	}

	err = s.repository.SaveCleanedText(ctx, resultID, analyzeResponse.CleanedText)
	if err != nil {
		log.Fatalf("failed to save cleaned text: %v", err)
	}

	log.Infof("successfully save analaysis data")
}
