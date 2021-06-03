package resultrankings

import (
	"context"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/types"
	"github.com/jponc/rank-analyse/pkg/sns"
	"github.com/jponc/rank-analyse/pkg/zenserp"
)

type Service struct {
	zenserpClient *zenserp.Client
	repository    *dbrepository.Repository
	snsClient     *sns.Client
}

func NewService(zenserpClient *zenserp.Client, repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		zenserpClient: zenserpClient,
		repository:    repository,
		snsClient:     snsClient,
	}

	return s
}

func (s *Service) ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent) {
	if s.zenserpClient == nil {
		log.Fatalf("zenserpClient not defined")
	}

	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.snsClient == nil {
		log.Fatalf("snsClient not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.ProcessKeywordMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	crawl, err := s.repository.CreateCrawl(
		ctx,
		msg.Keyword,
		msg.SearchEngine,
		msg.Device,
	)

	if err != nil {
		log.Fatalf("failed to create crawl: %v", err)
	}

	log.Infof("successfully created crawl with ID: %s", crawl.ID.String())

	res, err := s.zenserpClient.Search(
		ctx,
		msg.Keyword,
		msg.SearchEngine,
		msg.Device,
		msg.Count,
	)
	if err != nil {
		log.Fatalf("unable to query data from zenserp using keyword: %s", msg.Keyword)
	}

	resultItems := &types.ResultItemArray{}
	err = resultItems.Unmarshal(res)
	if err != nil {
		log.Fatalf("unable to unmarshal crawl result to result items: %v", err)
	}

	log.Infof("successfully unmarshalled zenserp res with length: %d", len(*resultItems))

	errorMsgs := []string{}
	results := []*types.Result{}

	// Iterate all result items and store to database
	for _, item := range *resultItems {
		// Store ResultItem to DB

		// Don't create a Result if there's no link
		if item.ItemURL == "" {
			continue
		}

		result, err := s.repository.CreateResult(ctx, crawl.ID, item.ItemURL, item.Title, item.Description, item.Position)
		if err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		} else {
			results = append(results, result)
		}
	}

	// Send ResultCreated message for all Results
	for _, result := range results {
		resultCreatedMsg := eventschema.ResultCreatedMessage{
			ResultID: result.ID.String(),
		}
		if err = s.snsClient.Publish(ctx, eventschema.ResultCreated, resultCreatedMsg); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		} else {
			log.Infof("Publishing result created: %s", result.ID.String())
		}
	}

	log.Infof("crawl results successfully created for keyword: %s, errors: %d", msg.Keyword, len(errorMsgs))

	if len(errorMsgs) > 0 {
		log.Errorf("errors encountered: %s", strings.Join(errorMsgs, "; "))
	}

	err = s.repository.Close()
	if err != nil {
		log.Fatalf("error closing connection: %v", err)
	}
}
