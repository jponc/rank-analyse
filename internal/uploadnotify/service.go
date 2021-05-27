package uploadnotify

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/api/eventschema"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/internal/repository/s3repository"
)

type ResultLink struct {
	LinkText string `json:"link_text"`
	LinkURL  string `json:"link_url"`
}

type ResultData struct {
	Title       string       `json:"title"`
	Position    int          `json:"position"`
	URL         string       `json:"url"`
	Links       []ResultLink `json:"links"`
	BodyText    string       `json:"body_text"`
	QueryDate   time.Time    `json:"query_date"`
	QuerySearch string       `json:"query_search"`
}

type Service struct {
	dbrepository *dbrepository.Repository
	s3repository *s3repository.Repository
}

func NewService(dbrepository *dbrepository.Repository, s3repository *s3repository.Repository) *Service {
	s := &Service{
		dbrepository: dbrepository,
		s3repository: s3repository,
	}

	return s
}

func (s *Service) CrawlFinishedUploadFileAndNotifyUser(ctx context.Context, snsEvent events.SNSEvent) {
	if s.dbrepository == nil {
		log.Fatalf("dbrepository not defined")
	}

	if s.s3repository == nil {
		log.Fatalf("s3repository not defined")
	}

	if err := s.dbrepository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.CrawlFinishedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	crawlID, err := uuid.FromString(msg.CrawlID)
	if err != nil {
		log.Fatalf("failed to get Crawl UUID: %v", err)
	}

	crawl, err := s.dbrepository.GetCrawl(ctx, crawlID)
	if err != nil {
		log.Fatalf("failed to get Crawl: %v", err)
	}

	results, err := s.dbrepository.GetCrawlResults(ctx, crawl.ID)
	if err != nil {
		log.Fatalf("failed to get crawl results: %v", err)
	}

	resultDatas := []ResultData{}

	// Iterate each result and get resultDatas
	for _, result := range *results {
		// Get ExtractInfo
		extractInfo, err := s.dbrepository.GetExtractInfo(ctx, result.ID)
		if err != nil {
			log.Fatalf("failed to get extract info (%s): %v", result.ID, err)
		}

		// Get ExtractLinks
		extractLinks, err := s.dbrepository.GetExtractLinks(ctx, result.ID)
		if err != nil {
			log.Fatalf("failed to get extract links (%s): %v", result.ID, err)
		}

		// Create ResultData
		links := []ResultLink{}

		for _, extractLink := range *extractLinks {
			links = append(links, ResultLink{
				LinkText: extractLink.Text,
				LinkURL:  extractLink.LinkURL,
			})
		}

		d := ResultData{
			Title:       extractInfo.Title,
			Position:    result.Position,
			URL:         result.Link,
			BodyText:    extractInfo.Content,
			QueryDate:   result.CreatedAt,
			QuerySearch: crawl.Keyword,
			Links:       links,
		}

		// Append to resultDatas
		resultDatas = append(resultDatas, d)

	}

	// Marshal to JSON
	resultDatasJson, err := json.Marshal(&resultDatas)
	if err != nil {
		log.Fatalf("failed to marshal result datas to JSON: %v", err)
	}

	// Upload to S3
	err = s.s3repository.UploadCrawlResults(ctx, crawl.ID, string(resultDatasJson))
	if err != nil {
		log.Fatalf("failed to upload results datas to S3: %v", err)
	}

	err = s.dbrepository.Close()
	if err != nil {
		log.Fatalf("error closing connection: %v", err)
	}

	log.Infof("finished uploading to s3 (%s)", crawl.ID)
}
